package telexcom

import (
	"context"

	"github.com/microcosm-cc/bluemonday"
)

func (txc *TelexCom) ProcessTelexChromaInputRequest(ctx context.Context, req TelexChatPayload) error {
	if lastMessageToTelex == req.Message {
		return nil // no need to process message
	}

	p := bluemonday.StrictPolicy()
	userQuery := p.Sanitize(req.Message)
	if lastMessageToTelex == userQuery {
		return nil // no need to process message
	}

	htmlStrippedQuery, task := processQuery(userQuery)

	switch task {
	case caseUpload:
		err := txc.ProcessUploadCmd2(ctx, htmlStrippedQuery, req.ChannelID, req.OrgId)
		if err != nil {
			return err
		}
	case caseHelp:
		err := txc.ProcessHelpCmd2(ctx, userQuery, req.ChannelID, req.OrgId)
		if err != nil {
			return err
		}
	}

	return nil
}
