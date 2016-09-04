package main

import (
	"io"
	"log"
	"os"

	"github.com/denbeigh2000/imstor"
	"github.com/denbeigh2000/imstor/thumbnailer"
)

const (
	inFile  = "test.jpg"
	outFile = "out.jpg"
)

func main() {
	thumbnailer := thumbnailers.NewLocal()

	f, err := os.Open(inFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	out, err := thumbnailer.Thumbnail(f, imstor.Size{LongEdge: 300})
	if err != nil {
		log.Fatal(err)
	}

	g, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	_, err = io.Copy(g, out)
	if err != nil {
		log.Fatal(err)
	}
}
