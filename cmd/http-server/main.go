package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/denbeigh2000/imstor/app"
	imstorHTTP "github.com/denbeigh2000/imstor/handlers/http"
	"github.com/denbeigh2000/imstor/stores/disk"

	"github.com/gorilla/mux"
)

const (
	DefaultAPIPrefix    = "/api"
	DefaultListenPort   = 8080
	DefaultFrontendPort = 3000
)

func frontendURL() *url.URL {
	port := DefaultFrontendPort
	uri, err := url.Parse(fmt.Sprintf("http://localhost:%v/", port))
	if err != nil {
		panic(err)
	}
	return uri
}

func frontendHandler() http.Handler {
	uri := frontendURL()
	return httputil.NewSingleHostReverseProxy(uri)
}

func apiHandler() http.Handler {
	//store := memory.NewStore()
	imageStore := disk.NewStore("./.imstor")
	thumbStore := disk.NewThumbStore("./.imstor/thumbs")

	userAPI := imstor.NewUserAPI(imageStore, thumbStore)

	imageAPI := userAPI.(imstor.UserImageAPI)
	thumbAPI := userAPI.(imstor.UserThumbnailAPI)

	httpAPI := imstorHTTP.NewHTTPAPI(imageAPI, thumbAPI, false)

	prefix := DefaultAPIPrefix
	return http.StripPrefix(prefix, httpAPI)
}

func main() {
	fe := frontendHandler()
	api := apiHandler()

	router := mux.NewRouter()
	router.PathPrefix("/api/").Handler(api)
	router.PathPrefix("/").Handler(fe)

	listenURL := fmt.Sprintf(":%v", DefaultListenPort)

	server := http.Server{
		Addr:    listenURL,
		Handler: router,
	}

	log.Printf("Listening on %v...\n", listenURL)

	server.ListenAndServe()
}
