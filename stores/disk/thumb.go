package disk

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/denbeigh2000/imstor"
)

type diskThumbStore struct {
	Directory string
}

func NewThumbStore(directory string) imstor.ThumbnailStore {
	return diskThumbStore{Directory: directory}
}

func (d diskThumbStore) createParentDir(ID imstor.ID) string {
	return filepath.Join(d.Directory, string(ID))
}

func (d diskThumbStore) createPath(t imstor.Thumbnail) string {
	return filepath.Join(d.createParentDir(t.Parent), t.Size.Key())
}

func (d diskThumbStore) init() {
	os.MkdirAll(d.Directory, 0755)
}

func (d diskThumbStore) exists(t imstor.Thumbnail) bool {
	_, err := os.Stat(d.createPath(t))
	return !os.IsNotExist(err)
}

func (d diskThumbStore) LinkThumb(ID imstor.ID, size imstor.Size, r io.Reader) (imstor.Thumbnail, error) {
	t := imstor.Thumbnail{Parent: ID, Size: size}
	if d.exists(t) {
		return imstor.Thumbnail{}, imstor.ThumbnailExistsErr(ID)
	}

	os.MkdirAll(d.createParentDir(t.Parent), 0755)

	f, err := os.Create(d.createPath(t))
	if err != nil {
		return imstor.Thumbnail{}, err
	}
	defer f.Close()

	if r == nil {
		return imstor.Thumbnail{}, imstor.EmptyBodyErr{}
	}

	n, err := io.Copy(f, r)
	if n == 0 {
		return imstor.Thumbnail{}, imstor.EmptyBodyErr{}
	}
	if err != nil {
		return imstor.Thumbnail{}, err
	}

	return t, nil
}

func (d diskThumbStore) DownloadThumb(t imstor.Thumbnail) (io.ReadCloser, error) {
	if !d.exists(t) {
		return nil, imstor.NoSuchThumbnailSizeErr(t.Parent)
	}

	f, err := os.Open(d.createPath(t))
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (d diskThumbStore) RetrieveThumbs(ID imstor.ID) ([]imstor.Thumbnail, error) {
	fileInfos, err := ioutil.ReadDir(d.createParentDir(ID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	var thumbs []imstor.Thumbnail
	for _, file := range fileInfos {
		key, err := imstor.FromKey(file.Name())
		if err != nil {
			// Probably file in the wrong place/malformed
			continue
		}

		thumbs = append(thumbs, imstor.Thumbnail{
			Parent: ID,
			Size:   key,
		})
	}

	return thumbs, nil
}
