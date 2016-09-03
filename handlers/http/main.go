package http

import (
	"net/http"

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

	return server.ListenAndServe()
}
