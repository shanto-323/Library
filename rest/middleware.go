package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

var SECRET_KEY = []byte("SET_YOUR_JWT_KEY")

type SignInDetails struct {
	Email     string
	User_id   string
	User_type string
	jwt.StandardClaims
}

type contextKey string

const userTypeKey = contextKey("user_type")

func JwtMiddleWere(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r_token := r.Header.Get("token")
		if r_token == "" {
			WriteJson(w, http.StatusNotAcceptable, "nil token")
			return
		}

		claims, err := ValidateToken(r_token)
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorExpired != 0 {
					WriteJson(w, http.StatusUnauthorized, "token expired")
					return
				}
			}
			WriteJson(w, http.StatusForbidden, err)
			return
		}

		ctx := context.WithValue(r.Context(), userTypeKey, claims.User_type)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateToken(signedToken string) (*SignInDetails, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignInDetails{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("token encryption method not matching")
			}
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignInDetails)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}
