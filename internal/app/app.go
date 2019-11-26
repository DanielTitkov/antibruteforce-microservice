package app

import (
	"context"
	"sync"

	"github.com/DanielTitkov/antibruteforce-microservice/internal/app/config"
	"github.com/DanielTitkov/antibruteforce-microservice/internal/app/grpc"
	"go.uber.org/zap"
)

// Application holds app settings and meta
type Application struct {
	config *config.AppConfig
	logger *zap.SugaredLogger
}

func New(
	config *config.AppConfig,
	logger *zap.SugaredLogger,
) *Application {
	return &Application{
		config: config,
		logger: logger,
	}
}

func (app *Application) Run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(1)
	go grpc.New(app.logger, app.config).Start(ctx)

	// select {
	// case <-grpcServerErrCh:
	// 	<-httpServerErrCh
	// case <-httpServerErrCh:
	// 	<-grpcServerErrCh
	// }
	app.logger.Info("App started")
	wg.Wait()
}

func (app *Application) Shutdown() {
	_ = app.logger.Sync()
}
