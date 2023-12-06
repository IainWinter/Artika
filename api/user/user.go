package user

import (
	"artika/api/data"
	"errors"
)

// This file is now acting as the logic layer on-top of the the database
// Do not allow database errors to bubble up past these functions

var FailedToDecodeJWTErr = errors.New("Failed to decode JWT")
var FailedToRegisterUser = errors.New("Failed to register user")
var FailedToRegisterSession = errors.New("Failed to register session")
var FailedToDeleteSessionErr = errors.New("Failed to delete session")
var FailedToValidateSessionErr = errors.New("Failed to validate session")
var FailedToFindUserFromSessionIDErr = errors.New("Failed to find user from session ID")

var database = data.NewArrayDatabaseConnection()

// Register a jwt with the system.
// If the user doesn't exist, create them.
// If the session doesn't exist, create it.
// If the session does exist, update its expiration time.
// Return the session.
func RegisterJWT(jwt string) (data.UserSession, error) {
	userInfo, err := DecodeJWT(jwt)
	if err != nil {
		return data.UserSession{}, FailedToDecodeJWTErr
	}

	err = database.CreateUserIfNotExists(userInfo)
	if err != nil {
		return data.UserSession{}, FailedToRegisterUser
	}

	userSession, err := database.GetOrCreateSession(userInfo)
	if err != nil {
		return data.UserSession{}, FailedToRegisterSession
	}

	return userSession, nil
}

func DeleteSession(sessionID string) error {
	err := database.DeleteSession(sessionID)
	if err != nil {
		return FailedToDeleteSessionErr
	}

	return nil
}

// Return true if the session exists and is of recent creation.
// Otherwise return false.
func IsSessionValid(sessionID string) (bool, error) {
	isValid, err := database.IsSessionValid(sessionID)

	// If the session is not valid, but the error is not front a database issue
	// Then just return false, do not error
	if err != nil && err != data.SessionExpiredErr && err != data.SessionNotFoundErr {
		return false, FailedToValidateSessionErr
	}

	return isValid, nil
}

func GetUserFromValidSessionID(sessionID string) (data.UserInfo, error) {
	userInfo, err := database.GetUserFromValidSessionID(sessionID)
	if err != nil {
		return data.UserInfo{}, FailedToFindUserFromSessionIDErr
	}

	return userInfo, nil
}
