package imstor

import (
	"fmt"
	"strconv"
	"time"
)

type ID string

type Metadata struct {
	Added time.Time
}

func FromKey(key string) (Size, error) {
	u, err := strconv.ParseUint(key, 10, 64)
	if err != nil {
		return Size{}, err
	}

	return Size{LongEdge: uint(u)}, nil
}

type Size struct {
	LongEdge uint
}

func (s Size) Key() string {
	if s.LongEdge <= 0 {
		panic("Need either (width AND height) OR long edge to serialise size")
	}
	return strconv.FormatUint(uint64(s.LongEdge), 10)
}

type Image struct {
	ID
	/* Size // Maybe later */
	Metadata
}

type Thumbnail struct {
	Parent ID
	Size
}

func (t Thumbnail) Key() string {
	sizeStr := t.Size.Key()
	return fmt.Sprintf("%v_%v", t.Parent, sizeStr)
}

func NewImage(ID ID) Image {
	return Image{
		ID: ID,
		Metadata: Metadata{
			Added: time.Now(),
		},
	}
}
