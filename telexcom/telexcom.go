package telexcom

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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

func (txc *TelexCom) GenerateResponseToQuery(ctx context.Context, query, channelID string) error {
	// TODO: refactor this function, write tests
	// Call the AI Service
	// Call the Repository ser
	respPayload := TelexResponsePayload{
		Message:   "We are still building... hold on a little",
		EventName: "info",
		Username:  "TelexAI",
		Status:    "info",
	}

	url := fmt.Sprintf("https://ping.telex.im/v1/webhooks/%s", channelID)
	fmt.Println(url)
	data, err := json.Marshal(respPayload)
	if err != nil {
		fmt.Printf("error marshalling struct: %v\n", err)
	}
	reader := bytes.NewReader(data)

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Post(url, "application/json", reader)
	if err != nil {
		fmt.Printf("error sending POST request to Telex: %v\n", err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println("Response from telex is: ", string(body))
	// slog.Info("Sent error log to telex")

	return nil
}
