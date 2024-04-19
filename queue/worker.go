package queue  

import (  
   "context"  
   "log"  
   "time"  
   "webhook/sender"  
   redisClient "webhook/redis"  
)

func ProcessWebhooks(ctx context.Context, webhookQueue chan redisClient.WebhookPayload) {  
	for payload := range webhookQueue {  
	   go func(p redisClient.WebhookPayload) {  
		  backoffTime := time.Second 
		  maxBackoffTime := time.Hour 
		  retries := 0  
		  maxRetries := 5  
 
		  for {  
			 err := sender.SendWebHook(p.Data, p.Url, p.WebhookId)  
			 if err == nil {  
				break  
			 }  
			 log.Println("Error sending webhook:", err)  
 
			 retries++  
			 if retries >= maxRetries {  
				log.Println("Max retries reached. Giving up on webhook:", p.WebhookId)  
				break  
			 }  
 
			 time.Sleep(backoffTime)  
 
			 backoffTime *= 2  
			 log.Println(backoffTime)  
			 if backoffTime > maxBackoffTime {  
				backoffTime = maxBackoffTime  
			 }  
		  }  
	   }(payload)  
	}  
 }