package imstor

import (
	"io"

	"github.com/denbeigh2000/imstor"
	"github.com/denbeigh2000/imstor/thumbnailer"
	"github.com/denbeigh2000/imstor/thumbnailer/thumbnailers"
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

	imgAPI := userImageAPI{Store: img, Thumbnailer: thumber}
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
}

func (a userImageAPI) CreateImage(r io.Reader) (imstor.Image, error) {
	imageID, err := a.Store.Create()
	if err != nil {
		return imstor.Image{}, err
	}

	img, err := a.Store.Upload(imageID, r)
	if err != nil {
		return imstor.Image{}, err
	}

	a.Thumbnailer.Queue(thumbnailer.Request{
		ID:   imageID,
		Size: imstor.Size{LongEdge: 300},
	})

	return img, nil
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
