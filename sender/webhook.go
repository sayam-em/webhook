package sender

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type Payload struct {
	Event string
	Date string
	Id string
	Payment string
}


func SendWebHook(data interface{}, url string, webhookId string) error {
	jsonBytes,err := json.Marshal(data)

	if err != nil {
		return err
	}
	req,err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))


	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp,err := client.Do(req)

	if err != nil {
		return err
	}

	defer func (Body io.ReadCloser)  {

		if err := Body.Close(); err != nil {
			log.Println("Error closing response body:", err)
		}
		
	}(resp.Body)


	status:= "failed"
	if resp.StatusCode == http.StatusOK {
		status = "delivered"
	}

	log.Println((status))

	if status == "failed" {
		return errors.New(status)
	}

	return nil
}