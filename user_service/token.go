package userservice

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// Dont do that
var SECRET_KEY = []byte("SET_YOUR_JWT_KEY")

type SignInDetails struct {
	Email     string
	User_id   string
	User_type string
	jwt.StandardClaims
}

func CreateTokens(email string, user_id string, user_type string) (string, string, error) {
	claims := &SignInDetails{
		Email:     email,
		User_id:   user_id,
		User_type: user_type,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	r_claims := &SignInDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(SECRET_KEY)
	if err != nil {
		return "", "", err
	}
	r_token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, r_claims).SignedString(SECRET_KEY)
	if err != nil {
		return "", "", err
	}
	return token, r_token, nil
}
