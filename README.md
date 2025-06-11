# Twit - Plataforma de Microblogging

Implementación de una plataforma de microblogging estilo Twitter siguiendo principios de Clean Architecture.

## Guía de Instalación y Ejecución con Docker

### Prerequisitos

- [Docker](https://www.docker.com/get-started) (20.10.0 o superior)
- [Docker Compose](https://docs.docker.com/compose/install/) (v2.0.0 o superior)
- Git

### Paso 1: Clonar el repositorio

```bash
git clone https://github.com/juanmalvarez3/twit.git
cd twit
```

### Paso 2: Configuración del entorno

Copia el archivo de configuración de ejemplo:

```bash
cp .env.example .env
```

> **Nota**: Por defecto, la configuración incluida está optimizada para desarrollo local con Docker.

### Paso 3: Construir e iniciar los servicios

```bash
# Construir las imágenes
docker-compose build

# Iniciar todos los servicios en modo detached
docker-compose up -d
```

Este comando iniciará los siguientes servicios:
- **API REST** en http://localhost:8080
- **Workers** para procesamiento asíncrono
  - orchestrate-fanout: Distribuye tweets a seguidores
  - update-timeline: Actualiza timelines individuales
  - process-new-follow: Procesa nuevas relaciones de seguimiento
  - populate-cache: Prepara caché de timelines
  - rebuild-timeline: Reconstruye timelines
- **LocalStack** (emulador de AWS) en puerto 4566
  - DynamoDB: Para almacenar tweets, timelines, follows y usuarios
  - SNS: Para notificaciones
  - SQS: Para colas de mensajes
- **Redis** en puerto 6379 para caché de timelines

### Paso 4: Verificar que todos los servicios estén funcionando

```bash
# Ver el estado de todos los contenedores
docker-compose ps

# Verificar logs
docker-compose logs -f

# Ver logs de un servicio específico (ej: api)
docker-compose logs -f api
```

### Paso 5: Realizar peticiones a la API

Puedes usar curl o cualquier cliente HTTP como Postman para probar los endpoints:

```bash
# Crear un tweet
curl -X POST http://localhost:8080/api/v1/tweets \
  -H "Content-Type: application/json" \
  -d '{"userId": "user123", "content": "¡Hola mundo!"}'

# Seguir a un usuario
curl -X POST http://localhost:8080/api/v1/follows \
  -H "Content-Type: application/json" \
  -d '{"followerId": "user456", "followedId": "user123"}'

# Obtener el timeline de un usuario
curl http://localhost:8080/api/v1/timelines/user456
```

### Paso 6: Gestión de los servicios

```bash
# Detener todos los servicios (preservando volúmenes)
docker-compose down

# Detener y eliminar volúmenes (útil para empezar de cero)
docker-compose down -v

# Reiniciar un servicio específico
docker-compose restart api

# Ver logs en tiempo real
docker-compose logs -f
```

### Paso 7: Purgar datos de caché (Redis)

Si necesitas limpiar la caché de Redis:

```bash
docker-compose exec redis redis-cli FLUSHALL
```

### Solución de problemas comunes

- **LocalStack no inicializa correctamente**: Verifica los logs con `docker-compose logs -f localstack` y asegúrate de que Docker tenga suficientes recursos asignados

- **Errores de conexión**: Asegúrate de que los servicios estén en la misma red de Docker y que las variables de entorno están configuradas correctamente

- **Problemas de caché**: Intenta purgar Redis como se indica en el Paso 7

## Arquitectura

El sistema está diseñado siguiendo los principios de Clean Architecture con los siguientes componentes:

- **Casos de uso**: Lógica de aplicación independiente de la infraestructura
- **Servicios de dominio**: Reglas de negocio y lógica específica del dominio
- **Repositorios**: Interfaces para acceso a datos
- **Adaptadores**: Implementaciones concretas para HTTP, DynamoDB, Redis, SNS/SQS

## API Endpoints

- `POST /api/v1/tweets`
  - Crear un nuevo tweet
  - Body: `{"userId": "user123", "content": "¡Hola mundo!"}`

- `POST /api/v1/follows`
  - Seguir a un usuario
  - Body: `{"followerId": "user123", "followedId": "user456"}`

- `GET /api/timelines/{userID}`
  - Obtener el timeline de un usuario
  - Parámetros opcionales: `limit`, `cursor`

## Estructura del proyecto

```
twit/
├── cmd/                      # Punto de entrada de la aplicación
│   ├── twitter/              # Servicios principales
│       ├── http/             # API HTTP REST
│       ├── snssqs/           # Procesadores de mensajes SNS/SQS
│       │   ├── follows/      # Procesador de eventos de follows
│       │   └── tweets/       # Procesador de eventos de tweets
│       └── sqs/              # Workers para procesamiento asíncrono
│           ├── orchestratefanout/ # Distribución de tweets a seguidores
│           ├── populatecache/    # Población de caché de timeline
│           ├── rebuildtimeline/  # Reconstrucción de timelines
│           └── updatetimeline/   # Actualización de timelines
├── internal/                 # Código privado de la aplicación
│   ├── adapters/             # Adaptadores para infraestructura
│   │   ├── queue/           # Adaptadores para colas de mensajes
│   │   ├── redis/           # Adaptador para Redis
│   │   └── sns/             # Adaptadores para SNS
│   └── domains/              # Dominios de negocio
│       └── twitter/          # Dominio principal
│           ├── follow/        # Subdominio de seguimientos
│           ├── timeline/      # Subdominio de timeline
│           └── tweet/         # Subdominio de tweets
├── pkg/                      # Código público reutilizable
│   ├── config/               # Gestión de configuración
│   ├── errors/               # Manejo de errores
│   └── logger/               # Sistema de logging
├── docker/                   # Archivos Docker y scripts
├── localstack/               # Configuración de servicios locales
├── docker-compose.yml        # Definición de servicios Docker
└── setup-local.sh           # Script para configuración local
```

Cada dominio (tweet, follow, timeline) sigue una estructura interna consistente:

```
dominio/
├── domain/           # Definiciones y entidades del dominio
├── repository/       # Interfaces e implementaciones de repositorios
├── service/          # Servicios del dominio
└── usecases/         # Casos de uso para operaciones específicas
    ├── usecase1/     # Implementación de caso de uso
    │   ├── interfaces.go # Interfaces requeridas
    │   ├── provide.go   # Proveedor de dependencias
    │   └── exec.go      # Implementación del caso de uso
    └── usecase2/
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