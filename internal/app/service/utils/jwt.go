package utils

import (
	"encoding/json"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type (
	Claims struct {
		Data any `json:"data"`
		jwt.StandardClaims
	}
)

func Sign(data any) (signatureJWT string, exp int64, err error) {
	secret := []byte(os.Getenv("JWT_SECRET_KEY"))
	bt, err := json.Marshal(data)
	if err != nil {
		return
	}
	exp = time.Now().Add(7 * 24 * time.Hour).Unix() // * 7 days
	encryptedData, err := EncryptAES256CBC(string(bt), os.Getenv("JWT_ENCRYPT_KEY"), os.Getenv("JWT_ENCRYPT_IV"))
	if err != nil {
		return
	}
	claims := &Claims{
		encryptedData,
		jwt.StandardClaims{
			Issuer:    os.Getenv("JWT_ISSUER"),
			ExpiresAt: exp,
		},
	}
	signatureJWT, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	return
}

func Verify(token string) (claims *Claims, err error) {
	secret := []byte(os.Getenv("JWT_SECRET_KEY"))
	claims = &Claims{}
	_, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return
	}
	encryptedString, err := DecryptAES256CBC(claims.Data.(string), os.Getenv("JWT_ENCRYPT_KEY"), os.Getenv("JWT_ENCRYPT_IV")) // * should be string
	if err != nil {
		return
	}
	var data any
	err = json.Unmarshal([]byte(encryptedString), &data)
	if err != nil {
		return
	}
	claims.Data = data

	return

}
