package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BerkCicekler/e-commerce-audio-api/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "userID"

func WithJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.GetTokenFromRequest(r)

		token, err := ValidateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		isExpired := validateIsJWTExpire(claims)

		if isExpired {
			log.Println("token has expired")
			utils.WriteError(w, http.StatusUnauthorized, errors.New("token expired"))
			return
		}

		str := claims["userID"].(string)

		// Add the user to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, str)
		r = r.WithContext(ctx)

		// Call the function if the token is valid
		handlerFunc(w, r)
	}
}

func CreateJWT(userID string) (string, string, error) {
	accessTokenExpiration := time.Second * time.Duration(6000)
	refreshTokenExpiration := time.Second * time.Duration(600000)

	// Access Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"expiresAt": time.Now().Add(accessTokenExpiration).Unix(),
	})

	// Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"expiresAt": time.Now().Add(refreshTokenExpiration).Unix(),
	})

	// JWT Secret
	secret := []byte(os.Getenv("JWT_SECRET"))

	// Sign Access Token
	accessTokenString, err := accessToken.SignedString(secret)
	if err != nil {
		return "", "", err
	}

	// Sign Refresh Token
	refreshTokenString, err := refreshToken.SignedString(secret)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}

func validateIsJWTExpire(claims jwt.MapClaims) bool {
	// Check token expiration
	expiresAt, ok := claims["expiresAt"].(float64) // jwt.MapClaims stores numbers as float64
	if !ok {
		log.Println("missing or invalid expiresAt in token")
		return true
	}

	// Current time
	now := time.Now().Unix()
	if int64(expiresAt) < now {
		log.Println("token has expired")
		return true
	}
	return false
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(UserKey).(string)
	if !ok {
		return ""
	}

	return userID
}
