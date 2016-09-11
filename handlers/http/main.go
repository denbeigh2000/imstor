package http

import (
	"log"
	"net/http"
	"net/http/pprof"

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

	if debug {
		go func() {
			mux := http.NewServeMux()
			mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
			mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
			mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
			mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
			mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

			addr := ":1337"
			log.Printf("Debug listening on %v", addr)

			debugServer := http.Server{Addr: addr, Handler: mux}
			log.Fatalf("pprof failed: %v", debugServer.ListenAndServe())
		}()
	}

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

	return server.ListenAndServe()
}
