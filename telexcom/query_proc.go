package telexcom

import (
	"fmt"
	"strings"
)

func processQuery(query string) (string, string) {
	var remainingContent string
	var task string
	queryToLower := strings.ToLower(query)
	if strings.HasPrefix(queryToLower, "/upload") {
		remainingContent = strings.TrimPrefix(query, "/upload")
		task = "/upload"
		fmt.Println("Processing upload with:", strings.TrimSpace(remainingContent))
	} else {
		remainingContent = query
		task = "/help"
	}
	return remainingContent, task
}
