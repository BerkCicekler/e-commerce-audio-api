package image

import (
	"net/http"

	"github.com/gorilla/mux"
)

type ImageServiceHandler struct {
	BaseDir string
}

func (h *ImageServiceHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/image/{filename}", h.handleImage).Methods("GET")
}

func (s *ImageServiceHandler) handleImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	imagePath := s.BaseDir + "/images/" + filename

	w.Header().Set("Content-Type", "image/png")

	http.ServeFile(w, r, imagePath)
}
