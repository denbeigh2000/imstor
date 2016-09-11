package thumbnailer

import (
	"fmt"
	"image"
	"io"

	"github.com/nfnt/resize"

	"github.com/denbeigh2000/imstor"
)

func NewLocal() Thumbnailer {
	return Local{
		pngEncoder:  pngEncoder{},
		jpegEncoder: jpegEncoder{},
	}
}

type Local struct {
	InterpFn resize.InterpolationFunction

	// Encoders (at least the jpeg encoder) have a tendency to leak,
	// at least the JPEG encoder certainly does.
	pngEncoder
	jpegEncoder
}

func (l Local) decode(src io.Reader) (image.Image, string, error) {
	return image.Decode(src)
}

func (l Local) encode(src image.Image, codec string) (io.Reader, error) {
	realCodec, err := imstor.GuessCodec(codec)
	if err != nil {
		return nil, err
	}

	var e encoder
	switch realCodec {
	case imstor.PNG:
		e = l.pngEncoder
	case imstor.JPEG:
		e = l.jpegEncoder
	default:
		return nil, fmt.Errorf("Unsupported format %v", codec)
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
