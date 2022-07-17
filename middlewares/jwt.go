package middlewares

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	privateKey = []byte("key")
)

const (
	expMinites = 60 * 24 * 3
)

func GenerateJwt() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * expMinites).Unix(),
		Issuer:    "foo",
	}
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func ValidateJwt(tokenString string) (jwt.MapClaims, error) {
	parser := new(jwt.Parser)
	token, err := parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return privateKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
