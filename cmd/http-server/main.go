package main

import (
	imstorLib "github.com/denbeigh2000/imstor"
	"github.com/denbeigh2000/imstor/app"
	"github.com/denbeigh2000/imstor/handlers/http"
	"github.com/denbeigh2000/imstor/stores/memory"
)

func main() {
	store := memory.NewStore()
	httpAPI := http.HTTPAPI{}

	userAPI := imstor.NewUserAPI(store, store.(imstorLib.ThumbnailStore))

	imstorApp := imstor.Imstor{
		UserAPI: userAPI,
		Servers: []imstor.Server{httpAPI},
	}

	imstorApp.Serve()
}
