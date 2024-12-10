package basket

import (
	"fmt"
	"net/http"

	"github.com/BerkCicekler/e-commerce-audio-api/model"
	"github.com/BerkCicekler/e-commerce-audio-api/repository"
	"github.com/BerkCicekler/e-commerce-audio-api/service/auth"
	"github.com/BerkCicekler/e-commerce-audio-api/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BasketServiceHandler struct {
	repository *repository.BasketRepo
}

func BasketServiceNewHandler(repository *repository.BasketRepo) *BasketServiceHandler {
	return &BasketServiceHandler{
		repository: repository,
	}
}

func (h *BasketServiceHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/basket/", auth.WithJWTAuth(h.handleGetUserBasket)).Methods("GET")
	router.HandleFunc("/basket/update", auth.WithJWTAuth(h.handleUpdateBasketItem)).Methods("POST")
	router.HandleFunc("/basket/add", auth.WithJWTAuth(h.handleAddItemToBasket)).Methods("POST")
	router.HandleFunc("/basket/removeOne", auth.WithJWTAuth(h.handleRemoveFromBasket)).Methods("DELETE")
	router.HandleFunc("/basket/removeAll", auth.WithJWTAuth(h.handleClearUserBasket)).Methods("DELETE")
}

func (h *BasketServiceHandler) handleGetUserBasket(w http.ResponseWriter, r *http.Request) {
	uId, err := primitive.ObjectIDFromHex(auth.GetUserIDFromContext(r.Context()))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userBasket, err := h.repository.FetchUserBasket(&uId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, 200, map[string]interface{}{"basket": userBasket})
}

func (h *BasketServiceHandler) handleAddItemToBasket(w http.ResponseWriter, r *http.Request) {
	var requestModel model.BasketAddToBasketRequestModel

	uId, err := primitive.ObjectIDFromHex(auth.GetUserIDFromContext(r.Context()))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.ParseJSON(r, &requestModel); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	productObId, err := primitive.ObjectIDFromHex(requestModel.Product)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	basketId, err := h.repository.AddToBasket(&productObId, &uId)
	fmt.Println(basketId.Hex())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, 200, map[string]string{"id": basketId.Hex()})
}

func (h *BasketServiceHandler) handleUpdateBasketItem(w http.ResponseWriter, r *http.Request) {
	var requestModel model.BasketUpdateRequestModel

	if err := utils.ParseJSON(r, &requestModel); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := h.repository.AddToItemCount(requestModel.Count, &requestModel.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, 200, nil)
}

func (h *BasketServiceHandler) handleRemoveFromBasket(w http.ResponseWriter, r *http.Request) {
	var requestModel model.BasketUpdateRequestModel

	if err := utils.ParseJSON(r, &requestModel); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := h.repository.DeleteItem(&requestModel.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, 200, nil)
}

func (h *BasketServiceHandler) handleClearUserBasket(w http.ResponseWriter, r *http.Request) {
	uId, err := primitive.ObjectIDFromHex(auth.GetUserIDFromContext(r.Context()))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.repository.ClearUserBasket(&uId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, 200, nil)
}
