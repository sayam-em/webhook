package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8" 
)


type WebhookPayload struct {


	Url string `json:"url"`
	WebhookId string `json:"WebhookId"`
	Data struct {
		Id      string `json:"id"`
		Lat     float64 `json:"lat"`
		Long    float64 `json:"long"`
		Created string `json:"created"`
	} `json:"data"`
}

func Subscribe(ctx context.Context, client *redis.Client, webhookQueue chan WebhookPayload) error {  

	pubSub := client.Subscribe(ctx, "coordinates")

	defer func (pubSub *redis.PubSub) {
		if err := pubSub.Close(); err != nil {
			log.Println("Error closing PubSub:",err)
		}
	}(pubSub)

	var payload WebhookPayload 


	for {
      msg, err := pubSub.ReceiveMessage(ctx)  
      if err != nil {  
         return err 
      } 

		err = json.Unmarshal([]byte(msg.Payload), &payload)

		if err != nil {  
			log.Println("Error unmarshalling payload:", err)  
			continue // Continue with the next message if there's an error unmarshalling  
		 }  

		 webhookQueue <- payload
	}







}