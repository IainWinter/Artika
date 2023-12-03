package data

import "errors"

var SessionExpiredErr = errors.New("Session expired")
var SessionNotFoundErr = errors.New("Session not found")
var UserNotFoundErr = errors.New("User not found")

type DatabaseConnectionInterface interface {
	Connect(connectionString string) error

	// Get or create a session from a UserInfo.
	// If one already exists, return that and update its creation time.
	// Can also return database connection errors.
	GetOrCreateSession(user UserInfo) (UserSession, error)

	// Delete a session.
	// If one doesn't exist, return SessionNotFoundErr.
	// Can also return database connection errors.
	DeleteSession(sessionID string) error

	// Return true if the session exists and is of recent creation.
	// If the session doesn't exist, return SessionNotFoundErr.
	// If the session is too old, return SessionExpiredErr.
	// Can also return database connection errors.
	IsSessionValid(sessionID string) (bool, error)

	// Create a user if one doesn't exist.
	// If one does exist, do nothing.
	// Can also return database connection errors.
	CreateUserIfNotExists(user UserInfo) error

	// Get a user from an already validated sessionID.
	// If the session doesn't exist, return SessionNotFoundErr.
	// If the user doesn't exist, return UserNotFoundErr.
	// Can also return database connection errors.
	GetUserFromValidSessionID(sessionID string) (UserInfo, error)
}
