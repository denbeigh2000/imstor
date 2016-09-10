package imstor

import (
	"fmt"
	"io"
)

type KeyNotFoundErr ID

func (e KeyNotFoundErr) Error() string {
	return fmt.Sprintf("Key not found: %v", string(e))
}

type NotUploadedYetErr ID

func (e NotUploadedYetErr) Error() string {
	return fmt.Sprintf("Image has been created but not uploaded: %v", string(e))
}

type EmptyBodyErr struct{}

func (e EmptyBodyErr) Error() string {
	return fmt.Sprintf("Cannot process image with empty body")
}

type AlreadyUploadedErr ID

func (e AlreadyUploadedErr) Error() string {
	return fmt.Sprintf("Image has already been uploaded: %v", string(e))
}

type ThumbnailExistsErr ID

func (e ThumbnailExistsErr) Error() string {
	return fmt.Sprintf("Thumbnail already exists: %v", string(e))
}

type NoSuchThumbnailSizeErr ID

func (e NoSuchThumbnailSizeErr) Error() string {
	return fmt.Sprintf("No thumbnail with that size for key: %v", string(e))
}

type ImageSink interface {
	Create() (ID, error)
	Upload(Image, io.Reader) (Image, error)
}

type ImageSource interface {
	Retrieve(ID) (Image, error)
	Download(ID) (io.ReadCloser, error)
}

type Store interface {
	ImageSource
	ImageSink
}

type ThumbnailSource interface {
	// Retrieve the metadata about this thumbnail
	RetrieveThumbs(ID) ([]Thumbnail, error)
	// Download a thumbnail matching the given size. Returns a not found error
	// if no such sized thumbnail exists
	DownloadThumb(Thumbnail) (io.ReadCloser, error)
}

type ThumbnailSink interface {
	// Upload a thumbnail, and link it to the image with the given ID
	LinkThumb(ID, Size, io.Reader) (Thumbnail, error)
}

type ThumbnailStore interface {
	ThumbnailSource
	ThumbnailSink
}
