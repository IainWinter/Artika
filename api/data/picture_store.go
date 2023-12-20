package data

import "mime/multipart"

// prolly should not use mime

type PictureStoreInterface interface {
	// StorePicture stores a picture and returns the URI
	// Can also return filesystem errors.
	StorePicture(file multipart.File, header *multipart.FileHeader) (string, error)

	// GetPicture returns a picture from the store
	// Can also return filesystem errors.
	GetPicture(pictureURI string) (PictureFileData, error)
}
