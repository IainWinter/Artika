package prop

import (
	"artika/api/data"
	"artika/api/user"
	"errors"
)

var FailedToCreateIndexPagePropsErr = errors.New("Failed to create index page props")

type IndexProps struct {
	IsSessionValid bool
	UserInfo       data.UserInfo
}

func GetIndexPagePropsFromSessionID(sessionID string) (IndexProps, error) {
	isSessionValid, err := user.IsSessionValid(sessionID)
	if err != nil {
		return IndexProps{}, FailedToCreateIndexPagePropsErr
	}

	if isSessionValid {
		userInfo, err := user.GetUserFromValidSessionID(sessionID)
		if err != nil {
			return IndexProps{}, FailedToCreateIndexPagePropsErr
		}

		return IndexProps{
			IsSessionValid: isSessionValid,
			UserInfo:       userInfo,
		}, nil
	}

	return IndexProps{}, nil
}
