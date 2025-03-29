package format

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func HashFile(file multipart.File) (string, error) {
	hasher := sha256.New()

	// Copy file content into the hasher
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	// Convert hash to a hexadecimal string
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func ComputeFileHash(file multipart.File) (string, error) {
	fileHash, hashErr := HashFile(file)
	if hashErr != nil {
		return "", fmt.Errorf("failed to hash file: %v", hashErr)
	}

	// Reset file pointer to the beginning after reading the file to the end
	if seeker, ok := file.(io.Seeker); ok {
		_, err := seeker.Seek(0, io.SeekStart)
		if err != nil {
			return "", fmt.Errorf("failed to seek file: %v", err)
		}
	}

	return fileHash, nil
}

func DownloadFile(url string, filePath string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func CreateMultipartFileHeader(filePath string) (*multipart.FileHeader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	header := &multipart.FileHeader{
		Filename: stat.Name(),
		Size:     stat.Size(),
	}

	return header, nil
}
