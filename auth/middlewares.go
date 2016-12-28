package auth

import (
	"fmt"
	"net/http"

	"errors"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

func RequireTokenAuthentication(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authBackend := InitJWTAuthenticationBackend()

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			return authBackend.PublicKey, nil
		}
	})

	if err == nil && token.Valid {
		fmt.Println(token.Claims.(jwt.MapClaims)["sub"])
		next(w, r)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func GetTokenFromRequest(r *http.Request) (string, error) {
	authBackend := InitJWTAuthenticationBackend()
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			return authBackend.PublicKey, nil
		}
	})
	if err == nil && token.Valid {
		return token.Claims.(jwt.MapClaims)["sub"].(string), nil
	} else {
		return "", errors.New("Not authorized")
	}
}
