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

type Size struct {
	Width  uint
	Height uint
}

func (s Size) Key() string {
	switch {
	case s.Width > 0 && s.Height > 0:
		return fmt.Sprintf("%vx%v", s.Width, s.Height)
	case s.Height > 0:
		return strconv.FormatUint(uint64(s.Height), 10)
	case s.Width > 0:
		return strconv.FormatUint(uint64(s.Width), 10)
	default:
		panic("Can't generate a key of a thumbnail with zero size")
	}
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
