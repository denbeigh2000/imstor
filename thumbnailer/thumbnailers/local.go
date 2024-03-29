package thumbnailers

import (
	"io"
	"log"
	"sync"

	"github.com/denbeigh2000/imstor"
	"github.com/denbeigh2000/imstor/thumbnailer"
)

const DefaultConcurrency = 4

func NewLocalThumbnailer(t thumbnailer.Thumbnailer, source imstor.ImageSource, sink imstor.ThumbnailSink, bufferSize int) thumbnailer.AsyncThumbnailer {
	thumber := local{
		Thumbnailer: t,
		BufferSize:  bufferSize,

		ImageSource:   source,
		ThumbnailSink: sink,

		in:   make(chan thumbnailer.Request, bufferSize),
		out:  make(chan thumbnailer.Result, bufferSize),
		errs: make(chan error),
	}

	go thumber.handleErrs()
	thumber.loop()
	return thumber
}

type local struct {
	thumbnailer.Thumbnailer
	BufferSize  int
	Concurrency int

	imstor.ImageSource
	imstor.ThumbnailSink

	in   chan thumbnailer.Request
	out  chan thumbnailer.Result
	errs chan error
}

func (l local) Queue(req thumbnailer.Request) {
	l.in <- req
}

func (l local) extractThumbnail(r io.ReadCloser, size imstor.Size) (io.Reader, error) {
	defer r.Close()

	result, err := l.Thumbnail(r, size)
	return result, err
}

func (l local) loop() {
	concurrency := l.Concurrency
	if concurrency == 0 {
		concurrency = DefaultConcurrency
	}

	wg := &sync.WaitGroup{}
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(i int) {
			log.Printf("%v Waiting for thumbnail requests...", i)
			for req := range l.in {
				img, err := l.Download(req.ID)
				if err != nil {
					l.errs <- err
					continue
				}
				log.Printf("%v: %v Thumbnailing", req.ID, i)
				result, err := l.extractThumbnail(img, req.Size)
				if err != nil {
					l.errs <- err
					continue
				}
				l.out <- thumbnailer.Result{
					Request: req,
					Reader:  result,
				}
			}

			wg.Done()
		}(i)

		go func() {
			wg.Wait()
			close(l.out)
		}()
	}

	for i := 0; i < concurrency; i++ {
		go func() {
			for res := range l.out {
				log.Printf("%v: Linking thumbnail", res.ID)
				_, err := l.LinkThumb(res.ID, res.Size, res.Reader)
				if err != nil {
					l.errs <- err
					continue
				}
			}
		}()
	}
}

func (l local) handleErrs() {
	for err := range l.errs {
		log.Printf("Error: %v", err)
	}
}
