package app

import (
	v1handler "dnk.com/hoc-golang/internal/handler/v1"
	"dnk.com/hoc-golang/internal/repository"
	"dnk.com/hoc-golang/internal/routes"
	v1routes "dnk.com/hoc-golang/internal/routes/v1"
	v1service "dnk.com/hoc-golang/internal/service/v1"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule(ctx *ModuleContext) *UserModule {
	// Inittialize repository
	userRepo := repository.NewSqlUserRepository(ctx.DB)

	// Inittialize service
	userService := v1service.NewUserService(userRepo, ctx.Redis)

	// Inittialize handler
	userHandler := v1handler.NewUserHandler(userService)

	// Inittialize routes
	userRoutes := v1routes.NewUserRoutes(userHandler)
	return &UserModule{routes: userRoutes}
}

func (m *UserModule) Routes() routes.Route {
	return m.routes
}
