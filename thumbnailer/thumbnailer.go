package thumbnailer

import (
	"io"

	"github.com/denbeigh2000/imstor"
)

type Thumbnailer interface {
	Thumbnail(io.Reader, imstor.Size) (io.Reader, error)
}
