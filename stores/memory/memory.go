package memory

import (
	"bytes"
	"io"
	"io/ioutil"
	"sync"

	"github.com/denbeigh2000/imstor"
	"github.com/satori/go.uuid"
)

func NewStore() imstor.Store {
	return &store{
		store: make(map[imstor.ID]storeEntry),
	}
}

type store struct {
	sync.RWMutex
	store map[imstor.ID]storeEntry
}

type storeEntry struct {
	imstor.Image
	Thumbs map[string]thumbEntry
	Data   []byte
}

type thumbEntry struct {
	imstor.Thumbnail
	Data []byte
}

func (s *store) exists(key imstor.ID) bool {
	s.RLock()
	defer s.RUnlock()

	_, ok := s.store[key]
	return ok
}

func (s *store) uploaded(key imstor.ID) bool {
	s.RLock()
	defer s.RUnlock()
	entry, ok := s.store[key]
	if !ok {
		return false
	}

	return entry.Data != nil
}

func (s *store) retrieve(key imstor.ID) (storeEntry, bool) {
	s.RLock()
	defer s.RUnlock()
	entry, ok := s.store[key]
	return entry, ok
}

func (s *store) key() imstor.ID {
	for {
		key := imstor.ID(uuid.NewV4().String())
		if !s.exists(key) {
			return key
		}
	}
}

func (s *store) Create() (imstor.ID, error) {
	key := s.key()

	entry := storeEntry{
		Image:  imstor.NewImage(key),
		Thumbs: make(map[string]thumbEntry),
	}

	s.Lock()
	defer s.Unlock()
	s.store[key] = entry

	return key, nil
}

func (s *store) Upload(key imstor.ID, r io.Reader) (imstor.Image, error) {
	var img imstor.Image
	if !s.exists(key) {
		return img, imstor.KeyNotFoundErr(key)
	}

	if s.uploaded(key) {
		return img, imstor.AlreadyUploadedErr(key)
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return img, err
	}

	s.Lock()
	defer s.Unlock()

	entry := s.store[key]
	entry.Data = data
	s.store[key] = entry
	img = entry.Image

	return img, nil
}

func (s *store) Retrieve(key imstor.ID) (imstor.Image, error) {
	s.RLock()
	defer s.RUnlock()

	entry, ok := s.store[key]
	if !ok {
		return imstor.Image{}, imstor.KeyNotFoundErr(key)
	}

	return entry.Image, nil
}

func (s *store) Download(key imstor.ID) (io.Reader, error) {
	s.RLock()
	defer s.RUnlock()

	entry, ok := s.store[key]
	if !ok {
		return nil, imstor.KeyNotFoundErr(key)
	}

	if entry.Data == nil {
		return nil, imstor.NotUploadedYetErr(key)
	}

	return bytes.NewReader(entry.Data), nil
}

func (s *store) LinkThumb(ID imstor.ID, size imstor.Size, r io.Reader) (t imstor.Thumbnail, err error) {
	if !s.uploaded(ID) {
		err = imstor.KeyNotFoundErr(ID)
		return
	}

	t.Parent = ID
	t.Size = size

	thumbKey := t.Key()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	s.Lock()
	defer s.Unlock()

	entry := s.store[ID]
	_, ok := entry.Thumbs[thumbKey]
	if ok {
		err = imstor.ThumbnailExistsErr(thumbKey)
		return
	}

	entry.Thumbs[thumbKey] = thumbEntry{
		Thumbnail: t,
		Data:      data,
	}

	s.store[t.Parent] = entry

	return t, nil
}

func (s *store) RetrieveThumbs(ID imstor.ID) ([]imstor.Thumbnail, error) {
	entry, ok := s.retrieve(ID)
	if !ok {
		return nil, imstor.KeyNotFoundErr(ID)
	}

	thumbs := entry.Thumbs

	if len(thumbs) == 0 {
		return nil, nil
	}

	thumbsCopy := make([]imstor.Thumbnail, 0, len(thumbs))
	for _, thumb := range thumbs {
		thumbsCopy = append(thumbsCopy, thumb.Thumbnail)
	}

	return thumbsCopy, nil
}

func (s *store) DownloadThumb(t imstor.Thumbnail) (io.Reader, error) {
	entry, ok := s.retrieve(t.Parent)
	if !ok {
		return nil, imstor.KeyNotFoundErr(t.Parent)
	}

	key := t.Key()
	thumb, ok := entry.Thumbs[key]
	if !ok {
		return nil, imstor.NoSuchThumbnailSizeErr(t.Parent)
	}

	// Prevent people from changing this buffer while we stream
	s.RLock()
	defer s.RUnlock()
	return bytes.NewReader(thumb.Data), nil
}
