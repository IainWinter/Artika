package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func writeGoogleFile(file string, data []byte) error {
	return os.WriteFile(file, data, 0644)
}

func fetchGoogleCerts() (string, error) {
	response, err := http.Get("https://www.googleapis.com/oauth2/v3/certs")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to fetch Google certificates, status code: %d", response.StatusCode)
	}

	jsonData, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func getGoogleCerts() (string, error) {
	file := "certs/google_jwt_oauth.json"
	fileInfo, err := os.Stat(file)

	if os.IsNotExist(err) || time.Since(fileInfo.ModTime()) > 24*time.Hour {
		jsonString, err := fetchGoogleCerts()
		if err != nil {
			return "", err
		}

		err = writeGoogleFile(file, []byte(jsonString))
		if err != nil {
			return "", err
		}

		return jsonString, nil
	}

	b, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
