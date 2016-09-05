package http

import (
	"net/http"
	_ "net/http/pprof"

	app "github.com/denbeigh2000/imstor/app"

	"github.com/gorilla/mux"
)

type HTTPAPI struct{}

func (a HTTPAPI) Serve(api app.UserAPI) error {
	router := mux.NewRouter()

	imageHandler := NewHandler(api.(app.UserImageAPI))
	thumbHandler := NewThumbnailHandler(api.(app.UserThumbnailAPI))

	router.PathPrefix("/thumb/").Handler(thumbHandler)
	router.PathPrefix("/").Handler(imageHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	debugServer := http.Server{Addr: ":1337"}
	go debugServer.ListenAndServe()

	return server.ListenAndServe()
}
