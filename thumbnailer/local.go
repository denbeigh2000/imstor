package thumbnailer

import (
	"image"
	"io"

	"github.com/nfnt/resize"

	"github.com/denbeigh2000/imstor"
)

func NewLocal() Thumbnailer {
	return Local{}
}

type Local struct {
	InterpFn resize.InterpolationFunction
}

func (l Local) decode(src io.Reader) (image.Image, string, error) {
	return image.Decode(src)
}

func (l Local) encode(src image.Image, codec string) (io.Reader, error) {
	var e encoder
	switch codec {
	case "png":
		e = pngEncoder{}
	case "jpg", "jpeg":
		e = jpegEncoder{}
	}

	return e.Encode(src)
}

func (l Local) resize(src image.Image, size imstor.Size) (image.Image, error) {
	var height, width uint
	b := src.Bounds()
	srcHeight, srcWidth := b.Max.Y, b.Max.X
	if srcHeight > srcWidth {
		height = size.LongEdge
	} else {
		width = size.LongEdge
	}

	return resize.Resize(width, height, src, l.InterpFn), nil
}

func (l Local) Thumbnail(src io.Reader, size imstor.Size) (io.Reader, error) {
	img, codec, err := l.decode(src)
	if err != nil {
		return nil, err
	}

	thumb, err := l.resize(img, size)
	if err != nil {
		return nil, err
	}

	out, err := l.encode(thumb, codec)
	if err != nil {
		return nil, err
	}

	return out, nil
}
