package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BerkCicekler/e-commerce-audio-api/model"
	"github.com/BerkCicekler/e-commerce-audio-api/repository"
	"github.com/BerkCicekler/e-commerce-audio-api/service/auth"
	"github.com/BerkCicekler/e-commerce-audio-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
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
	router.HandleFunc("/user/google/login", h.handleGoogleLogin).Methods("GET")
	router.HandleFunc("/user/google/callback", h.handleGoogleCallback).Methods("GET")
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

	token, err := auth.CreateJWT(u.ID.Hex())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	userResponse := model.UserLoginResponseFromUser(u)
	userResponse.Token = token

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
		ID: obId,
		UserName: user.UserName,
		Email: user.Email,
		PhoneNumber: user.PhoneNumber,
		Password: hashedPassword,
	})
	_ = result

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	
	token, err := auth.CreateJWT(obId.Hex())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	userResponse := model.UserLoginResponseFromUser(&user)
	userResponse.Token = token

	utils.WriteJSON(w, http.StatusCreated, userResponse)
}

// GOOGLE 

func (h *UserServiceHandler) handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := auth.GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *UserServiceHandler) handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := auth.GoogleOauthConfig.Exchange(context.Background(), code)

	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := auth.GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var gmailData model.GmailData
	if err := json.NewDecoder(resp.Body).Decode(&gmailData); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	user, err := h.repository.FindUserByEmail(gmailData.Email)

	if user == nil {
		user = &model.User{ID: primitive.NewObjectID(), Email: gmailData.Email, UserName: gmailData.GivenName}
		_, err = h.repository.InsertUser(user)
		if err != nil{
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	}
	
	JWTtoken, err := auth.CreateJWT(user.ID.Hex())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	userResponse := model.UserLoginResponseFromUser(user)
	userResponse.Token = JWTtoken

	utils.WriteJSON(w, http.StatusCreated, userResponse)

}