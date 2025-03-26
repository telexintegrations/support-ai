package telexcom

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/telexintegrations/support-ai/format"
	chromadb "github.com/telexintegrations/support-ai/internal/repository/chromaDB"
)

func extractTexts(results []chromadb.SearchResult) string {
	var allTexts []string

	for _, result := range results {
		allTexts = append(allTexts, result.Text...) // Append all text values
	}

	return strings.Join(allTexts, "\n") // Join with newlines (or ", " if preferred)
}

func (txc *TelexCom) ProcessHelpCmd2(ctx context.Context, query, chanID string, orgID string) error {

	res, err := txc.cdb.SearchVectorFromContentEmbedding(ctx, query, 3, orgID)

	if err != nil {
		slog.Error("failed to search for relevant content embeddings")
		lastMessageToTelex = failedToProcessQueryMsg
		txc.SendMessageToTelex(ctx, lastMessageToTelex, chanID)
		return ErrSemanticVectorSearchFailed
	}

	similarTextToQuery := extractTexts(res)
	fmt.Printf("final response from search is %s:", similarTextToQuery)

	lastMessageToTelex, err = txc.aisvc.FineTunedResponse(query, similarTextToQuery)
	fmt.Printf("ai response is %s", lastMessageToTelex)
	if err != nil {
		slog.Error("failed to generate AI response")
		lastMessageToTelex = failedToProcessQueryMsg
		txc.SendMessageToTelex(ctx, lastMessageToTelex, chanID)
		return ErrFailedToGetAIResponse
	}
	txc.SendMessageToTelex(ctx, lastMessageToTelex, chanID)
	return nil
}

func (txc *TelexCom) ProcessUploadCmd2(ctx context.Context, content, channelID, orgID string) error {

	chunks := format.ChunkTextByParagraph(content, 20)

	va := chromadb.ChromaContentEmbeddings{OrgId: orgID, ContentChunks: chunks}
	fmt.Printf("%+v", va)
	cotx := context.TODO()

	err := txc.cdb.InsertIntoChromaEmbeddingCollection(cotx, va)

	// err = txc.db.InsertIntoEmbeddingCollection(ctx, chunks, chunkEmbeddings, dbmodel.OrgMetaData{ID: chanID})
	if err != nil {
		slog.Error("failed to insertIntoDBCollections", "details", err)
		lastMessageToTelex = failedToProcessUploadMsg
		txc.SendMessageToTelex(ctx, lastMessageToTelex, channelID)
		return ErrFailedToUploadContextToAI
	}
	slog.Info("successfull, inserting chunk embeddings into db")
	lastMessageToTelex = successUploadMsg
	txc.SendMessageToTelex(ctx, lastMessageToTelex, channelID)
	return nil
}
