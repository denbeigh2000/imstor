package validator

import (
	"image"
	"io"

	_ "image/jpeg"
	_ "image/png"

	"github.com/denbeigh2000/imstor"
)

type Validator interface {
	Validate(io.Reader) (imstor.ImageInfo, error)
}

func NewLocal() Validator {
	return local{}
}

type local struct{}

func (l local) Validate(data io.Reader) (imstor.ImageInfo, error) {
	config, codecStr, err := image.DecodeConfig(data)
	if err != nil {
		return imstor.ImageInfo{}, err
	}

	size := imstor.FromImageConfig(config)
	codec, err := imstor.GuessCodec(codecStr)
	if err != nil {
		return imstor.ImageInfo{}, err
	}

	return imstor.ImageInfo{
		Codec: codec,
		Size:  size,
	}, nil
}
