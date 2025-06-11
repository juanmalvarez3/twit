#!/bin/bash
set -euo pipefail

echo "Iniciando configuración de servicios AWS locales..."

# Variables de configuración
AWS_REGION="us-east-1"
ENDPOINT_URL="http://localstack:4566"
AWS_ACCOUNT_ID="000000000000"  # ID de cuenta ficticio para LocalStack

# Credenciales de prueba
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=$AWS_REGION

# Función de utilidad para verificar si un recurso existe
resource_exists() {
  local output
  output=$(eval "$1" 2>&1 || true)
  echo "$output" | grep -v "ResourceNotFoundException" | grep -v "NonExistentQueue" > /dev/null
}

# Crear tablas DynamoDB
echo "Creando tablas DynamoDB..."

# Tabla tweets
if ! resource_exists "aws dynamodb describe-table --table-name tweets --endpoint-url $ENDPOINT_URL"; then
  aws dynamodb create-table \
    --table-name tweets \
    --attribute-definitions \
      AttributeName=tweet_id,AttributeType=S \
      AttributeName=created_at,AttributeType=S \
    --key-schema \
      AttributeName=tweet_id,KeyType=HASH \
      AttributeName=created_at,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --endpoint-url $ENDPOINT_URL
  echo "Tabla 'tweets' creada"
else
  echo "La tabla 'tweets' ya existe"
fi

# Tabla timelines
if ! resource_exists "aws dynamodb describe-table --table-name timelines --endpoint-url $ENDPOINT_URL"; then
  aws dynamodb create-table \
    --table-name timelines \
    --attribute-definitions \
      AttributeName=user_id,AttributeType=S \
      AttributeName=created_at_tweet_id,AttributeType=S \
    --key-schema \
      AttributeName=user_id,KeyType=HASH \
      AttributeName=created_at_tweet_id,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --endpoint-url $ENDPOINT_URL
  
  # Habilitar TTL
  aws dynamodb update-time-to-live \
    --table-name timelines \
    --time-to-live-specification "Enabled=true, AttributeName=ttl" \
    --endpoint-url $ENDPOINT_URL
  
  echo "Tabla 'timelines' creada con TTL habilitado"
else
  echo "La tabla 'timelines' ya existe"
fi

# Tabla follows
if ! resource_exists "aws dynamodb describe-table --table-name follows --endpoint-url $ENDPOINT_URL"; then
  aws dynamodb create-table \
    --table-name follows \
    --attribute-definitions \
      AttributeName=follower_id,AttributeType=S \
      AttributeName=followed_id,AttributeType=S \
    --key-schema \
      AttributeName=follower_id,KeyType=HASH \
      AttributeName=followed_id,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --endpoint-url $ENDPOINT_URL
  echo "Tabla 'follows' creada"
else
  echo "La tabla 'follows' ya existe"
fi

# Tabla users
if ! resource_exists "aws dynamodb describe-table --table-name users --endpoint-url $ENDPOINT_URL"; then
  aws dynamodb create-table \
    --table-name users \
    --attribute-definitions \
      AttributeName=user_id,AttributeType=S \
    --key-schema \
      AttributeName=user_id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --endpoint-url $ENDPOINT_URL
  echo "Tabla 'users' creada"
else
  echo "La tabla 'users' ya existe"
fi

# Crear temas SNS
echo "Creando temas SNS..."

# Tema tweets
if ! resource_exists "aws sns list-topics --endpoint-url $ENDPOINT_URL | grep tweets"; then
  aws sns create-topic --name tweets --endpoint-url $ENDPOINT_URL
  echo "Tema 'tweets' creado"
else
  echo "El tema 'tweets' ya existe"
fi

# Tema follows
if ! resource_exists "aws sns list-topics --endpoint-url $ENDPOINT_URL | grep follows"; then
  aws sns create-topic --name follows --endpoint-url $ENDPOINT_URL
  echo "Tema 'follows' creado"
else
  echo "El tema 'follows' ya existe"
fi

# Crear colas SQS
echo "Creando colas SQS..."

# Función para crear cola y su DLQ
create_queue_with_dlq() {
  local queue_name=$1
  local dlq_name="${queue_name}-dlq"
  local dlq_url
  
  # Crear DLQ
  if ! resource_exists "aws sqs get-queue-url --queue-name $dlq_name --endpoint-url $ENDPOINT_URL"; then
    aws sqs create-queue --queue-name "$dlq_name" --endpoint-url $ENDPOINT_URL
    echo "Cola DLQ '$dlq_name' creada"
  else
    echo "La cola DLQ '$dlq_name' ya existe"
  fi
  
  # Obtener ARN de DLQ
  dlq_url=$(aws sqs get-queue-url --queue-name "$dlq_name" --endpoint-url $ENDPOINT_URL --output text --query 'QueueUrl' || echo "$ENDPOINT_URL/$AWS_ACCOUNT_ID/$dlq_name")
  dlq_arn=$(aws sqs get-queue-attributes --queue-url "$dlq_url" --attribute-names QueueArn --endpoint-url $ENDPOINT_URL --output text --query 'Attributes.QueueArn' || echo "arn:aws:sqs:$AWS_REGION:$AWS_ACCOUNT_ID:$dlq_name")
  
  # Crear cola principal con redirección a DLQ
  if ! resource_exists "aws sqs get-queue-url --queue-name $queue_name --endpoint-url $ENDPOINT_URL"; then
    aws sqs create-queue \
      --queue-name "$queue_name" \
      --attributes "{\"RedrivePolicy\":\"{\\\"deadLetterTargetArn\\\":\\\"$dlq_arn\\\",\\\"maxReceiveCount\\\":\\\"5\\\"}\"}" \
      --endpoint-url $ENDPOINT_URL
    echo "Cola '$queue_name' creada con redirección a DLQ"
  else
    echo "La cola '$queue_name' ya existe"
  fi
}

# Crear las 5 colas principales y sus DLQs
create_queue_with_dlq "orchestrate-fanout"
create_queue_with_dlq "update-timeline"
create_queue_with_dlq "process-new-follow"
create_queue_with_dlq "populate-cache"
create_queue_with_dlq "rebuild-timeline"

# Suscribir colas a temas SNS
echo "Suscribiendo colas a temas SNS..."

# Función para suscribir cola a tema
subscribe_queue_to_topic() {
  local queue_name=$1
  local topic_name=$2
  local queue_url
  local queue_arn
  local topic_arn
  
  # Obtener ARNs
  queue_url=$(aws sqs get-queue-url --queue-name "$queue_name" --endpoint-url $ENDPOINT_URL --output text --query 'QueueUrl')
  queue_arn=$(aws sqs get-queue-attributes --queue-url "$queue_url" --attribute-names QueueArn --endpoint-url $ENDPOINT_URL --output text --query 'Attributes.QueueArn')
  topic_arn=$(aws sns list-topics --endpoint-url $ENDPOINT_URL --output text --query "Topics[?contains(TopicArn, '$topic_name')].TopicArn")
  
  # Verificar si ya existe una suscripción
  if ! aws sns list-subscriptions --endpoint-url $ENDPOINT_URL | grep -q "$queue_arn"; then
    # Suscribir cola a tema
    aws sns subscribe \
      --topic-arn "$topic_arn" \
      --protocol sqs \
      --notification-endpoint "$queue_arn" \
      --endpoint-url $ENDPOINT_URL
    
    # Configurar política de acceso de la cola para permitir mensajes del tema SNS
    aws sqs set-queue-attributes \
      --queue-url "$queue_url" \
      --attributes "{\"Policy\":\"{\\\"Version\\\":\\\"2012-10-17\\\",\\\"Statement\\\":[{\\\"Effect\\\":\\\"Allow\\\",\\\"Principal\\\":{\\\"Service\\\":\\\"sns.amazonaws.com\\\"},\\\"Action\\\":\\\"sqs:SendMessage\\\",\\\"Resource\\\":\\\"$queue_arn\\\",\\\"Condition\\\":{\\\"ArnEquals\\\":{\\\"aws:SourceArn\\\":\\\"$topic_arn\\\"}}}]}\"}" \
      --endpoint-url $ENDPOINT_URL
    
    echo "Cola '$queue_name' suscrita a tema '$topic_name'"
  else
    echo "La cola '$queue_name' ya está suscrita a un tema SNS"
  fi
}

# Suscripciones específicas
subscribe_queue_to_topic "orchestrate-fanout" "tweets"
subscribe_queue_to_topic "process-new-follow" "follows"

echo "Configuración de servicios AWS locales completada"
