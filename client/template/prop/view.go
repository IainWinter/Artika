package prop

import (
	"artika/api/data"
	"artika/api/user"
	"errors"
)

var FailedToCreateViewPropsErr = errors.New("Failed to create view props")

type ViewProps struct {
	IsSessionValid bool
	UserInfo       data.UserInfo
}

func GetViewPropsFromSessionID(sessionID string) (ViewProps, error) {
	isSessionValid, err := user.IsSessionValid(sessionID)
	if err != nil {
		return ViewProps{}, FailedToCreateViewPropsErr
	}

	if isSessionValid {
		userInfo, err := user.GetUserFromValidSessionID(sessionID)
		if err != nil {
			return ViewProps{}, FailedToCreateViewPropsErr
		}

		return ViewProps{
			IsSessionValid: isSessionValid,
			UserInfo:       userInfo,
		}, nil
	}

	return ViewProps{}, nil
}
