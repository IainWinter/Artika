package props

import (
	"artika/api/data"
	"artika/api/user"
	"artika/client/template/pages"
	"errors"
)

var FailedToCreateIndexPagePropsErr = errors.New("Failed to create index page props")

type IndexProps struct {
	IsSessionValid bool
	UserInfo       data.UserInfo
}

func GetIndexPagePropsFromSessionID(sessionID string) (pages.IndexProps, error) {
	isSessionValid, err := user.IsSessionValid(sessionID)
	if err != nil {
		return pages.IndexProps{}, FailedToCreateIndexPagePropsErr
	}

	if isSessionValid {
		userInfo, err := user.GetUserFromValidSessionID(sessionID)
		if err != nil {
			return pages.IndexProps{}, FailedToCreateIndexPagePropsErr
		}

		return pages.IndexProps{
			IsSessionValid: isSessionValid,
			UserInfo:       userInfo,
		}, nil
	}

	return pages.IndexProps{}, nil
}
