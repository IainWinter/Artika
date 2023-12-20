package data

type Picture struct {
	PictureID string
	UserID    string
	Hash      string
	URI       string
}

type PictureCreateInfo struct {
	URI  string
	Hash string
}

type PictureFileData struct {
	Data []byte
}
