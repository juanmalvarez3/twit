#!/bin/bash
set -euo pipefail

echo "==== Configurando entorno local para el proyecto Twit ===="

# Verificar si docker y docker-compose están instalados
if ! command -v docker &> /dev/null; then
    echo "Error: Docker no está instalado. Por favor instálalo antes de continuar."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "Error: Docker Compose no está instalado. Por favor instálalo antes de continuar."
    exit 1
fi

# Asegurar que el script de inicialización tenga permisos de ejecución
chmod +x ./localstack/init-aws.sh

# Crear archivo .env desde el ejemplo si no existe
if [ ! -f .env ]; then
    echo "Creando archivo .env desde .env.example..."
    cp .env.example .env
    echo "Archivo .env creado. Por favor, revísalo y ajusta los valores según sea necesario."
fi

# Detener contenedores si están corriendo
echo "Deteniendo contenedores si están corriendo..."
docker-compose down

# Iniciar servicios
echo "Iniciando servicios locales (Redis y LocalStack)..."
docker-compose up -d

echo "Esperando a que los servicios estén listos..."
sleep 10

# Verificar estado de los servicios
echo "Verificando estado de LocalStack..."
curl -s http://localhost:4566/_localstack/health

echo "Verificando estado de Redis..."
docker exec redis-twit redis-cli ping

echo "==== Configuración local completada ===="
echo "Servicios disponibles:"
echo "- LocalStack: http://localhost:4566"
echo "- Redis: localhost:6379"
echo ""
echo "Para detener los servicios: docker-compose down"
