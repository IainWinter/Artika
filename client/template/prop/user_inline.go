package prop

import "artika/api/data"

type UserProps struct {
	Name       string
	PictureURI string
}

func GetUserPropsFromUserInfo(user data.UserInfo) UserProps {
	return UserProps{
		Name:       user.FirstName,
		PictureURI: user.PictureURI,
	}
}
