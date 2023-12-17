package data

import (
	"time"

	"github.com/google/uuid"
)

type UserInfo struct {
	UniqueID           string `json:"uniqueID"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Email              string `json:"email"`
	IsEmailVerified    bool   `json:"isEmailVerified"`
	PictureURI         string `json:"pictureURI"`
	UnixTimeCreated    int64  `json:"unixTimeCreated"`
	UnixTimeExpires    int64  `json:"unixTimeExpires"`
	HasShippingAddress bool   `json:"HasShippingAddress"`
	IsDesigner         bool   `json:"isDesigner"`
}

type UserInfoPublic struct {
	UniqueID   string `json:"uniqueID"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	PictureURI string `json:"pictureURI"`
}

type UserSession struct {
	UniqueID        string `json:"uniqueID"`
	UserID          string `json:"userID"`
	UnixTimeExpires int64  `json:"unixTimeExpires"`
}

func CreateSessionFromUser(user UserInfo) UserSession {
	var uniqueSessionID = uuid.New().String()
	var time = time.Now().Unix()

	var session = UserSession{
		UniqueID:        uniqueSessionID,
		UserID:          user.UniqueID,
		UnixTimeExpires: time + 60*60,
	}

	return session
}
