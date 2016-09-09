package http

import (
	"net/http"
	_ "net/http/pprof"

	app "github.com/denbeigh2000/imstor/app"

	"github.com/gorilla/mux"
)

type HTTPAPI struct {
	http.Handler

	Debug bool
}

func NewHTTPAPI(imageAPI app.UserImageAPI, thumbAPI app.UserThumbnailAPI, debug bool) HTTPAPI {
	router := mux.NewRouter()

	imageHandler := NewHandler(imageAPI)
	thumbHandler := NewThumbnailHandler(thumbAPI)

	router.PathPrefix("/thumb/").Handler(thumbHandler)
	router.PathPrefix("/").Handler(imageHandler)

	return HTTPAPI{
		Handler: router,
		Debug:   debug,
	}
}

func (a HTTPAPI) Serve(api app.UserAPI) error {
	server := http.Server{
		Addr:    ":8080",
		Handler: a,
	}

	if a.Debug {
		debugServer := http.Server{Addr: ":1337"}
		go debugServer.ListenAndServe()
	}

	return server.ListenAndServe()
}
