package data

import (
	"time"

	"github.com/google/uuid"
)

type ArrayDatabaseConnection struct {
	sessions  []UserSession
	users     []UserInfo
	workItems []WorkItemInfo
	pictures  []Picture
}

func NewArrayDatabaseConnection() DatabaseConnectionInterface {
	database := &ArrayDatabaseConnection{
		sessions: []UserSession{},
		users:    []UserInfo{},
	}

	return database
}

func (db *ArrayDatabaseConnection) Connect(connectionString string) error {
	// The array database connection doesn't need to connect to anything
	return nil
}

func (db *ArrayDatabaseConnection) GetOrCreateSession(user UserInfo) (UserSession, error) {
	var currentTime = time.Now().Unix()

	for i, s := range db.sessions {
		if s.UserID == user.UniqueID {
			db.sessions[i].UnixTimeExpires = currentTime + 60*60
			return db.sessions[i], nil
		}
	}

	var session = CreateSessionFromUser(user)
	db.sessions = append(db.sessions, session)

	return session, nil
}

func (db *ArrayDatabaseConnection) DeleteSession(sessionID string) error {
	for i, s := range db.sessions {
		if s.UniqueID == sessionID {
			db.sessions = append(db.sessions[:i], db.sessions[i+1:]...)
			return nil
		}
	}

	return SessionNotFoundErr
}

func (db *ArrayDatabaseConnection) IsSessionValid(sessionID string) (bool, error) {
	var currentTime = time.Now().Unix()

	for _, s := range db.sessions {
		if s.UniqueID == sessionID {
			if s.UnixTimeExpires <= currentTime {
				return false, SessionExpiredErr
			}

			return true, nil
		}
	}

	return false, SessionNotFoundErr
}

func (db *ArrayDatabaseConnection) CreateUserIfNotExists(user UserInfo) error {
	for _, u := range db.users {
		if u.UniqueID == user.UniqueID {
			return nil
		}
	}

	db.users = append(db.users, user)

	return nil
}

func (db *ArrayDatabaseConnection) GetUserForValidSessionID(sessionID string) (UserInfo, error) {
	var sessionIndex = -1
	for i, s := range db.sessions {
		if s.UniqueID == sessionID {
			sessionIndex = i
			break
		}
	}

	if sessionIndex == -1 {
		return UserInfo{}, SessionNotFoundErr
	}

	for i, u := range db.users {
		if u.UniqueID == db.sessions[sessionIndex].UserID {
			return db.users[i], nil
		}
	}

	return UserInfo{}, UserNotFoundErr
}

func (db *ArrayDatabaseConnection) GetAllPublicDesigners() ([]UserInfoPublic, error) {
	var designers []UserInfoPublic

	for _, u := range db.users {
		if u.IsDesigner {
			var designer UserInfoPublic
			designer.UniqueID = u.UniqueID
			designer.FirstName = u.FirstName
			designer.LastName = u.LastName
			designer.PictureURI = u.PictureURI

			designers = append(designers, designer)
		}
	}

	return designers, nil
}

func (db *ArrayDatabaseConnection) EnableUserAsDesignerForValidSessionID(sessionID string) error {
	userInfo, err := db.GetUserForValidSessionID(sessionID)

	if err != nil {
		return err
	}

	for i, u := range db.users {
		if u.UniqueID == userInfo.UniqueID {
			db.users[i].IsDesigner = true
			return nil
		}
	}

	return UserNotFoundErr
}

func (db *ArrayDatabaseConnection) UpdateUserInfoForValidSessionID(sessionID string, userInfoUpdate UserInfoUpdate) error {
	userInfo, err := db.GetUserForValidSessionID(sessionID)

	if err != nil {
		return err
	}

	for i, u := range db.users {
		if u.UniqueID == userInfo.UniqueID {
			db.users[i].FirstName = userInfoUpdate.FirstName
			db.users[i].LastName = userInfoUpdate.LastName
			//db.users[i].PictureURI = userInfoUpdate.PictureURI
			return nil
		}
	}

	return UserNotFoundErr
}

func (db *ArrayDatabaseConnection) CreateWorkItemForValidSessionID(sessionID string, workItemCreateInfo WorkItemCreateInfo) (WorkItemInfo, error) {
	userInfo, err := db.GetUserForValidSessionID(sessionID)

	if err != nil {
		return WorkItemInfo{}, err
	}

	var workItem = WorkItemInfo{
		WorkID:        uuid.New().String(),
		CreatorUserID: userInfo.UniqueID,
		Title:         workItemCreateInfo.Title,
		Description:   workItemCreateInfo.Description,
		PictureIDs:    workItemCreateInfo.PictureIDs,
	}

	db.workItems = append(db.workItems, workItem)

	return workItem, nil
}

func (db *ArrayDatabaseConnection) DeleteWorkItemForValidSessionID(sessionID string, workItemID string) error {
	userInfo, err := db.GetUserForValidSessionID(sessionID)

	if err != nil {
		return err
	}

	for i, w := range db.workItems {
		if w.WorkID == workItemID {
			if w.CreatorUserID != userInfo.UniqueID {
				return WorkItemNotFoundErr
			}

			db.workItems = append(db.workItems[:i], db.workItems[i+1:]...)
			return nil
		}
	}

	return WorkItemNotFoundErr
}

func (db *ArrayDatabaseConnection) GetAllWorkItemsForValidSessionID(sessionID string) ([]WorkItemInfo, error) {
	userInfo, err := db.GetUserForValidSessionID(sessionID)

	if err != nil {
		return []WorkItemInfo{}, err
	}

	var workItems []WorkItemInfo

	for _, w := range db.workItems {
		if w.CreatorUserID == userInfo.UniqueID {
			workItems = append(workItems, w)
		}
	}

	return workItems, nil
}

func (db *ArrayDatabaseConnection) GetOrCreatePictureForValidSessionID(sessionID string, pictureCreateInfo PictureCreateInfo) (Picture, error) {
	userInfo, err := db.GetUserForValidSessionID(sessionID)

	if err != nil {
		return Picture{}, err
	}

	for _, p := range db.pictures {
		if p.Hash == pictureCreateInfo.Hash {
			return p, nil
		}
	}

	var picture = Picture{
		PictureID: uuid.New().String(),
		UserID:    userInfo.UniqueID,
		URI:       pictureCreateInfo.URI,
		Hash:      pictureCreateInfo.Hash,
	}

	db.pictures = append(db.pictures, picture)

	return picture, nil
}

func (db *ArrayDatabaseConnection) GetPictureFromPictureID(pictureID string) (Picture, error) {
	for _, p := range db.pictures {
		if p.PictureID == pictureID {
			return p, nil
		}
	}

	return Picture{}, PictureNotFoundErr
}
