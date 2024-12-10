package user

import (
	"fmt"
	"net/http"

	"github.com/BerkCicekler/e-commerce-audio-api/model"
	"github.com/BerkCicekler/e-commerce-audio-api/repository"
	"github.com/BerkCicekler/e-commerce-audio-api/service/auth"
	"github.com/BerkCicekler/e-commerce-audio-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserServiceHandler struct {
	repository repository.UsersRepo
}

func UserServiceNewHandler(repository repository.UsersRepo) *UserServiceHandler {
	return &UserServiceHandler{repository: repository}
}

func (h *UserServiceHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/user/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/user/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/user/oauth", h.handleOauth).Methods("POST")
	router.HandleFunc("/user/refreshToken", auth.WithJWTAuth(h.handleRefreshToken)).Methods("POST")
}

func (h *UserServiceHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	u, err := h.repository.FindUserByEmail(user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !auth.ComparePasswords(u.Password, user.Password) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	token, refreshToken, err := auth.CreateJWT(u.ID.Hex())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	userResponse := model.UserLoginResponseFromUser(u)
	userResponse.Token = token
	userResponse.RefreshToken = refreshToken

	utils.WriteJSON(w, http.StatusOK, userResponse)
}

func (h *UserServiceHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// check if user exists
	_, err := h.repository.FindUserByEmail(user.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", user.Email))
		return
	}

	// hash password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	obId := primitive.NewObjectID()

	result, err := h.repository.InsertUser(&model.User{
		ID:       obId,
		UserName: user.UserName,
		Email:    user.Email,
		Password: hashedPassword,
	})
	_ = result

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	token, refreshToken, err := auth.CreateJWT(obId.Hex())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	userResponse := model.UserLoginResponseFromUser(&user)
	userResponse.Token = token
	userResponse.RefreshToken = refreshToken

	utils.WriteJSON(w, http.StatusCreated, userResponse)
}

func (h *UserServiceHandler) handleOauth(w http.ResponseWriter, r *http.Request) {
	var userRequest model.OAuthUser
	if err := utils.ParseJSON(r, &userRequest); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(userRequest); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	var obId primitive.ObjectID

	// check if user exists
	user, err := h.repository.FindUserByEmail(userRequest.Email)
	if user == nil {
		obId = primitive.NewObjectID()
		_, err := h.repository.InsertOAuthUser(&model.OAuthUser{
			ID:       obId,
			OAuthID:  userRequest.OAuthID,
			UserName: userRequest.UserName,
			Email:    userRequest.Email,
		})
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Server error: %v"))
			return
		}
	} else {
		obId = user.ID
	}

	token, refreshToken, err := auth.CreateJWT(obId.Hex())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	userResponse := model.UserLoginResponseFromUser(&model.User{
		UserName: userRequest.UserName,
		Email:    userRequest.Email,
	})
	userResponse.Token = token
	userResponse.RefreshToken = refreshToken

	utils.WriteJSON(w, http.StatusCreated, userResponse)
}

func (h *UserServiceHandler) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	uId, err := primitive.ObjectIDFromHex(auth.GetUserIDFromContext(r.Context()))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.repository.FindUserById(uId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("User doesn't exist"))
		return
	}

	token, refreshToken, err := auth.CreateJWT(uId.Hex())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	userResponse := model.UserLoginResponseFromUser(user)
	userResponse.Token = token
	userResponse.RefreshToken = refreshToken
	utils.WriteJSON(w, http.StatusCreated, userResponse)

}
