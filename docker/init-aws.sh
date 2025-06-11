#!/bin/sh

# Script para inicializar recursos AWS en LocalStack
echo "Esperando a que LocalStack esté listo..."
MAX_RETRIES=30
RETRY_COUNT=0

# Intentar con un endpoint más específico que devuelve un estado más confiable
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  if curl -s http://localstack:4566/health | grep -q "\"running\": true"; then
    echo "LocalStack está listo!"
    break
  elif curl -s http://localstack:4566/_localstack/health | grep -q "\"ready\": true"; then
    echo "LocalStack está listo (endpoint alternativo)!"
    break
  fi
  
  echo "Esperando a LocalStack... (intento $RETRY_COUNT de $MAX_RETRIES)"
  RETRY_COUNT=$((RETRY_COUNT + 1))
  sleep 2
  
  # Si llegamos al máximo de intentos, continuamos de todos modos
  if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo "Alcanzado número máximo de intentos. Continuando de todos modos..."
  fi
done

echo "Creando tablas DynamoDB..."

# Crear tabla de tweets con GSI para búsquedas por usuario
echo "Creando tabla 'tweets'..."
aws --endpoint-url=http://localstack:4566 --region us-east-1 dynamodb create-table \
  --table-name tweets \
  --attribute-definitions \
      AttributeName=id,AttributeType=S \
      AttributeName=user_id,AttributeType=S \
      AttributeName=created_at,AttributeType=S \
  --key-schema AttributeName=id,KeyType=HASH \
  --global-secondary-indexes \
      "[{\
          \"IndexName\": \"user_id-created_at-index\",\
          \"KeySchema\": [{\"AttributeName\":\"user_id\",\"KeyType\":\"HASH\"}, {\"AttributeName\":\"created_at\",\"KeyType\":\"RANGE\"}],\
          \"Projection\": {\"ProjectionType\":\"ALL\"},\
          \"ProvisionedThroughput\": {\"ReadCapacityUnits\":5,\"WriteCapacityUnits\":5}\
        }]" \
  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 || echo "Error al crear tabla tweets, puede que ya exista"

# Crear tabla de usuarios
echo "Creando tabla 'users'..."
aws --endpoint-url=http://localstack:4566 --region us-east-1 dynamodb create-table \
  --table-name users \
  --attribute-definitions AttributeName=id,AttributeType=S \
  --key-schema AttributeName=id,KeyType=HASH \
  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 || echo "Error al crear tabla users, puede que ya exista"

# Crear tabla de follows con GSI para consultas bidireccionales
echo "Creando tabla 'follows' con índices secundarios globales..."
aws --endpoint-url=http://localstack:4566 --region us-east-1 dynamodb create-table \
  --table-name follows \
  --attribute-definitions \
      AttributeName=follower_id,AttributeType=S \
      AttributeName=followed_id,AttributeType=S \
  --key-schema AttributeName=follower_id,KeyType=HASH AttributeName=followed_id,KeyType=RANGE \
  --global-secondary-indexes \
      "[\
        {\
          \"IndexName\": \"followed_id-index\",\
          \"KeySchema\": [{\"AttributeName\":\"followed_id\",\"KeyType\":\"HASH\"}],\
          \"Projection\": {\"ProjectionType\":\"ALL\"},\
          \"ProvisionedThroughput\": {\"ReadCapacityUnits\":5,\"WriteCapacityUnits\":5}\
        }\
      ]" \
  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 || echo "Error al crear tabla follows, puede que ya exista"

# Crear tabla de timelines
echo "Creando tabla 'timelines'..."
aws --endpoint-url=http://localstack:4566 --region us-east-1 dynamodb create-table \
  --table-name timelines \
  --attribute-definitions AttributeName=user_id,AttributeType=S AttributeName=tweet_id,AttributeType=S \
  --key-schema AttributeName=user_id,KeyType=HASH AttributeName=tweet_id,KeyType=RANGE \
  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 || echo "Error al crear tabla timelines, puede que ya exista"

echo "Listando tablas DynamoDB creadas:"
aws --endpoint-url=http://localstack:4566 --region us-east-1 dynamodb list-tables

echo "Creando colas SQS..."
aws --endpoint-url=http://localstack:4566 --region us-east-1 sqs create-queue --queue-name orchestrate-fanout || echo "Error al crear cola orchestrate-fanout, puede que ya exista"
aws --endpoint-url=http://localstack:4566 --region us-east-1 sqs create-queue --queue-name update-timeline || echo "Error al crear cola update-timeline, puede que ya exista"
aws --endpoint-url=http://localstack:4566 --region us-east-1 sqs create-queue --queue-name process-new-follow || echo "Error al crear cola process-new-follow, puede que ya exista"
aws --endpoint-url=http://localstack:4566 --region us-east-1 sqs create-queue --queue-name populate-cache || echo "Error al crear cola populate-cache, puede que ya exista"
aws --endpoint-url=http://localstack:4566 --region us-east-1 sqs create-queue --queue-name rebuild-timeline || echo "Error al crear cola rebuild-timeline, puede que ya exista"

echo "Creando temas SNS..."
aws --endpoint-url=http://localstack:4566 --region us-east-1 sns create-topic --name tweets || echo "Error al crear tema tweets, puede que ya exista"
aws --endpoint-url=http://localstack:4566 --region us-east-1 sns create-topic --name follows || echo "Error al crear tema follows, puede que ya exista"

echo "Listando colas SQS creadas:"
aws --endpoint-url=http://localstack:4566 --region us-east-1 sqs list-queues

echo "Listando temas SNS creados:"
aws --endpoint-url=http://localstack:4566 --region us-east-1 sns list-topics

# Crear suscripciones SNS a SQS
echo "Creando suscripciones SNS-SQS..."

# Obtener ARNs de temas y colas de forma simplificada
TWEETS_TOPIC_ARN="arn:aws:sns:us-east-1:000000000000:tweets"
FOLLOWS_TOPIC_ARN="arn:aws:sns:us-east-1:000000000000:follows"

# Usar URLs y ARNs conocidos para simplificar
ORCHESTRATE_QUEUE_URL="http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/orchestrate-fanout"
PROCESS_FOLLOW_QUEUE_URL="http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/process-new-follow"

ORCHESTRATE_QUEUE_ARN="arn:aws:sqs:us-east-1:000000000000:orchestrate-fanout"
PROCESS_FOLLOW_QUEUE_ARN="arn:aws:sqs:us-east-1:000000000000:process-new-follow"

# Suscribir colas a temas
echo "Suscribiendo cola orchestrate-fanout al tema tweets..."
aws --endpoint-url=http://localstack:4566 --region us-east-1 sns subscribe \
  --topic-arn "$TWEETS_TOPIC_ARN" \
  --protocol sqs \
  --notification-endpoint "$ORCHESTRATE_QUEUE_ARN" \
  || echo "Error al crear suscripción de tweets a orchestrate-fanout, puede que ya exista"

echo "Suscribiendo cola process-new-follow al tema follows..."
aws --endpoint-url=http://localstack:4566 --region us-east-1 sns subscribe \
  --topic-arn "$FOLLOWS_TOPIC_ARN" \
  --protocol sqs \
  --notification-endpoint "$PROCESS_FOLLOW_QUEUE_ARN" \
  || echo "Error al crear suscripción de follows a process-new-follow, puede que ya exista"

# Con LocalStack no es necesario configurar políticas SQS para SNS
# Las suscripciones funcionarán directamente
echo "Nota: En LocalStack no se requieren políticas para las colas SQS, las suscripciones funcionan directamente."

echo "Listando suscripciones SNS creadas:"
aws --endpoint-url=http://localstack:4566 --region us-east-1 sns list-subscriptions

echo "Recursos AWS creados correctamente"
exit 0
