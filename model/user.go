package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	UserName string             `json:"userName,omitempty" bson:"userName" validate:"omitempty,min=5,max=25"`
	Password string             `json:"password,omitempty" bson:"password"`
	Email    string             `json:"email,omitempty" bson:"email" validate:"required,email"`
}

type OAuthUser struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	OAuthID  string             `bson:"oAuthId" json:"oAuthId,omitempty"`
	UserName string             `json:"userName,omitempty" bson:"userName" validate:"omitempty,min=5,max=25"`
	Password string             `json:"password,omitempty" bson:"password"`
	Email    string             `json:"email,omitempty" bson:"email" validate:"required,email"`
}

type UserLoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	UserName     string `json:"userName,omitempty"`
	Email        string `json:"email,omitempty"`
}

func UserLoginResponseFromUser(user *User) *UserLoginResponse {
	return &UserLoginResponse{
		UserName: user.UserName,
		Email:    user.Email,
	}
}
