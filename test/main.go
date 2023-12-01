package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

const googleTokenInfoURL = "https://www.googleapis.com/oauth2/v3/certs"

func main() {
	// Replace YOUR_GOOGLE_TOKEN with the actual Google JWT token you want to decode.
	tokenString := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImU0YWRmYjQzNmI5ZTE5N2UyZTExMDZhZjJjODQyMjg0ZTQ5ODZhZmYiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiI2MTAwMTc1MTA2ODMtc2JzYzRiNTViOWxkbnJvamFkZTgwY3IzdmJmMnVra3YuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI2MTAwMTc1MTA2ODMtc2JzYzRiNTViOWxkbnJvamFkZTgwY3IzdmJmMnVra3YuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMTMyMjkzNjQxMjYxNjM1MTg5NTQiLCJlbWFpbCI6ImlhaW53aW50ZXIxQGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYmYiOjE3MDE0MDA5MzYsIm5hbWUiOiJJYWluIFdpbnRlciIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS9BQ2c4b2NJcE8wLW54RGNGQ1B5dDZmTGtsSnV5X0kwWk54bEZuWm1zS1BsbjI2dWh0bTg9czk2LWMiLCJnaXZlbl9uYW1lIjoiSWFpbiIsImZhbWlseV9uYW1lIjoiV2ludGVyIiwibG9jYWxlIjoiZW4iLCJpYXQiOjE3MDE0MDEyMzYsImV4cCI6MTcwMTQwNDgzNiwianRpIjoiNTg4MGM1NWU3MTZhOGJiOTAyNTMzNzQ5NjFhZjI3MDQ3MjY0ZmZhZSJ9.Zl8zzGabwyffKjc0XzFf5hbHW5MjJ7XWQAJtlLHa639U2jVqCYOfCgwzzgdePEtAgZMslCqePyC9uOYkyUmObKAeHXgd09G54JPnVKRX5IdHJym7Iqb4mvPE_QRA1XlmA1Pcg-0dqL9SeOLTPcpk5DWaqNyh2GWzW5gxqOlFfphhsLPxDog8Vd-9fIPyBhKlPTUtEy2Gai-CvegNVy5JFY7wo7Q58gEZfyXuUCbe7rYOJ6JHQi8Y5vkEQFyfbpsFs_7kXVUDbDkiL5RU0axDAc3onnRBNoQ_H0mHUJO01Lo6kcBeu09SUTCvrHSVlDGY9rO905eLNoiw-CcEQE7uxA"
	// Decode the Google JWT token.
	claims, err := decodeGoogleToken(tokenString)
	if err != nil {
		log.Fatal(err)
	}

	// Print the decoded claims.
	fmt.Println("Decoded Claims:")
	printClaims(claims)
}

func decodeGoogleToken(tokenString string) (jwt.MapClaims, error) {
	// Fetch Google's public keys for token verification.
	keys, err := fetchGooglePublicKeys()
	if err != nil {
		return nil, err
	}

	// Parse the token without verification to get the key ID.
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	// Find the appropriate key for verification.
	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("Key ID not found in token header")
	}

	key, exists := keys[keyID]
	if !exists {
		return nil, fmt.Errorf("Key not found for Key ID: %s", keyID)
	}

	// Parse and verify the token using the key.
	parsedToken, err := jwt.Parse(tokenString, func(parsedToken *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid.
	if !parsedToken.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	// Extract claims from the parsed token.
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Invalid claims format")
	}

	return claims, nil
}

func fetchGooglePublicKeys() (map[string]interface{}, error) {
	// Fetch Google's public keys used for token verification.
	resp, err := http.Get(googleTokenInfoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to fetch Google public keys: %s", resp.Status)
	}

	// Read and parse the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response.
	var keys map[string]interface{}
	err = json.Unmarshal(body, &keys)
	if err != nil {
		return nil, err
	}

	// Map the kid to the public key
	keysMap := make(map[string]interface{})
	for item := range keys["keys"].([]interface{}) {
		var kid = keys["keys"].([]interface{})[item].(map[string]interface{})["kid"].(string) // Go may actually stink
		var key = keys["keys"].([]interface{})[item].(map[string]interface{})["n"].(string)

		keysMap[kid] = key
	}

	return keysMap, nil
}

func printClaims(claims jwt.MapClaims) {
	for key, value := range claims {
		fmt.Printf("%s: %v\n", key, value)
	}
}
