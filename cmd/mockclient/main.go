package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/DanielTitkov/antibruteforce-microservice/api"
	"google.golang.org/grpc"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	// https://stackoverflow.com/questions/7998302/graphing-a-processs-memory-usage
	log.Println("START")

	var maxResTime, minResTime time.Duration
	minResTime = time.Duration(1) * time.Minute

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	client := api.NewABServiceClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), 540*time.Second)
	defer cancel()

	for {
		login := randStringRunes(3)
		password := randStringRunes(3)
		// ip := "12.13.145.124"
		ip := randStringRunes(3)

		log.Printf("Making request: %s, %s, %s", login, password, ip)
		start := time.Now()

		attemptResponse, err := client.Attempt(ctx, &api.AttemptRequest{Login: login, Password: password, Ip: ip})
		if err != nil {
			log.Printf("Error occured during attempt request: %v", err)
			break
		}
		elapsed := time.Since(start)

		if elapsed > maxResTime {
			maxResTime = elapsed
		} else if elapsed < minResTime {
			minResTime = elapsed
		}
		log.Printf("Got response in %v, STATUS: %s, OK: %v", elapsed, attemptResponse.Status, attemptResponse.Ok)
		time.Sleep(100 * time.Microsecond) // 10 000 / second
	}

	log.Println("---------------------------------------")
	log.Printf("Response time stats: min %v, max %v", minResTime, maxResTime)
	log.Println("EXIT")
}
