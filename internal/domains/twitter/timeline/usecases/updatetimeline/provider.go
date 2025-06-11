package updatetimeline

import (
	srvtimeline "github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/service"
	"github.com/juanmalvarez3/twit/pkg/config"
	"github.com/juanmalvarez3/twit/pkg/logger"
)

func Provide() *UseCase {
	cfg, err := config.New()
	if err != nil {
		panic("Error cargando configuración: " + err.Error())
	}

	log, err := logger.New(cfg.Log.Level, cfg.Log.Environment)
	if err != nil {
		panic("Error inicializando logger: " + err.Error())
	}

	return New(srvtimeline.Provide(), log)
}
