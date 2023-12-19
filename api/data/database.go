package data

import "errors"

var SessionExpiredErr = errors.New("Session expired")
var SessionNotFoundErr = errors.New("Session not found")
var UserNotFoundErr = errors.New("User not found")
var InvalidInfoErr = errors.New("Invalid info")

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
	GetUserForValidSessionID(sessionID string) (UserInfo, error)

	// Get all public designers.
	// Can also return database connection errors.
	GetAllPublicDesigners() ([]UserInfoPublic, error)

	// Enable a user as a designer by looking it up by sessionID.
	// If the session doesn't exist, return SessionNotFoundErr.
	// If the user doesn't exist, return UserNotFoundErr.
	// Can also return database connection errors.
	EnableUserAsDesignerForValidSessionID(sessionID string) error

	// Update a user's info from a sessionID.
	// If the session doesn't exist, return SessionNotFoundErr.
	// If the user doesn't exist, return UserNotFoundErr.
	// Can also return database connection errors.
	UpdateUserInfoForValidSessionID(sessionID string, userInfoUpdate UserInfoUpdate) error

	// Create a new work item and return its ID.
	// If the session doesn't exist, return SessionNotFoundErr.
	// If the user doesn't exist, return UserNotFoundErr.
	// If the work item info is invalid, return InvalidInfoErr.
	// Can also return database connection errors.
	CreateWorkItemForValidSessionID(sessionID string, workItemCreateInfo WorkItemCreateInfo) (WorkItemInfo, error)

	// Get all work items belonging to a user.
	// If the session doesn't exist, return SessionNotFoundErr.
	// If the user doesn't exist, return UserNotFoundErr.
	// Can also return database connection errors.
	GetAllWorkItemsForSessionID(sessionID string) ([]WorkItemInfo, error)

	// Register a picture with the database
	// If the user doesn't exist, return UserNotFoundErr.
	// If the picture info is invalid, return InvalidInfoErr.
	// Can also return database connection errors.
	CreatePictureForValidSessionID(sessionID string, pictureCreateInfo PictureCreateInfo) (Picture, error)
}
