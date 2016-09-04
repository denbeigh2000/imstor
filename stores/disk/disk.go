package disk

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/denbeigh2000/imstor"

	"github.com/satori/go.uuid"
)

const dataFile = ".imstor.json"

func NewStore(directory string) imstor.Store {
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		panic(err)
	}

	return diskStore{Directory: directory}
}

func touch(paths ...string) error {
	for _, path := range paths {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		f.Close()
	}

	return nil
}

type diskStore struct {
	Directory string
}

func (d diskStore) createPath(id imstor.ID) (string, string) {
	base := filepath.Join(d.Directory, string(id))
	return fmt.Sprintf("%v.blob", base), fmt.Sprintf("%v.meta", base)
}

func (d diskStore) exists(id imstor.ID) bool {
	_, meta := d.createPath(id)

	_, err := os.Stat(meta)
	return !os.IsNotExist(err)
}

func (d diskStore) Create() (imstor.ID, error) {
	id := imstor.ID(uuid.NewV4().String())
	img := imstor.NewImage(id)

	_, path := d.createPath(id)

	f, err := os.Create(path)
	if err != nil {
		return imstor.ID(""), err
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(img)
	if err != nil {
		return imstor.ID(""), err
	}

	log.Printf("%v: Created on disk", id)

	return imstor.ID(id), nil
}

func (d diskStore) Upload(key imstor.ID, data io.Reader) (imstor.Image, error) {
	log.Printf("%v: Uploading to disk", key)
	path, _ := d.createPath(key)
	f, err := os.Create(path)
	if err != nil {
		return imstor.Image{}, err
	}
	defer f.Close()

	n, err := io.Copy(f, data)
	if err != nil {
		return imstor.Image{}, err
	}
	if n == 0 {
		return imstor.Image{}, imstor.EmptyBodyErr{}
	}

	log.Printf("%v: Uploaded to disk", key)

	img := imstor.NewImage(key)
	return img, nil
}

func (d diskStore) Download(key imstor.ID) (io.ReadCloser, error) {
	if !d.exists(key) {
		return nil, imstor.KeyNotFoundErr(key)
	}

	path, _ := d.createPath(key)
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		// Something something race condition
		return nil, imstor.KeyNotFoundErr(key)
	}

	return f, err
}

func (d diskStore) Retrieve(key imstor.ID) (imstor.Image, error) {
	if !d.exists(key) {
		return imstor.Image{}, imstor.KeyNotFoundErr(key)
	}

	_, path := d.createPath(key)

	f, err := os.Open(path)
	if err != nil {
		return imstor.Image{}, err
	}

	img := imstor.Image{}
	err = json.NewDecoder(f).Decode(&img)
	if err != nil {
		return imstor.Image{}, err
	}

	return img, nil
}
