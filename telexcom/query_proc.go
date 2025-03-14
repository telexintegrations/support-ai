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
	} else if strings.HasPrefix(queryToLower, "/help") {
		remainingContent = strings.TrimPrefix(query, "/help")
		task = "/help"
		fmt.Println("Processing help with:", strings.TrimSpace(remainingContent))
	} else if strings.HasPrefix(queryToLower, "/change-contex") {
		remainingContent = strings.TrimPrefix(query, "/change-context")
		task = "/change-context"
		fmt.Println("Processing /change-context:", strings.TrimSpace(remainingContent))
	} else if strings.HasPrefix(queryToLower, "use") {
		task = ""
	} else {
		remainingContent = query
		task = "use"
	}
	return remainingContent, task
}
