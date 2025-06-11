package repository

import (
	"context"
	"fmt"

	"github.com/juanmalvarez3/twit/pkg/dynamodb"
	pkgLogger "github.com/juanmalvarez3/twit/pkg/logger"
	pkgRedis "github.com/juanmalvarez3/twit/pkg/redis"
)

func Provide() *TimelineRepository {
	dynamo, err := dynamodb.Provide(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	log, err := pkgLogger.ProvideError()
	if err != nil {
		fmt.Println(err)
	}

	return NewTimelineRepository(dynamo, pkgRedis.Provide(), "timelines", log)
}
