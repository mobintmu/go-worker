package app

import (
	"go-worker/internal/config"
	"go-worker/internal/health"
	"go-worker/internal/pkg/logger"
	productController "go-worker/internal/product/controller"
	productService "go-worker/internal/product/service"
	"go-worker/internal/server"
	"go-worker/internal/storage/cache"
	"go-worker/internal/storage/sql"
	"go-worker/internal/storage/sql/migrate"
	"go-worker/internal/storage/sql/sqlc"

	"go.uber.org/fx"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(
			logger.NewLogger,
			config.NewConfig,
			sql.InitialDB,
			//server
			health.New,
			server.NewGinEngine,
			server.CreateHTTPServer,
			server.CreateGRPCServer,
			//db
			migrate.NewRunner, // migration runner
			sqlc.New,
			//cache
			cache.NewClient,
			cache.NewCacheStore,
			//controller
			productController.NewAdmin,
			productController.NewClient,
			productController.NewGRPC,
			//service
			productService.New,
		),
		fx.Invoke(
			server.RegisterRoutes,
			server.StartHTTPServer,
			server.StartGRPCServer,
			//migration
			migrate.RunMigrations,
			//life cycle
			logger.RegisterLoggerLifecycle,
			server.GRPCLifeCycle,
		),
	)
}
