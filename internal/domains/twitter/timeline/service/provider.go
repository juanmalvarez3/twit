package service

import (
	"fmt"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/repository"
	pkgLogger "github.com/juanmalvarez3/twit/pkg/logger"
)

func Provide() Service {
	log, err := pkgLogger.ProvideError()
	if err != nil {
		fmt.Println(err)
	}

	return New(repository.Provide(), log)
}
