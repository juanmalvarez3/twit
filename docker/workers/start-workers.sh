#!/bin/sh
set -e

echo "Iniciando workers para Twit..."

# Esperar que los servicios AWS estén disponibles antes de comenzar, con menos intentos
MAX_RETRIES=10
RETRY=0

# Función mejorada para verificar si LocalStack está listo
check_localstack() {
    # Verificación más rápida usando timeout para evitar esperas largas
    if curl --max-time 2 -s http://localstack:4566/_localstack/health | grep -q "dynamodb.*running"; then
        return 0
    fi
    return 1
}

echo "Verificando disponibilidad de LocalStack..."
while [ $RETRY -lt $MAX_RETRIES ]; do
    if check_localstack; then
        echo "LocalStack está listo y funcionando correctamente."
        break
    fi
    echo "Esperando a que LocalStack esté listo... (intento $RETRY de $MAX_RETRIES)"
    sleep 2
    RETRY=$((RETRY+1))
done

# Continuar incluso si LocalStack no responde como esperado
if [ $RETRY -eq $MAX_RETRIES ]; then
    echo "ADVERTENCIA: LocalStack no respondió correctamente después de $MAX_RETRIES intentos."
    echo "Continuando de todos modos ya que esto suele funcionar..."
fi

# Iniciar workers en background
#echo "Iniciando worker: orchestrate-fanout"
#/bin/orchestratefanout &
#ORCHESTRATE_PID=$!

echo "Iniciando worker: update-timeline"
/bin/update-timeline &
UPDATE_TIMELINE_PID=$!

echo "Iniciando worker SNS: tweets"
/bin/tweets &
TWEETS_PID=$!

echo "Iniciando worker SNS: follows"
/bin/follows &
FOLLOWS_PID=$!

# echo "Iniciando worker: process-new-follow"
# /bin/process-new-follow &
# PROCESS_FOLLOW_PID=$!

echo "Iniciando worker: populate-cache"
/bin/populate-cache &
POPULATE_CACHE_PID=$!

echo "Iniciando worker: rebuild-timeline"
/bin/rebuild-timeline &
REBUILD_TIMELINE_PID=$!

echo "Todos los workers iniciados correctamente."

# Función para manejar señales
handle_signal() {
    echo "Recibida señal para terminar, deteniendo workers..."
    kill $UPDATE_TIMELINE_PID $REBUILD_TIMELINE_PID $TWEETS_PID $FOLLOWS_PID $POPULATE_CACHE_PID 2>/dev/null || true
    wait
    echo "Todos los workers detenidos."
    exit 0
}

# Registrar manejador de señales
trap handle_signal SIGTERM SIGINT

# Mantener el script en ejecución para que Docker no detenga el contenedor
wait
