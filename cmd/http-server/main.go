package main

import (
	"log"
	"net/http"

	imstorHTTP "github.com/denbeigh2000/imstor/handlers/http"
	"github.com/denbeigh2000/imstor/stores/memory"
)

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: imstorHTTP.NewHandler(memory.NewStore()),
	}

	log.Fatal(server.ListenAndServe())
}
