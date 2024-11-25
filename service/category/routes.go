package category

import (
	"fmt"
	"net/http"

	"github.com/BerkCicekler/e-commerce-audio-api/repository"
	"github.com/BerkCicekler/e-commerce-audio-api/service/auth"
	"github.com/BerkCicekler/e-commerce-audio-api/utils"
	"github.com/gorilla/mux"
)

type CategoryServiceHandler struct{
	repository *repository.CategoriesRepo
}

func CategoriesServiceNewHandler(repository *repository.CategoriesRepo) *CategoryServiceHandler {
	return &CategoryServiceHandler{repository: repository}
}

func (h *CategoryServiceHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/categories/", auth.WithJWTAuth(h.handleCategories) ).Methods("GET")
}

func (h *CategoryServiceHandler) handleCategories(w http.ResponseWriter, r *http.Request) {
	
	categories, err :=  h.repository.GetCategories()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("couldn't fetch the categories"))
		return
	}

	utils.WriteJSON(w, 200, map[string]interface{}{"categories": categories})
}