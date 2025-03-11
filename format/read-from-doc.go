package format

import (
	"fmt"
	"log"
	"os"

	"github.com/nguyenthenguyen/docx"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

// isPDF checks if the file is a PDF.
func isPDF(filePath string) bool {
	return len(filePath) > 4 && filePath[len(filePath)-4:] == ".pdf"
}

// isDOCX checks if the file is a DOCX.
func isDOCX(filePath string) bool {
	return len(filePath) > 5 && filePath[len(filePath)-5:] == ".docx"
}

func ExtractTextFromDocx(filePath string) (string, error) {
	doc, err := docx.ReadDocxFile(filePath)

	if err != nil {
		return "", err
	}
	defer doc.Close()

	text := doc.Editable().GetContent()
	return text, nil
}

func ExtractTextFromPDF(pdfPath string) (string, error) {
	f, err := os.Open(pdfPath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF file: %w", err)
	}
	defer f.Close()

	// Load the PDF document
	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return "", fmt.Errorf("failed to read PDF: %w", err)
	}

	// Get total pages
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", fmt.Errorf("failed to get page count: %w", err)
	}

	// Extract text from each page
	var extractedText string
	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			log.Printf("Warning: could not get page %d: %v", i, err)
			continue
		}

		ex, err := extractor.New(page)
		if err != nil {
			log.Printf("Warning: could not create extractor for page %d: %v", i, err)
			continue
		}

		text, err := ex.ExtractText()
		if err != nil {
			log.Printf("Warning: could not extract text from page %d: %v", i, err)
			continue
		}

		extractedText += text + "\n"
	}

	return extractedText, nil
}

func ExtractText(path string) (string, error) {
	if isPDF(path) {
		return ExtractTextFromPDF(path)
	} else if isDOCX(path) {
		return ExtractTextFromDocx(path)
	} else {
		return "", fmt.Errorf("unsupported file format: %s", path)
	}
}
