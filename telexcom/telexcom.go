package telexcom

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	BaseTelexAPIProduction = "https://api.telex.im/api/v1/"
	BaseTelexAPIStaging    = "https://api.staging.telex.im/api/v1/"
)

type TelexCom struct {
	// Have the AI interface,
	// Have the Repository interface
}

func NewTelexCom() *TelexCom {
	return &TelexCom{}
}

func (txc *TelexCom) GenerateResponseToQuery(ctx context.Context, response, channelID string) error {
	// TODO: refactor this function, write tests
	fmt.Println("Generating response to telex")

	if channelID == ""{
		channelID = "01956858-b67d-7654-80cd-68abd43b57f2"
	}

	respPayload := gin.H{
		"message": response,
		"username": "Support AI",
        "event_name": "Query Support",
        "status": "success",
	}

	url := fmt.Sprintf("https://ping.telex.im/v1/webhooks/%s", channelID)
	// fmt.Println(url)

	data, err := json.Marshal(respPayload)
	if err != nil {
		fmt.Printf("error marshalling struct: %v\n", err)
	}
	// reader := bytes.NewReader(data)
	fmt.Println(string(data))
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	fmt.Println("Posting to telex")
	
	res, err := client.Post(url, "application/json", bytes.NewBuffer(data))
	fmt.Println("Posted to telex")
	if err != nil {
		fmt.Printf("error sending POST request to Telex: %v\n", err)
	}
	defer res.Body.Close()

	// body, _ := io.ReadAll(res.Body)
	// fmt.Println("Response from telex is: ", string(body))
	// slog.Info("Sent error log to telex")

	return nil
}
