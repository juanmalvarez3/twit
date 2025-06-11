package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	AWS      AWSConfig
	DynamoDB DynamoDBConfig
	Redis    RedisConfig
	SNS      SNSConfig
	SQS      SQSConfig
	Cache    CacheConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port string
}

type AWSConfig struct {
	Region    string
	Endpoint  string
	AccessKey string
	SecretKey string
}

type DynamoDBConfig struct {
	TweetsTable    string
	UsersTable     string
	FollowsTable   string
	TimelinesTable string
}

type RedisConfig struct {
	Host        string
	Port        string
	Password    string
	DB          int
	TimelineTTL int
}

type SNSConfig struct {
	TweetsTopic  string
	FollowsTopic string
}

type SQSConfig struct {
	OrchestrateQueue     string
	UpdateTimelineQueue  string
	ProcessFollowQueue   string
	PopulateCacheQueue   string
	RebuildTimelineQueue string
}

type CacheConfig struct {
	Enabled bool
	TTL     int
}

type LogConfig struct {
	Level       string
	Environment string
}

func New() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		AWS: AWSConfig{
			Region:    getEnv("AWS_REGION", "us-east-1"),
			Endpoint:  getEnv("AWS_ENDPOINT", "http://localhost:4566"),
			AccessKey: getEnv("AWS_ACCESS_KEY_ID", "test"),
			SecretKey: getEnv("AWS_SECRET_ACCESS_KEY", "test"),
		},
		DynamoDB: DynamoDBConfig{
			TweetsTable:    getEnv("DYNAMODB_TWEETS_TABLE", "tweets"),
			UsersTable:     getEnv("DYNAMODB_USERS_TABLE", "users"),
			FollowsTable:   getEnv("DYNAMODB_FOLLOWS_TABLE", "follows"),
			TimelinesTable: getEnv("DYNAMODB_TIMELINES_TABLE", "timelines"),
		},
		Redis: RedisConfig{
			Host:        getEnv("REDIS_HOST", "localhost"),
			Port:        getEnv("REDIS_PORT", "6379"),
			Password:    getEnv("REDIS_PASSWORD", ""),
			DB:          getEnvAsInt("REDIS_DB", 0),
			TimelineTTL: getEnvAsInt("REDIS_TIMELINE_TTL", 3600*24),
		},
		SNS: SNSConfig{
			TweetsTopic:  getEnv("SNS_TWEETS_TOPIC", "arn:aws:sns:us-east-1:000000000000:tweets"),
			FollowsTopic: getEnv("SNS_FOLLOWS_TOPIC", "arn:aws:sns:us-east-1:000000000000:follows"),
		},
		SQS: SQSConfig{
			OrchestrateQueue:     getEnv("SQS_ORCHESTRATE_FANOUT_QUEUE", "http://localstack:4566/000000000000/orchestrate-fanout"),
			UpdateTimelineQueue:  getEnv("SQS_UPDATE_TIMELINE_QUEUE", "http://localstack:4566/000000000000/update-timeline"),
			ProcessFollowQueue:   getEnv("SQS_PROCESS_NEW_FOLLOW_QUEUE", "http://localstack:4566/000000000000/process-new-follow"),
			PopulateCacheQueue:   getEnv("SQS_POPULATE_CACHE_QUEUE", "http://localstack:4566/000000000000/populate-cache"),
			RebuildTimelineQueue: getEnv("SQS_REBUILD_TIMELINE_QUEUE", "http://localstack:4566/000000000000/rebuild-timeline"),
		},
		Cache: CacheConfig{
			Enabled: getEnvAsBool("CACHE_ENABLED", true),
			TTL:     getEnvAsInt("CACHE_TTL", 3600),
		},
		Log: LogConfig{
			Level:       getEnv("LOG_LEVEL", "info"),
			Environment: getEnv("APP_ENV", "development"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
