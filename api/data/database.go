package data

type DatabaseConnectionInterface interface {
	Connect(connectionString string) error

	IsSessionValid(userID string, sessionID string) (bool, error)
	CreateSession(userID string) (string, error)
	DeleteSession(userID string, sessionID string) error
}

type DatabaseConnectionMock struct {
}
