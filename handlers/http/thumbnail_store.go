package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/denbeigh2000/imstor"
	app "github.com/denbeigh2000/imstor/app"

	"github.com/gorilla/mux"
)

type ThumbnailHandler struct {
	app.UserThumbnailAPI

	router *mux.Router
}

func NewThumbnailHandler(a app.UserThumbnailAPI) *ThumbnailHandler {
	handler := &ThumbnailHandler{
		UserThumbnailAPI: a,
	}

	router := mux.NewRouter()

	router.HandleFunc("/thumb/{id}", handler.HandleRetrieve).Methods(http.MethodGet)
	// router.HandleFunc("/thumb/{id}/{size}", handler.HandleLink).Methods(http.MethodPut)
	router.HandleFunc("/thumb/{id}/{size}/download", handler.HandleDownload).Methods(http.MethodGet)

	handler.router = router

	return handler
}

func (h *ThumbnailHandler) vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

// Not deleting yet so I have something to reference when I write a HTTP API for the store
// func (h *ThumbnailHandler) HandleLink(w http.ResponseWriter, r *http.Request) {
// 	imageID := imstor.ID(h.vars(r)["id"])
// 	size, err := imstor.FromKey(h.vars(r)["size"])
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 	}
//
// 	thumb, err := h.LinkThumb(imageID, size, r.Body)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
//
// 	err = json.NewEncoder(w).Encode(thumb)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func (h *ThumbnailHandler) HandleRetrieve(w http.ResponseWriter, r *http.Request) {
	imageID := imstor.ID(h.vars(r)["id"])

	thumbs, err := h.RetrieveThumbnails(imageID)
	switch err.(type) {
	case nil:
		json.NewEncoder(w).Encode(thumbs)
	case imstor.NoSuchThumbnailSizeErr, imstor.KeyNotFoundErr, imstor.NotUploadedYetErr:
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ThumbnailHandler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	imageID := imstor.ID(h.vars(r)["id"])
	size, err := imstor.FromKey(h.vars(r)["size"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	thumb := imstor.Thumbnail{Parent: imageID, Size: size}

	reader, err := h.DownloadThumbnail(thumb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer reader.Close()

	if reader == nil {
		http.Error(w, fmt.Sprintf("Item %v should be here, but it is empty for some reason.", imageID), http.StatusInternalServerError)
	}

	_, err = io.Copy(w, reader)
	if err != nil {
		panic(err)
	}
}

func (h *ThumbnailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
