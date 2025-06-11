FROM golang:1.22-alpine AS builder

# Instalar dependencias del sistema y herramientas de compilación
RUN apk add --no-cache git ca-certificates tzdata make

# Configurar directorio de trabajo
WORKDIR /app

# Copia solo archivos de dependencias para mejor caché
COPY go.mod go.sum ./

# Descarga dependencias en una capa separada 
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download

# Copiar archivos específicos de estructura primero (solo los necesarios para compilación)
COPY pkg/ ./pkg/
COPY internal/ ./internal/
COPY cmd/ ./cmd/

# Crear directorio para los binarios
RUN mkdir -p /bin/workers

# Compilación paralela de workers SQS - mucho más rápido
RUN echo "Compilando workers en paralelo..." && \
    (CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/workers/orchestratefanout ./cmd/twitter/sqs/orchestratefanout/main.go & \
     CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/workers/updatetimeline ./cmd/twitter/sqs/updatetimeline/main.go & \
     CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/workers/populatecache ./cmd/twitter/sqs/populatecache/main.go & \
     wait)

# Compilación paralela de workers SNS
RUN echo "Compilando workers SNS en paralelo..." && \
    (CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/workers/tweets ./cmd/twitter/snssqs/tweets/main.go & \
     CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/workers/follows ./cmd/twitter/snssqs/follows/main.go & \
     wait)

# Imagen final
FROM alpine:latest

# Añadir dependencias básicas
RUN apk --no-cache add ca-certificates tzdata curl

# Copiar los binarios compilados desde el builder (en una sola capa)
COPY --from=builder /bin/workers/orchestratefanout /bin/orchestratefanout
COPY --from=builder /bin/workers/updatetimeline /bin/update-timeline
COPY --from=builder /bin/workers/tweets /bin/tweets
COPY --from=builder /bin/workers/follows /bin/follows
COPY --from=builder /bin/workers/populatecache /bin/populate-cache

# Copiar script de inicio para los workers
COPY ./docker/workers/start-workers.sh /bin/start-workers.sh
RUN chmod +x /bin/start-workers.sh

# Crear directorios para configuración
RUN mkdir -p /config

# Ejecutar todos los workers mediante el script
CMD ["/bin/start-workers.sh"]
