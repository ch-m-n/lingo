package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("supersecretkey")
type JWTClaim struct {
	Id 		 string	`json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(id string,email string, username string) (tokenString string, err error) {
	expirationTime := time.Now().Add(2*8760*time.Hour)
	claims:= &JWTClaim{
		Id: id,
		Email: email,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return
}