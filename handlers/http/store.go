package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/denbeigh2000/imstor"

	"github.com/gorilla/mux"
)

type Handler struct {
	imstor.Store

	router *mux.Router
}

func NewHandler(store imstor.Store) *Handler {
	router := mux.NewRouter()

	handler := &Handler{
		Store:  store,
		router: router,
	}

	router.HandleFunc("/", handler.CreateImage).Methods("POST")
	router.HandleFunc("/{id}", handler.RetrieveImage).Methods("GET")
	// not today
	// router.HandleFunc("/{id}", handler.UploadImage).Methods("PUT")
	router.HandleFunc("/download/{id}", handler.DownloadImage).Methods("GET")

	return handler
}

func (h *Handler) vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) RetrieveImage(w http.ResponseWriter, r *http.Request) {
	imageID := imstor.ID(h.vars(r)["id"])
	image, err := h.Retrieve(imageID)

	switch err.(type) {
	case nil:
		json.NewEncoder(w).Encode(image)
	case imstor.KeyNotFoundErr, imstor.NotUploadedYetErr:
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) CreateImage(w http.ResponseWriter, r *http.Request) {
	imageID, err := h.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	img, err := h.Upload(imageID, r.Body)
	switch err.(type) {
	case nil:
		json.NewEncoder(w).Encode(img)
	case imstor.AlreadyUploadedErr:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) DownloadImage(w http.ResponseWriter, r *http.Request) {
	imageID := imstor.ID(h.vars(r)["id"])

	reader, err := h.Download(imageID)
	switch err.(type) {
	case nil:
		n, err := io.Copy(w, reader)
		log.Printf("Wrote %v bytes with error %v", n, err)
		if err != nil {
			// Not much we can do here - we've already written a successful
			// status code, and http.Server will catch this for us.
			panic(err)
		}
	case imstor.KeyNotFoundErr, imstor.NotUploadedYetErr:
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
