package format

import (
	"encoding/json"
	"regexp"
	"strings"
)

func FormatResponse(data interface{}) ([]byte, error) {
	formattedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}
	return formattedJSON, nil
}

func CleanText(text string) string {
	// Remove XML tags
	text = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(text, " ")

	// Remove non-printable characters
	text = regexp.MustCompile(`[\x00-\x1F\x7F]+`).ReplaceAllString(text, " ")

	// Fix concatenated words (letters, numbers, and uppercase transitions)
	text = regexp.MustCompile(`([a-zA-Z])([0-9])|([0-9])([a-zA-Z])|([a-z])([A-Z])`).ReplaceAllString(text, `$1$3$5 $2$4$6`)

	// Normalize spaces (convert multiple spaces to a single space)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Ensure paragraph spacing (preserve double line breaks)
	text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")

	// Trim leading/trailing spaces
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
