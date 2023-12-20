package data

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"image"
	"image/jpeg"
	"mime/multipart"
	"os"
)

var FailedToStoreImageInvalidFormatErr = errors.New("Failed to store image, invalid format")
var FailedToStoreImageFilesystemErr = errors.New("Failed to store image, filesystem error")
var FailedToGetImageFileNotFound = errors.New("Failed to get image, file not found")

type FilesystemPictureStore struct {
	Directory string
}

func NewFilesystemPictureStore(directory string) PictureStoreInterface {
	store := &FilesystemPictureStore{
		Directory: directory,
	}

	return store
}

func (store *FilesystemPictureStore) StorePicture(file multipart.File, header *multipart.FileHeader) (string, error) {
	// should generate hash to stop duplicates on disk, can use as id
	// this layer doesn't care about max size

	// read data into array

	var data = make([]byte, header.Size)
	file.Read(data)

	// could use simpler hash function as this is just for dup check
	var hash = sha256.Sum256(data)
	var id = hex.EncodeToString(hash[:])

	var filepath = store.Directory + "/" + id + ".jpeg"

	// don't write if already exists
	if _, err := os.Stat(filepath); err == nil {
		return filepath, nil
	}

	image, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", FailedToStoreImageInvalidFormatErr
	}

	// Create a new file
	fileOut, err := os.Create(filepath)
	defer fileOut.Close()
	if err != nil {
		return "", FailedToStoreImageFilesystemErr
	}

	err = jpeg.Encode(fileOut, image, nil)
	if err != nil {
		return "", FailedToStoreImageFilesystemErr
	}

	return filepath, nil
}

func (store *FilesystemPictureStore) GetPicture(pictureURI string) (PictureFileData, error) {
	var filepath = pictureURI

	data, err := os.ReadFile(filepath)
	if err != nil {
		return PictureFileData{}, FailedToGetImageFileNotFound
	}

	return PictureFileData{
		Data: data,
	}, nil
}
