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

var sessions []UserSession = []UserSession{}
var users []UserInfo = []UserInfo{}

func HasUser(userId string) bool {
	for _, u := range users {
		if u.UniqueID == userId {
			return true
		}
	}

	return false
}

func GetSession(sessionId string) (UserSession, error) {
	for _, s := range sessions {
		if s.UniqueID == sessionId {
			return s, nil
		}
	}

	return UserSession{}, fmt.Errorf("Session not found")
}

func GetUser(userId string) (UserInfo, error) {
	for _, u := range users {
		if u.UniqueID == userId {
			return u, nil
		}
	}

	return UserInfo{}, fmt.Errorf("User not found")
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

	if !HasUser(user.UniqueID) {
		users = append(users, user)
	}

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

func GetUserFromSessionID(sessionID string) (*UserInfo, error) {
	var session UserSession
	for _, s := range sessions {
		if s.UniqueID == sessionID {
			session = s
		}
	}

	var user *UserInfo
	for _, u := range users {
		if u.UniqueID == session.UserID {
			user = &u
		}
	}

	if user == nil {
		return nil, fmt.Errorf("User not found")
	}

	return user, nil
}
