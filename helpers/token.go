package helpers

import (
	// "log"
	//  "encoding/json"

	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"

	// "net/http"
	"time"
)

type TokenClaims struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

var sampleSecretKey = []byte("MySecretKey")

func GenerateJWT(userId string, userName string) (string, error) {

	expireToken := time.Now().Add(time.Hour * 2).Unix()

	// Set-up claims
	claims := TokenClaims{
		ID:   userId,
		Name: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// claims := token.Claims.(jwt.MapClaims)
	// claims["exp"] = time.Now().Add(2 * time.Hour)
	// claims["authorized"] = true
	// claims["id"] = userId
	// claims["name"] = userName
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJwt(receivedToken string) (string, string, error) {

	if receivedToken == "" {
		return "", "", errors.New("No token in request")
	}

	var keyfunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected siging method")
		}
		return []byte(sampleSecretKey), nil
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(receivedToken, claims, keyfunc)

	if err != nil {
		fmt.Print(err)
		return "", "", err
	}

	id := claims["id"].(string)
	name := claims["name"].(string)

	return id, name, nil

}
