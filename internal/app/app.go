package app

import (
	"context"
	"os"
	"sync"

	"github.com/DanielTitkov/antibruteforce-microservice/internal/app/config"
	"github.com/DanielTitkov/antibruteforce-microservice/internal/app/grpc"
	"github.com/DanielTitkov/antibruteforce-microservice/internal/bucketstorage"
	"go.uber.org/zap"
)

// Application holds app settings and meta
type Application struct {
	config *config.AppConfig
	logger *zap.SugaredLogger
	bs     *bucketstorage.BucketStorage
}

func New(
	config *config.AppConfig,
	logger *zap.SugaredLogger,
	bs *bucketstorage.BucketStorage,
) *Application {
	return &Application{
		config: config,
		logger: logger,
		bs:     bs,
	}
}

func (app *Application) Run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(1)
	go grpc.New(app.logger, app.config, app.bs).Start(ctx)

	// select {
	// case <-grpcServerErrCh:
	// 	<-httpServerErrCh
	// case <-httpServerErrCh:
	// 	<-grpcServerErrCh
	// }
	app.logger.Infof("App started with pid %v", os.Getpid())
	wg.Wait()
}

func (app *Application) Shutdown() {
	_ = app.logger.Sync()
}
