package user

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// This file is concerned with taking a UserInfo struct and creating or deleting a session.

type UserSession struct {
	UniqueID        string `json:"uniqueID"`
	UserID          string `json:"userID"`
	UnixTimeCreated int64  `json:"unixTimeCreated"`
	UnixTimeExpires int64  `json:"unixTimeExpires"`
}

var sessions []UserSession = []UserSession{
	{UniqueID: "123", UserID: "123", UnixTimeCreated: 123, UnixTimeExpires: 123},
}

func IsSessionValid(sessionId string) (bool, error) {
	for _, s := range sessions {
		if s.UniqueID == sessionId {
			return true, nil
		}
	}

	// This is not an error, when using the database connection that could return an error
	return false, nil
}

func CreateSession(user UserInfo) (UserSession, error) {
	var uniqueSessionID = uuid.New().String()
	var time = time.Now().Unix()

	var session = UserSession{
		UniqueID:        uniqueSessionID,
		UserID:          user.UniqueID,
		UnixTimeCreated: time,
		UnixTimeExpires: time + 60*60,
	}

	sessions = append(sessions, session)

	return session, nil
}

func DeleteSession(user UserInfo) error {
	var sessionIndex int = -1
	for i, s := range sessions {
		if s.UserID == user.UniqueID {
			sessionIndex = i
		}
	}

	if sessionIndex == -1 {
		return fmt.Errorf("Session not found")
	}

	sessions = append(sessions[:sessionIndex], sessions[sessionIndex+1:]...)

	return nil
}
