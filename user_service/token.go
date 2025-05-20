package userservice

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// Dont do that
var SECRET_KEY string = "SET_YOUR_JWT_KEY"

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
		return nil, fmt.Errorf("token is not valid")
	}
	return claims, nil
}
