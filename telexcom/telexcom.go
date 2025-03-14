package telexcom

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/telexintegrations/support-ai/aicom"
	"github.com/telexintegrations/support-ai/internal/repository"
)

var (
	ErrChannelIdNotExist          = errors.New("no channel id provided")
	ErrFailedToPostMessageToTelex = errors.New("failed to send message to telex")
)

const (
	telexWebhookBase         = "https://ping.telex.im/v1/webhooks"
	failedToProcessQueryMsg  = "sorry couldn't process your query, try again"
	failedToProcessUploadMsg = "sorry, failed to process your upload, try again"
	successUploadMsg         = "Content Uploaded, you can use /help to send queries"
	caseUpload               = "/upload"
	caseHelp                 = "/help"
	caseChangeContext        = "/change-context"
	caseUse                  = "use"
)

type TelexCom struct {
	aisvc      aicom.AIService
	db         repository.VectorRepo
	httpClient http.Client
}

var lastMessageToTelex string

func NewTelexCom(aiservice aicom.AIService, dbinterface repository.VectorRepo, client http.Client) *TelexCom {
	return &TelexCom{
		aisvc:      aiservice,
		db:         dbinterface,
		httpClient: client,
	}
}

func (txc *TelexCom) SendMessageToTelex(ctx context.Context, messaage, channelID string) error {
	if channelID == "" {
		return ErrChannelIdNotExist
	}

	respPayload := gin.H{
		"message":    messaage,
		"username":   "Support AI",
		"event_name": "Query Support",
		"status":     "success",
	}

	url := fmt.Sprintf("%s/%s", telexWebhookBase, channelID)

	data, err := json.Marshal(respPayload)
	if err != nil {
		fmt.Printf("error marshalling struct: %v\n", err)
	}

	slog.Info("posting message to telex")
	res, err := txc.httpClient.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("error sending POST request to Telex: %v\n", err)
		return ErrFailedToPostMessageToTelex
	}
	slog.Info("posted message to telex successfully")
	defer res.Body.Close()
	return nil
}

func (txc *TelexCom) ProcessTelexInputRequest(ctx context.Context, req TelexChatPayload) error {
	if lastMessageToTelex == req.Message {
		return nil // no need to process message
	}

	p := bluemonday.StrictPolicy()
	userQuery := p.Sanitize(req.Message)
	if lastMessageToTelex == userQuery {
		return nil // no need to process message
	}
	fmt.Println(lastMessageToTelex, req.Message, "this teling message")
	htmlStrippedQuery, task := processQuery(userQuery)

	switch task {
	case caseUpload:
		err := txc.processUploadCmd(ctx, htmlStrippedQuery, req.ChannelID)
		if err != nil {
			return err
		}
	case caseHelp:
		err := txc.processHelpCmd(ctx, htmlStrippedQuery, req.ChannelID)
		if err != nil {
			return err
		}
	case caseUse:
		err := txc.processManualMsg(ctx, req.ChannelID)
		if err != nil {
			return err
		}
	case caseChangeContext:
		err := txc.processChangeContextCmd(ctx, htmlStrippedQuery, req.ChannelID)
		if err != nil {
			return err
		}
	}
	return nil
}
