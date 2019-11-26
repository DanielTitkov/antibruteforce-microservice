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
	log.Println("START")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	client := api.NewABServiceClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	for {
		login := randStringRunes(3)
		password := randStringRunes(3)
		ip := "12.13.145.124"

		log.Printf("Making request: %s, %s, %s", login, password, ip)
		attemptResponse, err := client.Attempt(ctx, &api.AttemptRequest{Login: login, Password: password, Ip: ip})
		if err != nil {
			log.Printf("Error occured during attempt request: %v", err)
			break
		}
		log.Println(attemptResponse)
		time.Sleep(1 * time.Millisecond)
	}

	log.Println("EXIT")
}
