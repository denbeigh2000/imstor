package thumbnailers

import (
	"log"
	"sync"

	"github.com/denbeigh2000/imstor"
	"github.com/denbeigh2000/imstor/thumbnailer"
)

const DefaultConcurrency = 12

func NewLocalThumbnailer(t thumbnailer.Thumbnailer, source imstor.ImageSource, sink imstor.ThumbnailSink, bufferSize int) thumbnailer.AsyncThumbnailer {
	thumber := local{
		Thumbnailer: t,
		BufferSize:  bufferSize,

		ImageSource:   source,
		ThumbnailSink: sink,

		in:  make(chan thumbnailer.Request, bufferSize),
		out: make(chan thumbnailer.Result, bufferSize),
	}

	thumber.loop()
	go thumber.handleErrs()
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

func (l local) loop() {
	concurrency := l.Concurrency
	if concurrency == 0 {
		concurrency = DefaultConcurrency
	}

	wg := &sync.WaitGroup{}
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			log.Printf("Waiting for thumbnail requests...")
			for req := range l.in {
				log.Printf("%v: Received thumbnail request", req.ID)
				log.Printf("%v: Downloading image", req.ID)
				img, err := l.Download(req.ID)
				if err != nil {
					l.errs <- err
					continue
				}
				log.Printf("%v: Thumbnailing", req.ID)
				result, err := l.Thumbnail(img, req.Size)
				log.Printf("%v: Queueing for upload", req.ID)
				l.out <- thumbnailer.Result{
					Request: req,
					Reader:  result,
				}
			}

			wg.Done()
		}()

		go func() {
			wg.Wait()
			close(l.out)
		}()
	}

	for i := 0; i < concurrency; i++ {
		go func() {
			for res := range l.out {
				log.Printf("%v: Uploading image", res.ID)
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
