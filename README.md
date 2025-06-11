# Twit - Plataforma de Microblogging

Implementación de una plataforma de microblogging estilo Twitter siguiendo principios de Clean Architecture.

## Ejecución rápida (Docker)

Para ejecutar todo el proyecto con un solo comando:

```bash
# Clonar el repositorio
git clone https://github.com/juanmalvarez3/twit.git
cd twit

# Ejecutar con Docker Compose
docker-compose up -d
```

Esto iniciará:
- API REST en http://localhost:8080
- Servicios de workers para procesamiento asíncrono
- LocalStack (DynamoDB, SNS, SQS) en puerto 4566
- Redis en puerto 6379

Para detener todos los servicios:
```bash
docker-compose down
```

## Arquitectura

El sistema está diseñado siguiendo los principios de Clean Architecture con los siguientes componentes:

- **Casos de uso**: Lógica de aplicación independiente de la infraestructura
- **Servicios de dominio**: Reglas de negocio y lógica específica del dominio
- **Repositorios**: Interfaces para acceso a datos
- **Adaptadores**: Implementaciones concretas para HTTP, DynamoDB, Redis, SNS/SQS

## API Endpoints

- `POST /api/tweets`
  - Crear un nuevo tweet
  - Body: `{"user_id": "user123", "content": "¡Hola mundo!"}`

- `POST /api/follows`
  - Seguir a un usuario
  - Body: `{"follower_id": "user123", "followed_id": "user456"}`

- `GET /api/timelines/{userID}`
  - Obtener el timeline de un usuario
  - Parámetros opcionales: `limit`, `cursor`

## Estructura del proyecto

```
twit/
├── cmd/                      # Punto de entrada de la aplicación
│   ├── twitter/              # Servicios principales
│   │   ├── http/             # API HTTP
│   │   ├── sns/              # Procesadores de eventos SNS
│   │   ├── sqs/              # Procesadores de colas SQS
│   │   └── scheduled/        # Tareas programadas
├── internal/                 # Código privado de la aplicación
│   ├── domains/              # Dominios de negocio
│   │   ├── twitter/          # Dominio principal
│   │   │   ├── tweet/        # Subdominio de tweets
│   │   │   ├── follow/       # Subdominio de follows
│   │   │   ├── timeline/     # Subdominio de timeline
│   │   │   └── user/         # Subdominio de usuarios
│   └── adapters/             # Adaptadores para infraestructura
├── pkg/                      # Código público reutilizable
│   ├── config/               # Gestión de configuración
│   ├── errors/               # Manejo de errores
│   └── logger/               # Sistema de logging
├── docker/                   # Archivos Docker
└── localstack/               # Configuración de servicios locales
```

## Ejecución en modo desarrollo

Si prefieres ejecutar la aplicación sin Docker para desarrollo:

```bash
# Instalar dependencias
go mod download

# Levantar servicios de infraestructura (LocalStack, Redis)
./setup-local.sh

# Ejecutar el servidor API
go run cmd/twitter/http/main.go

# En otra terminal, ejecutar los workers (ejemplo para uno de ellos)
go run cmd/twitter/sqs/orchestratefanout/main.go
```

## Componentes técnicos

- **Infraestructura simulada**:
  - LocalStack (DynamoDB, SNS, SQS)
  - Redis para caché de timelines

- **Tablas de DynamoDB**:
  - `tweets`: Almacena todos los tweets (PK=tweet_id, SK=created_at)
  - `follows`: Relaciones entre usuarios (PK=follower_id, SK=followed_id)
  - `timelines`: Timeline por usuario (PK=user_id, SK=created_at_tweet_id)
  - `users`: Información de usuarios (PK=user_id)

- **Tópicos SNS**:
  - `tweets`: Notifica eventos relacionados con tweets
  - `follows`: Notifica eventos relacionados con relaciones de seguimiento

- **Colas SQS**:
  - `orchestrate-fanout`: Distribuye tweets a seguidores
  - `update-timeline`: Actualiza timelines individuales
  - `process-new-follow`: Procesa nuevas relaciones de seguimiento
  - `populate-cache`: Prepara caché de timelines
  - `rebuild-timeline`: Reconstruye timelines desde datos persistentes