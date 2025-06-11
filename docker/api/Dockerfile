FROM golang:1.22-alpine AS builder

# Instalar dependencias del sistema
RUN apk add --no-cache git ca-certificates tzdata

# Configurar directorio de trabajo
WORKDIR /app

# Copia solo archivos de dependencias para mejor caché
COPY go.mod go.sum ./

# Usa montado de caché para acelerar descargas de dependencias
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download

# Copiar solo lo necesario para compilar la API
COPY pkg/ ./pkg/
COPY internal/ ./internal/
COPY cmd/twitter/http/ ./cmd/twitter/http/

# Compilar la aplicación con optimizaciones
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/api -trimpath ./cmd/twitter/http/main.go

# Imagen final
FROM alpine:latest

# Añadir dependencias básicas
RUN apk --no-cache add ca-certificates tzdata

# Copiar el binario compilado desde el builder
COPY --from=builder /bin/api /bin/api

# Crear directorios para configuración
RUN mkdir -p /config

# Exponer el puerto
EXPOSE 8080

# Ejecutar la API
CMD ["/bin/api"]
