package main

import (
	"context"
	"log"
	"os"

	redisClient "github.com/sayam-em/webhook/redis"  
	"github.com/sayam-em/webhook/queue"  

	"github.com/go-redis/redis/v8"

)


func main () {

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDRESS"),
		Password: "",
		DB: 0,
	})
	
	webhookQueue := make(chan redisClient.WebhookPayload, 100)

	go queue.ProcessWebhooks(ctx, webhookQueue)

	err := redisClient.Subscribe(ctx, client, webhookQueue)

	if err != nil {  
		log.Println("Error:", err)  
	 }  
  
	 select {} 



}