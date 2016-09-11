package imstor

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/denbeigh2000/imstor"
	"github.com/denbeigh2000/imstor/thumbnailer"
	"github.com/denbeigh2000/imstor/thumbnailer/thumbnailers"

	"github.com/denbeigh2000/imstor/validator"
)

// Server is able to translate incoming requests to corresponding calls to the
// Imstor app. It should block until it is no longer serving requests.
type Server interface {
	Serve(UserAPI) error
}

func NewUserAPI(img imstor.Store, thumb imstor.ThumbnailStore) UserAPI {
	thumber := thumbnailers.NewLocalThumbnailer(
		thumbnailer.NewLocal(),
		img, thumb, 50,
	)

	v := validator.NewLocal()
	processor := validator.NewTimedLocalProcessor(v, 2*time.Minute, 150, 2)

	imgAPI := userImageAPI{Store: img, Thumbnailer: thumber, Validator: processor}
	thumbAPI := userThumbnailAPI{ThumbnailStore: thumb}

	return userAPI{
		UserImageAPI:     imgAPI,
		UserThumbnailAPI: thumbAPI,
	}
}

type userAPI struct {
	UserImageAPI
	UserThumbnailAPI
}

// Keeping these separate, because it's totally feasible that they be separated
// in future
type UserAPI interface {
	UserImageAPI
	UserThumbnailAPI
}

type UserImageAPI interface {
	CreateImage(io.Reader) (imstor.Image, error)
	RetrieveImage(imstor.ID) (imstor.Image, error)
	DownloadImage(imstor.ID) (io.ReadCloser, error)
}

type UserThumbnailAPI interface {
	RetrieveThumbnails(imstor.ID) ([]imstor.Thumbnail, error)
	DownloadThumbnail(imstor.Thumbnail) (io.ReadCloser, error)
}

type userImageAPI struct {
	Store       imstor.Store
	Thumbnailer thumbnailer.AsyncThumbnailer

	Validator validator.Processor
}

func (a userImageAPI) validateImage(r io.Reader) (imstor.ImageInfo, error) {
	req := validator.NewRequest(r)

	err := a.Validator.Process(req)
	if err != nil {
		return imstor.ImageInfo{}, err
	}

	log.Println("Waiting for validation response...")
	resp := req.Response()
	log.Printf("Received %v", resp)
	if resp.Err != "" {
		return imstor.ImageInfo{}, fmt.Errorf(resp.Err)
	}

	return resp.ImageInfo, nil
}

func (a userImageAPI) CreateImage(r io.Reader) (imstor.Image, error) {
	imageID, err := a.Store.Create()
	if err != nil {
		return imstor.Image{}, err
	}

	rBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return imstor.Image{}, err
	}

	r = bytes.NewReader(rBytes)
	info, err := a.validateImage(r)
	if err != nil {
		return imstor.Image{}, err
	}

	img := imstor.Image{
		ID: imageID,
		Metadata: imstor.Metadata{
			ImageInfo: info,
		},
	}

	r = bytes.NewReader(rBytes)
	newImg, err := a.Store.Upload(img, r)
	if err != nil {
		return imstor.Image{}, err
	}

	a.Thumbnailer.Queue(thumbnailer.Request{
		ID:   imageID,
		Size: imstor.Size{LongEdge: 300},
	})

	return newImg, nil
}

func (a userImageAPI) RetrieveImage(ID imstor.ID) (imstor.Image, error) {
	return a.Store.Retrieve(ID)
}

func (a userImageAPI) DownloadImage(ID imstor.ID) (io.ReadCloser, error) {
	return a.Store.Download(ID)
}

type userThumbnailAPI struct {
	ThumbnailStore imstor.ThumbnailStore
}

func (a userThumbnailAPI) RetrieveThumbnails(ID imstor.ID) ([]imstor.Thumbnail, error) {
	return a.ThumbnailStore.RetrieveThumbs(ID)
}

func (a userThumbnailAPI) DownloadThumbnail(t imstor.Thumbnail) (io.ReadCloser, error) {
	return a.ThumbnailStore.DownloadThumb(t)
}
