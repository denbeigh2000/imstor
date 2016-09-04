package thumbnailer

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

type encoder interface {
	Encode(image.Image) (io.Reader, error)
}

type jpegEncoder struct{}

func (j jpegEncoder) Encode(i image.Image) (io.Reader, error) {
	buf := &bytes.Buffer{}
	err := jpeg.Encode(buf, i, nil)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

type pngEncoder struct{}

func (p pngEncoder) Encode(i image.Image) (io.Reader, error) {
	buf := &bytes.Buffer{}
	err := png.Encode(buf, i)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
