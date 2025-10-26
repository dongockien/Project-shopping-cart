package app

import (
	v1handler "dnk.com/hoc-golang/internal/handler/v1"
	"dnk.com/hoc-golang/internal/repository"
	"dnk.com/hoc-golang/internal/routes"
	v1routes "dnk.com/hoc-golang/internal/routes/v1"
	v1service "dnk.com/hoc-golang/internal/service/v1"
	"dnk.com/hoc-golang/pkg/auth"
	"dnk.com/hoc-golang/pkg/cache"
	"dnk.com/hoc-golang/pkg/mail"
	"dnk.com/hoc-golang/pkg/rabbitmq"
)

type AuthModule struct {
	routes routes.Route
}

func NewAuthModule(ctx *ModuleContext, tokenService auth.TokenService, cacheService cache.RedisCacheService, mailService mail.EmailProviderService, rabbitmqService rabbitmq.RabbitMQService) *AuthModule {
	// Inittialize repository
	userRepo := repository.NewSqlUserRepository(ctx.DB)

	// Inittialize service
	authService := v1service.NewAuthService(userRepo, tokenService, cacheService, mailService, rabbitmqService)

	// Inittialize handler
	authHandler := v1handler.NewAuthHandler(authService)

	// Inittialize routes
	authRoutes := v1routes.NewAuthRoutes(authHandler)
	return &AuthModule{routes: authRoutes}
}

func (m *AuthModule) Routes() routes.Route {
	return m.routes
}
