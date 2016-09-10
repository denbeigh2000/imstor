package imstor

import (
	"fmt"
	"image"
	"strconv"
	"strings"
	"time"
)

type Codec int

const (
	Invalid Codec = iota // Nil/default value - reserved for empty error structs etc.
	JPEG
	PNG
)

type UnsupportedCodecErr string

func (err UnsupportedCodecErr) Error() string {
	return fmt.Sprintf("Unsupported codec type %v", string(err))
}

func GuessCodec(in string) (c Codec, err error) {
	switch strings.ToLower(in) {
	case "png":
		c = PNG
	case "jpg", "jpeg":
		c = JPEG
	default:
		err = UnsupportedCodecErr(in)
	}

	return
}

type ID string

type ImageInfo struct {
	Codec
	Size
}

type Metadata struct {
	Added time.Time

	ImageInfo
}

type Size struct {
	LongEdge, Height, Width uint
}

func FromImageConfig(config image.Config) (size Size) {
	size.Width = uint(config.Width)
	size.Height = uint(config.Height)

	if config.Width >= config.Height {
		size.LongEdge = uint(config.Width)
	} else {
		size.LongEdge = uint(config.Height)
	}

	return
}

// Only single-uint key encoding/decoding supported for now
func FromKey(key string) (Size, error) {
	u, err := strconv.ParseUint(key, 10, 64)
	if err != nil {
		return Size{}, err
	}

	return Size{LongEdge: uint(u)}, nil
}

// Only single-uint key encoding/decoding supported for now
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
