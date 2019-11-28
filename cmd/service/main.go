package main

import (
	"context"
	"log"

	"github.com/DanielTitkov/antibruteforce-microservice/internal/app"
	"github.com/DanielTitkov/antibruteforce-microservice/internal/app/config"
	"github.com/DanielTitkov/antibruteforce-microservice/internal/app/logger"
)

func main() {
	ctx := context.Context(context.Background())

	l, err := logger.New("./configs/logger.json")
	if err != nil {
		log.Fatalf("setting up logger failed: %v", err)
	}

	c, err := config.NewAppConfig("./configs/app.yaml")
	if err != nil {
		l.Fatalf("setting up config failed: %v", err)
	}

	application := app.New(c, l)
	application.Run(ctx)
}
