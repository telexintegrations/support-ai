package format

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/nguyenthenguyen/docx"
	"rsc.io/pdf"
)

func ExtractTextFromDocx(fileBytes []byte) (string, error) {
	reader := bytes.NewReader(fileBytes)
	doc, err := docx.ReadDocxFromMemory(reader, int64(len(fileBytes)))
	if err != nil {
		return "", fmt.Errorf("failed to read DOCX file: %w", err)
	}
	defer doc.Close()

	text := doc.Editable().GetContent()
	cleanedText := CleanText(text)

	return cleanedText, nil
}

func ExtractTextFromPDF(pdfBytes []byte) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic in PDF extraction:", r)
		}
	}()

	reader := bytes.NewReader(pdfBytes)
	pdfDoc, err := pdf.NewReader(reader, int64(len(pdfBytes)))
	if err != nil {
		return "", err
	}

	var extractedText strings.Builder
	numPages := pdfDoc.NumPage()
	for i := 1; i <= numPages; i++ {
		page := pdfDoc.Page(i)
		if page.V.IsNull() {
			continue
		}

		// Extract text content from the page
		content := page.Content()
		for _, text := range content.Text {
			extractedText.WriteString(text.S + "\n\n")
		}
	}

	cleanedText := CleanText(extractedText.String())
	return cleanedText, nil
}

func ExtractText(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader.Size == 0 {
		return "", fmt.Errorf("file is empty")
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read first 512 bytes to detect MIME type
	prevBuf := make([]byte, 512)
	_, err = file.Read(prevBuf)
	if err != nil {
		return "", fmt.Errorf("failed to read file header: %w", err)
	}

	mimeType := http.DetectContentType(prevBuf)
	fmt.Println("Detected MIME Type:", mimeType)

	// Reset file pointer
	file.Seek(0, io.SeekStart)

	// Read file into memory
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	fileBytes := buf.Bytes() // Convert to byte slice

	switch mimeType {
	case "application/pdf":
		return ExtractTextFromPDF(fileBytes)
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return ExtractTextFromDocx(fileBytes)
	case "application/zip":
		return ExtractTextFromDocx(fileBytes)
	case "application/msword":
		return ExtractTextFromDocx(fileBytes)
	default:
		return "", fmt.Errorf("unsupported file format: %s", mimeType)
	}
}
