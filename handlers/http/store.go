package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/denbeigh2000/imstor"
	app "github.com/denbeigh2000/imstor/app"

	"github.com/gorilla/mux"
)

type Handler struct {
	app.UserImageAPI

	router *mux.Router
}

func NewHandler(a app.UserImageAPI) *Handler {
	router := mux.NewRouter()

	handler := &Handler{
		UserImageAPI: a,
		router:       router,
	}

	router.HandleFunc("/", handler.HandleCreate).Methods(http.MethodPost)
	router.HandleFunc("/{id}", handler.HandleRetrieve).Methods(http.MethodGet)
	router.HandleFunc("/{id}/download", handler.HandleDownload).Methods(http.MethodGet)

	return handler
}

func (h *Handler) vars(r *http.Request) map[string]string {
	log.Println("Serving HTTP!")
	return mux.Vars(r)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) HandleRetrieve(w http.ResponseWriter, r *http.Request) {
	imageID := imstor.ID(h.vars(r)["id"])
	image, err := h.RetrieveImage(imageID)

	switch err.(type) {
	case nil:
		json.NewEncoder(w).Encode(image)
	case imstor.KeyNotFoundErr, imstor.NotUploadedYetErr:
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	img, err := h.CreateImage(r.Body)

	switch err.(type) {
	case nil:
		json.NewEncoder(w).Encode(img)
	case imstor.AlreadyUploadedErr:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	imageID := imstor.ID(h.vars(r)["id"])

	reader, err := h.DownloadImage(imageID)
	switch err.(type) {
	case nil:
		_, err := io.Copy(w, reader)
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
