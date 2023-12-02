package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
)

func _download_string(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Error making GET request:", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body:", err)
	}

	return string(body), nil
}

func decode_jwt(jwtb64 string) (jwt.MapClaims, error) {
	json_string, err := _download_string("https://www.googleapis.com/oauth2/v3/certs")
	if err != nil {
		return nil, err
	}

	json := json.RawMessage(json_string)

	jwks, err := keyfunc.NewJSON(json)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(jwtb64, jwks.Keyfunc)
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	return claims, nil
}
