package data

type Picture struct {
	PictureID string
	UserID    string
	URI       string
}

type PictureCreateInfo struct {
	URI string
}

type PictureFileData struct {
	PictureID string
	Data      []byte
}
