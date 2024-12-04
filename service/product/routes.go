package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/BerkCicekler/e-commerce-audio-api/model"
	"github.com/BerkCicekler/e-commerce-audio-api/repository"
	"github.com/BerkCicekler/e-commerce-audio-api/service/auth"
	"github.com/BerkCicekler/e-commerce-audio-api/utils"
	"github.com/gorilla/mux"
)

type ProductServiceHandler struct {
	repository *repository.ProductRepo
}

func ProductServiceNewHandler(repository *repository.ProductRepo) *ProductServiceHandler {
	return &ProductServiceHandler{repository: repository}
}

func (h *ProductServiceHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/shop/featured", auth.WithJWTAuth(h.handleCategories)).Methods("GET")
}

func (h *ProductServiceHandler) handleCategories(w http.ResponseWriter, r *http.Request) {

	search := r.URL.Query().Get("search")
	offset := r.URL.Query().Get("offset")
	var offsetInt int64
	if offset == "" {
		offsetInt = 0
	} else {
		offsetInt, _ = strconv.ParseInt(offset, 10, 64)
	}

	var productRequest model.ProductRequest
	if err := utils.ParseJSON(r, &productRequest); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	products, err := h.repository.FetchFeatured(search, offsetInt, &productRequest)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("couldn't fetch the Products"))
		return
	}

	utils.WriteJSON(w, 200, map[string]interface{}{"products": products})
}
