package user

import (
	"artika/api/data"
	"errors"
	"mime/multipart"
)

// This file is now acting as the logic layer on-top of the the database
// Do not allow database errors to bubble up past these functions

// Todo: Should prob rename to something more generic like "logic.go"
// Also move verification of sessions to here and out of the database

var FailedToDecodeJWTErr = errors.New("Failed to decode JWT")
var FailedToRegisterUser = errors.New("Failed to register user")
var FailedToRegisterSession = errors.New("Failed to register session")
var FailedToDeleteSessionErr = errors.New("Failed to delete session")
var FailedToValidateSessionErr = errors.New("Failed to validate session")
var FailedToFindUserFromSessionIDErr = errors.New("Failed to find user from session ID")
var FailedToEnableUserAsDesignerErr = errors.New("Failed to enable user as designer")
var FailedToUpdateUserInfoErr = errors.New("Failed to update user info")
var FailedToCreatePictureTooLargeErr = errors.New("Failed to create picture, too large")
var FailedToCreatePictureFailedToStoreErr = errors.New("Failed to create picture, failed to store")

var FailedToCreateWorkItemErr = errors.New("Failed to create work item")
var FailedToGetAllWorkItemsErr = errors.New("Failed to get all work items")

var FailedToGetImageNotFoundErr = errors.New("Failed to get image, not found")

var database = data.NewArrayDatabaseConnection()
var pictureStore = data.NewFilesystemPictureStore("./_pictures")

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

func GetUserForValidSessionID(sessionID string) (data.UserInfo, error) {
	userInfo, err := database.GetUserForValidSessionID(sessionID)
	if err != nil {
		return data.UserInfo{}, FailedToFindUserFromSessionIDErr
	}

	return userInfo, nil
}

func GetAllPublicDesigners() ([]data.UserInfoPublic, error) {
	designers, err := database.GetAllPublicDesigners()
	if err != nil {
		return []data.UserInfoPublic{}, err
	}

	return designers, nil
}

func EnableUserAsDesignerForValidSessionID(sessionID string) error {
	err := database.EnableUserAsDesignerForValidSessionID(sessionID)
	if err != nil {
		return FailedToEnableUserAsDesignerErr
	}

	return nil
}

func UpdateUserInfoForValidSessionID(sessionID string, userInfoUpdate data.UserInfoUpdate) error {
	err := database.UpdateUserInfoForValidSessionID(sessionID, userInfoUpdate)
	if err != nil {
		return FailedToUpdateUserInfoErr
	}

	return nil
}

// These functions are less about users and more about the other aspects of the app
// should move, except the database connection exists here...

func CreateWorkItemForValidSessionID(sessionID string, workItemCreateInfo data.WorkItemCreateInfo) (data.WorkItemInfo, error) {
	workItem, err := database.CreateWorkItemForValidSessionID(sessionID, workItemCreateInfo)
	if err != nil {
		return data.WorkItemInfo{}, FailedToCreateWorkItemErr
	}

	return workItem, nil
}

func DeleteWorkItemForValidSessionID(sessionID string, workItemID string) error {
	err := database.DeleteWorkItemForValidSessionID(sessionID, workItemID)
	if err != nil {
		return FailedToCreateWorkItemErr
	}

	return nil
}

func GetAllWorkItemsForValidSessionID(sessionID string) ([]data.WorkItemInfo, error) {
	workItems, err := database.GetAllWorkItemsForValidSessionID(sessionID)
	if err != nil {
		return []data.WorkItemInfo{}, FailedToGetAllWorkItemsErr
	}

	return workItems, nil
}

func StorePictureForValidSessionID(sessionID string, file multipart.File, header *multipart.FileHeader) (data.Picture, error) {
	if header.Size > 10e6 {
		return data.Picture{}, FailedToCreatePictureTooLargeErr
	}

	filename, err := pictureStore.StorePicture(file, header)
	if err != nil {
		return data.Picture{}, FailedToCreatePictureFailedToStoreErr
	}

	var pictureCreateInfo = data.PictureCreateInfo{
		URI: filename,
	}

	picture, err := database.GetOrCreatePictureForValidSessionID(sessionID, pictureCreateInfo)
	if err != nil {
		return data.Picture{}, FailedToCreatePictureFailedToStoreErr
	}

	return picture, nil
}

func GetPicture(pictureID string) (data.PictureFileData, error) {
	picture, err := database.GetPictureFromPictureID(pictureID)
	if err != nil {
		return data.PictureFileData{}, FailedToGetImageNotFoundErr
	}

	pictureData, err := pictureStore.GetPicture(picture.URI)
	if err != nil {
		return data.PictureFileData{}, FailedToGetImageNotFoundErr
	}

	return pictureData, nil
}
