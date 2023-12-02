package user

// This file is concerned with taking a JWT from oauth2 and decoding it into a UserInfo struct.

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserInfo struct {
	UniqueID        string `json:"uniqueID"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"isEmailVerified"`
	PictureURI      string `json:"pictureURI"`
	UnixTimeCreated int64  `json:"unixTimeCreated"`
	UnixTimeExpires int64  `json:"unixTimeExpires"`
}

func DecodeJWT(jwtb64 string) (UserInfo, error) {
	// Should cache this in a file, only need to update it based on cache headers
	json_string, err := downloadString("https://www.googleapis.com/oauth2/v3/certs")
	if err != nil {
		return UserInfo{}, err
	}

	json := json.RawMessage(json_string)

	jwks, err := keyfunc.NewJSON(json)
	if err != nil {
		return UserInfo{}, err
	}

	token, err := jwt.Parse(jwtb64, jwks.Keyfunc)
	if err != nil {
		return UserInfo{}, err
	}

	claims := token.Claims.(jwt.MapClaims)

	var userInfo = UserInfo{
		UniqueID:        claims["sub"].(string),
		FirstName:       claims["given_name"].(string),
		LastName:        claims["family_name"].(string),
		Email:           claims["email"].(string),
		IsEmailVerified: claims["email_verified"].(bool),
		PictureURI:      claims["picture"].(string),
		UnixTimeCreated: int64(claims["iat"].(float64)),
		UnixTimeExpires: int64(claims["exp"].(float64)),
	}

	return userInfo, nil
}

func downloadString(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Error making GET request: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %w", err)
	}

	return string(body), nil
}
