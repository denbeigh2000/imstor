package thumbnailer

import (
	"io"

	"github.com/denbeigh2000/imstor"
)

type Thumbnailer interface {
	Thumbnail(io.Reader, imstor.Size) (io.Reader, error)
}

type AsyncThumbnailer interface {
	Queue(Request)
}

type Request struct {
	imstor.ID
	imstor.Size
}

type Result struct {
	Request
	io.Reader
}
