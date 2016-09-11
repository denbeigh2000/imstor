package validator

import (
	"fmt"
	"io"
	"time"

	"github.com/denbeigh2000/imstor"
)

func NewRequest(data io.Reader) Request {
	return Request{
		Data: data,

		out: make(chan Response),
	}
}

type Request struct {
	Data io.Reader

	out chan Response
}

type Response struct {
	imstor.ImageInfo
	Err string
}

func (r Request) Respond(resp Response) {
	r.out <- resp
	close(r.out)
}

func (r Request) Response() Response {
	resp, ok := <-r.out
	if !ok {
		panic("not ok")
	}

	return resp
}

type Processor interface {
	Process(Request) error
	Stop()
}

func NewLocalProcessor(v Validator, bufferSize, concurrency int) Processor {
	processor := localProcessor{
		Validator: v,

		concurrency: concurrency,
		traffic:     make(chan Request, bufferSize),
	}

	processor.startWorkers()
	return processor
}

func NewTimedLocalProcessor(v Validator, timeout time.Duration, bufferSize, concurrency int) Processor {
	processor := localProcessor{
		Validator: v,

		timeout:     timeout,
		concurrency: concurrency,
		traffic:     make(chan Request, bufferSize),
	}

	processor.startWorkers()
	return processor
}

type localProcessor struct {
	Validator

	concurrency int
	traffic     chan Request

	timeout time.Duration
}

func (p localProcessor) Process(r Request) error {
	if p.timeout == 0 {
		p.traffic <- r
		return nil
	}

	select {
	case p.traffic <- r:
		return nil
	case <-time.After(p.timeout):
		return fmt.Errorf("Took too long to accept validation request")
	}
}

func (p localProcessor) Stop() {
	close(p.traffic)
}

func (p localProcessor) startWorkers() {
	for i := 0; i < p.concurrency; i++ {
		go p.work()
	}
}

func (p localProcessor) work() {
	for req := range p.traffic {
		info, err := p.Validate(req.Data)
		var errStr string
		if err != nil {
			errStr = err.Error()
		}

		resp := Response{
			ImageInfo: info,
			Err:       errStr,
		}

		req.Respond(resp)
	}
}
