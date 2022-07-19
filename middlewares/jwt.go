package middlewares

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	expMinites = 60 * 24 * 3 // 3 days
)

type sessionJwtClaims struct {
	SessionID string `json:"sessionId"`
	jwt.StandardClaims
}

func getPrivateKey() []byte {
	rootDir := os.Getenv("WORKDIR")
	signingKey, err := ioutil.ReadFile(rootDir + "/jwt.key")
	if err != nil {
		panic(err)
	}
	return signingKey
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	privateKey := getPrivateKey()
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return privateKey, nil
}

func GenerateJwt(sessionId string) (string, error) {
	privateKey := getPrivateKey()
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = sessionJwtClaims{
		SessionID: sessionId,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * expMinites).Unix(),
			Issuer:    "ebra",
		},
	}
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func ValidateJwt(tokenString string) (jwt.MapClaims, error) {
	parser := new(jwt.Parser)
	token, err := parser.Parse(tokenString, keyFunc)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
