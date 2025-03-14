package format

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/nguyenthenguyen/docx"
	"rsc.io/pdf"
)

// isPDF checks if the file is a PDF.
func isPDF(filePath string) bool {
	return len(filePath) > 4 && filePath[len(filePath)-4:] == ".pdf"
}

// isDOCX checks if the file is a DOCX.
func isDOCX(filePath string) bool {
	return len(filePath) > 5 && filePath[len(filePath)-5:] == ".docx"
}

func ExtractTextFromDocx(fileBytes []byte) (string, error) {
	reader := bytes.NewReader(fileBytes)
	doc, err := docx.ReadDocxFromMemory(reader, int64(len(fileBytes)))
	if err != nil {
		return "", fmt.Errorf("failed to read DOCX file: %w", err)
	}
	defer doc.Close()

	text := doc.Editable().GetContent()
	return text, nil
}

func ExtractTextFromPDF(pdfBytes []byte) (string, error) {
	reader := bytes.NewReader(pdfBytes)
	pdfDoc, err := pdf.NewReader(reader, int64(len(pdfBytes)))
	if err != nil {
		return "", err
	}

	var extractedText string
	numPages := pdfDoc.NumPage()
	for i := 1; i <= numPages; i++ {
		page := pdfDoc.Page(i)
		if page.V.IsNull() {
			continue
		}

		// Extract text content from the page
		content := page.Content()
		for _, text := range content.Text {
			extractedText += text.S
		}
	}
	return extractedText, nil
}

func ExtractText(fileHeader *multipart.FileHeader) (string, error) {
	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read file into memory
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	fileBytes := buf.Bytes() // Convert to byte slice

	// Determine file type and extract text
	if isPDF(fileHeader.Filename) {
		return ExtractTextFromPDF(fileBytes)
	} else if isDOCX(fileHeader.Filename) {
		return ExtractTextFromDocx(fileBytes)
	} else {
		return "", fmt.Errorf("unsupported file format: %s", fileHeader.Filename)
	}
}
