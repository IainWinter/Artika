package data

type WorkItemInfo struct {
	WorkID        string
	CreatorUserID string
	Title         string
	Description   string
}

type WorkItemPicture struct {
	PictureID string
	WorkID    string
}

type WorkItemCreateInfo struct {
	Title       string
	Description string
	PictureIDs  []string
}
