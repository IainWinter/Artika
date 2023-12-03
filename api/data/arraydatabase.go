package data

import "time"

type ArrayDatabaseConnection struct {
	sessions []UserSession
	users    []UserInfo
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

func (db *ArrayDatabaseConnection) GetUserFromValidSessionID(sessionID string) (UserInfo, error) {
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
