package main

import (
	"net/http"

	"github.com/denbeigh2000/imstor/app"
	imstorHTTP "github.com/denbeigh2000/imstor/handlers/http"
	"github.com/denbeigh2000/imstor/stores/disk"
)

func main() {
	//store := memory.NewStore()
	imageStore := disk.NewStore("./.imstor")
	thumbStore := disk.NewThumbStore("./.imstor/thumbs")

	userAPI := imstor.NewUserAPI(imageStore, thumbStore)

	imageAPI := userAPI.(imstor.UserImageAPI)
	thumbAPI := userAPI.(imstor.UserThumbnailAPI)

	httpAPI := imstorHTTP.NewHTTPAPI(imageAPI, thumbAPI, false)

	server := http.Server{
		Addr:    ":8080",
		Handler: httpAPI,
	}

	server.ListenAndServe()
}
