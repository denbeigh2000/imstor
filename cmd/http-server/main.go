package main

import (
	"github.com/denbeigh2000/imstor/app"
	"github.com/denbeigh2000/imstor/handlers/http"
	"github.com/denbeigh2000/imstor/stores/disk"
)

func main() {
	//store := memory.NewStore()
	httpAPI := http.HTTPAPI{}

	imageStore := disk.NewStore("./.imstor")
	thumbStore := disk.NewThumbStore("./.imstor/thumbs")

	userAPI := imstor.NewUserAPI(imageStore, thumbStore)

	imstorApp := imstor.Imstor{
		UserAPI: userAPI,
		Servers: []imstor.Server{httpAPI},
	}

	imstorApp.Serve()
}
