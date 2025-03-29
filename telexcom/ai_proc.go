package telexcom

import (
	"context"
	"errors"
	"log/slog"

	"github.com/telexintegrations/support-ai/format"
	"github.com/telexintegrations/support-ai/internal/repository/dbmodel"
)

var FirstMessageToTelex = "Hello! I'm your virtual support assistant. I'm here to help with any questions or issues you may have. How can I assist you today?"

var (
	ErrFailedToEmbedQuery         = errors.New("failed to get query vector embeddings")
	ErrSemanticVectorSearchFailed = errors.New("failed to search of relevant content embeddings")
	ErrFailedToGetAIResponse      = errors.New("failed to get ai response")
	ErrFailedToUploadContextToAI  = errors.New("failed to upload content to ai")
)

func (txc *TelexCom) processHelpCmd(ctx context.Context, query, chanID string) error {
	query_embedding, err := txc.aisvc.GetGeminiEmbedding(query)
	if err != nil {
		slog.Error("failed to process query, error getting vectors")
		lastMessageToTelex = failedToProcessQueryMsg
		txc.SendMessageToTelex(ctx, lastMessageToTelex, chanID)
		return ErrFailedToEmbedQuery
	}

	similarChunksForQuery, err := txc.db.SearchVectorFromContentEmbedding(ctx, query_embedding, dbmodel.OrgMetaData{
		ID: chanID,
	}, 10)

	if err != nil {
		slog.Error("failed to search for relevant content embeddings")
		lastMessageToTelex = failedToProcessQueryMsg
		txc.SendMessageToTelex(ctx, lastMessageToTelex, chanID)
		return ErrSemanticVectorSearchFailed
	}

	similarTextToQuery := ""
	for _, chunk := range similarChunksForQuery {
		similarTextToQuery += chunk.Content
	}
	slog.Info("Vector search completed")

	lastMessageToTelex, err = txc.aisvc.FineTunedResponse(query, similarTextToQuery)
	if err != nil {
		slog.Error("failed to generate AI response")
		lastMessageToTelex = failedToProcessQueryMsg
		txc.SendMessageToTelex(ctx, lastMessageToTelex, chanID)
		return ErrFailedToGetAIResponse
	}
	txc.SendMessageToTelex(ctx, lastMessageToTelex, chanID)
	return nil
}

func (txc *TelexCom) processUploadCmd(ctx context.Context, content, chanID string) error {
	chunks := format.ChunkTextByParagraph(content, 20)
	chunkEmbeddings, err := txc.getChunkEmbeddings(ctx, chunks, chanID)
	if err != nil {
		return err
	}

	err = txc.db.InsertIntoEmbeddingCollection(ctx, chunks, chunkEmbeddings, dbmodel.OrgMetaData{ID: chanID})
	if err != nil {
		slog.Error("failed to insertIntoDBCollections", "details", err)
		lastMessageToTelex = failedToProcessUploadMsg
		txc.SendMessageToTelex(ctx, lastMessageToTelex, chanID)
		return ErrFailedToUploadContextToAI
	}
	slog.Info("successful, inserting chunk embeddings into db")
	lastMessageToTelex = successUploadMsg
	txc.SendMessageToTelex(ctx, lastMessageToTelex, chanID)
	return nil
}

func (txc *TelexCom) processChangeContextCmd(ctx context.Context, content, chanID string) error {
	chunks := format.ChunkTextByParagraph(content, 20)
	chunkEmbeddings, err := txc.getChunkEmbeddings(ctx, chunks, chanID)
	if err != nil {
		return err
	}

	err = txc.db.ReplaceEmbeddingContextTxn(ctx, chunks, chunkEmbeddings, dbmodel.OrgMetaData{
		ID: chanID,
	})
	if err != nil {
		slog.Error("failed to change context for ai knowledge", "details", err)
		lastMessageToTelex = failedToProcessUploadMsg
		txc.SendMessageToTelex(ctx, lastMessageToTelex, chanID)
		return ErrFailedToUploadContextToAI
	}
	slog.Info("successfull, inserting chunk embeddings into db")
	lastMessageToTelex = successUploadMsg
	txc.SendMessageToTelex(ctx, lastMessageToTelex, chanID)
	return nil
}

func (txc *TelexCom) processManualMsg(ctx context.Context, chanId string) error {
	manual := `Use:\n/upload followed by your text to upload context
/help followed by a query to interact with AI
/change-context followed by text/document to change the AI's Knowledge base`
	err := txc.SendMessageToTelex(ctx, manual, chanId)
	if err != nil {
		slog.Error("failed to send message to telex", "details", err)
		return ErrFailedToPostMessageToTelex
	}
	slog.Info("success sending manual to telex")
	return nil
}

func (txc *TelexCom) getChunkEmbeddings(ctx context.Context, chunks []string, chanID string) ([][]float32, error) {
	chunkEmbeddings := [][]float32{}
	for _, chunk := range chunks {
		chunkEmbedding, err := txc.aisvc.GetGeminiEmbedding(chunk)
		if err != nil {
			slog.Error("failed to get chunkEmbeddings", "details", err)
			lastMessageToTelex = failedToProcessUploadMsg
			txc.SendMessageToTelex(ctx, lastMessageToTelex, chanID)
			return chunkEmbeddings, ErrFailedToUploadContextToAI
		}
		chunkEmbeddings = append(chunkEmbeddings, chunkEmbedding)
	}
	return chunkEmbeddings, nil
}

func (txc *TelexCom) SendFirstMessageToChannel(ctx context.Context, chanID string) error {
	err := txc.SendMessageToTelex(ctx, FirstMessageToTelex, chanID)
	if err != nil {
		return err
	}
	return nil
}
