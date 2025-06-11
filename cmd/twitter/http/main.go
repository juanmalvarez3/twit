package main

import (
	"context"
	"fmt"
	dmnfollow "github.com/juanmalvarez3/twit/internal/domains/twitter/follow/domain"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/follow/usecases/createfollow"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/juanmalvarez3/twit/pkg/config"
	"github.com/juanmalvarez3/twit/pkg/logger"

	"github.com/juanmalvarez3/twit/internal/adapters/queue"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/timeline/usecases/gettimeline"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/usecases/createtweet"
	"github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/usecases/gettweet"

	dmntweet "github.com/juanmalvarez3/twit/internal/domains/twitter/tweet/domain"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	appLogger, err := logger.New(cfg.Log.Level, cfg.Log.Environment)
	if err != nil {
		log.Fatalf("Error inicializando logger: %v", err)
	}
	defer appLogger.Sync()

	appLogger.Info("Iniciando servicio HTTP",
		zap.String("environment", cfg.Log.Environment),
		zap.String("logLevel", cfg.Log.Level))

	sqsAdapter, err := queue.NewAdapter(cfg)
	if err != nil {
		appLogger.Fatal("Error inicializando adaptador SQS", zap.Error(err))
	}

	populateTimelineCachePublisher := queue.NewPopulateTimelineCachePublisher(sqsAdapter, cfg.SQS.PopulateCacheQueue, appLogger)
	rebuildTimelinePublisher := queue.NewRebuildTimelinePublisher(sqsAdapter, cfg.SQS.RebuildTimelineQueue, appLogger)

	createTweetUC := createtweet.Provide()
	getTweetUC := gettweet.Provide()
	getTimelineUC := gettimeline.Provide(
		populateTimelineCachePublisher,
		rebuildTimelinePublisher,
		appLogger,
	)
	createFollowUC := createfollow.Provide(appLogger)

	deps := &RouterDependencies{
		CreateTweetUC:  createTweetUC,
		GetTweetUC:     getTweetUC,
		GetTimelineUC:  getTimelineUC,
		CreateFollowUC: createFollowUC,
		Logger:         appLogger,
	}

	router := setupRouter(deps)

	port, _ := strconv.Atoi(cfg.Server.Port)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	go func() {
		appLogger.Info("Servidor HTTP iniciado", zap.String("port", cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Error en servidor HTTP", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	appLogger.Info("Apagando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		appLogger.Fatal("Error en el cierre del servidor", zap.Error(err))
	}

	appLogger.Info("Servidor detenido")
}

type RouterDependencies struct {
	CreateTweetUC  createtweet.UseCase
	GetTweetUC     gettweet.UseCase
	GetTimelineUC  gettimeline.UseCase
	CreateFollowUC createfollow.UseCase
	Logger         logger.LoggerInterface
}

func setupRouter(deps *RouterDependencies) *gin.Engine {
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := engine.Group("/api/v1")
	{
		t := v1.Group("/tweets")
		{
			t.POST("/", func(c *gin.Context) {
				var tweetRequest dmntweet.Tweet
				if err := c.BindJSON(&tweetRequest); err != nil {
					deps.Logger.Error("Error deserializando request", zap.Error(err))
					c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo deserializar el request"})
					return
				}
				tweet, err := deps.CreateTweetUC.CreateTweet(c.Request.Context(), &tweetRequest)
				if err != nil {
					deps.Logger.Error("Error creando tweet", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el tweet"})
					return
				}
				c.JSON(http.StatusCreated, gin.H{tweet.ID: tweet})
			})
			t.GET("/:id", func(c *gin.Context) {
				id := c.Param("id")
				tweet, err := deps.GetTweetUC.GetTweet(c.Request.Context(), id)
				if err != nil {
					deps.Logger.Error("Error obteniendo tweet", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el tweet"})
					return
				}
				c.JSON(http.StatusOK, tweet)
			})
		}

		f := v1.Group("/follows")
		{
			f.POST("/", func(c *gin.Context) {
				var followRequest dmnfollow.Follow
				if err := c.BindJSON(&followRequest); err != nil {
					deps.Logger.Error("Error deserializando request", zap.Error(err))
					c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo deserializar el request"})
					return
				}

				err := deps.CreateFollowUC.CreateFollow(c.Request.Context(), followRequest)
				if err != nil {
					deps.Logger.Error("Error creando follow", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el follow"})
					return
				}
				c.JSON(http.StatusAccepted, gin.H{"message": "Follow creado!"})
			})
		}

		tl := v1.Group("/timeline")
		{
			tl.GET("/:user_id", func(c *gin.Context) {
				userID := c.Param("user_id")
				timeline, err := deps.GetTimelineUC.Exec(c.Request.Context(), userID)
				if err != nil {
					deps.Logger.Error("Error obteniendo timeline", zap.String("user_id", userID), zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el timeline"})
					return
				}
				c.JSON(http.StatusOK, timeline)
			})
		}
	}

	return engine
}
