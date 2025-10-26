package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dnk.com/hoc-golang/internal/config"
	"dnk.com/hoc-golang/internal/db"
	"dnk.com/hoc-golang/internal/db/sqlc"
	"dnk.com/hoc-golang/internal/routes"
	"dnk.com/hoc-golang/internal/utils"
	"dnk.com/hoc-golang/internal/validation"
	"dnk.com/hoc-golang/pkg/auth"
	"dnk.com/hoc-golang/pkg/cache"
	"dnk.com/hoc-golang/pkg/logger"
	"dnk.com/hoc-golang/pkg/mail"
	"dnk.com/hoc-golang/pkg/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Module interface {
	Routes() routes.Route
}
type Application struct {
	config  *config.Config
	router  *gin.Engine
	modules []Module
}

type ModuleContext struct {
	DB    sqlc.Querier
	Redis *redis.Client
}

func NewApplication(cfg *config.Config) (*Application, error) {
	if err := validation.InitValidator(); err != nil {

		logger.Log.Fatal().Err(err).Msg("Validator init failed")
		return nil, err
	}

	r := gin.Default()

	if err := db.InitDB(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Database init failed")

	}

	redisClient := config.NewRedisClient()
	cacheRedisService := cache.NewRedisCacheService(redisClient)
	tokenService := auth.NewJWTService(cacheRedisService)
	mailLogger := utils.NewLoggerWithPath("mail.log", "info")
	factory, err := mail.NewProviderFactory(mail.ProviderMailtrap)
	if err != nil {
		mailLogger.Error().Err(err).Msg("Failed to create mail provider factory")
		return nil, err
	}
	mailService, err := mail.NewMailService(cfg, mailLogger, factory)
	if err != nil {
		mailLogger.Error().Err(err).Msg("Failed to initialize mail service")
		return nil, err
	}
	rabbitmqLogger := utils.NewLoggerWithPath("worker.log", "info")
	rabbitmqService, _ := rabbitmq.NewRabbitMQService(
		utils.GetEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"), rabbitmqLogger,
	)

	ctx := &ModuleContext{
		DB:    db.DB,
		Redis: redisClient,
	}

	modules := []Module{
		NewUserModule(ctx),
		NewAuthModule(ctx, tokenService, cacheRedisService, mailService, rabbitmqService),
	}
	routes.RegisterRoutes(r, tokenService, cacheRedisService, getModulRoutes(modules)...)
	return &Application{
		config:  cfg,
		router:  r,
		modules: modules,
	}, nil

}

func (a *Application) Run() error {
	srv := &http.Server{
		Addr:    a.config.ServerAddress,
		Handler: a.router,
	}

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	// syscall.SIGINT -> ctrl + c
	// syscall.SIGTERM -> KIll service
	// syscall.SIGHUp -> Reload service

	go func() {
		logger.Log.Info().Msgf("Server is running at %s", a.config.ServerAddress)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	<-quit
	logger.Log.Warn().Msg("Shutdown signal received ...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Log.Info().Msg("Server exited gracefully")
	return nil
}
func getModulRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}
	return routeList
}
