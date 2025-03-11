package format

import (
	"regexp"
	"strings"
)

func CleanText(text string) string {
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	text = strings.ReplaceAll(text, "\n", " ")

	text = strings.TrimSpace(text)

	return text
}


func ChunkTextByParagraph(text string, maxWords int) []string {
	// ChunkTextByParagraph splits text into meaningful chunks while preserving sentence boundaries.
	sentenceRegex := regexp.MustCompile(`[.?!]\s+`)
	sentences := sentenceRegex.Split(text, -1)

	var chunks []string
	var currentChunk []string
	wordCount := 0

	for _, sentence := range sentences {
		words := strings.Fields(sentence)
		sentenceWordCount := len(words)

		// If adding this sentence exceeds the limit, store the current chunks
		if wordCount+sentenceWordCount > maxWords && len(currentChunk) > 0 {
			chunks = append(chunks, strings.Join(currentChunk, " "))
			currentChunk = []string{} // Reset chunk
			wordCount = 0
		}

		// Add sentence to the chunk
		currentChunk = append(currentChunk, sentence)
		wordCount += sentenceWordCount
	}

	// Append the last chunk if it contains any content
	if len(currentChunk) > 0 {
		chunks = append(chunks, strings.Join(currentChunk, " "))
	}

	return chunks
}