package main

import (
	"context"
	"log"

	"net/http"
	_ "net/http/pprof"

	"github.com/DanielTitkov/antibruteforce-microservice/internal/app"
	"github.com/DanielTitkov/antibruteforce-microservice/internal/app/config"
	"github.com/DanielTitkov/antibruteforce-microservice/internal/app/logger"
	"github.com/DanielTitkov/antibruteforce-microservice/internal/bucketstorage"
)

func main() {
	ctx := context.Context(context.Background())

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	l, err := logger.New("./configs/logger.json")
	if err != nil {
		log.Fatalf("setting up logger failed: %v", err)
	}

	c, err := config.NewAppConfig("./configs/app.yaml")
	if err != nil {
		l.Fatalf("setting up config failed: %v", err)
	}

	bs, err := bucketstorage.New([]string{"login", "password", "ip"}, 1000*c.Buckets.Clean)
	if err != nil {
		l.Fatalf("setting up bucket storage failed: %v", err)
	}

	application := app.New(c, l, bs)
	application.Run(ctx)
}
