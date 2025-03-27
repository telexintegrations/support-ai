package format

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
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
