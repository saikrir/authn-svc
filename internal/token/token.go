package token

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

type Token = string

type JWTManager struct {
	SeacretKey      []byte
	ValidityInHours int
}

func NewJWTManager(secretKey string, validityInHours int) *JWTManager {
	return &JWTManager{
		SeacretKey:      []byte(secretKey),
		ValidityInHours: validityInHours,
	}
}

func (j *JWTManager) NewToken(subject string) (Token, int64, error) {
	expires_at := time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"subject": subject,
			"exp":     expires_at,
		})

	tokenStr, err := token.SignedString(j.SeacretKey)
	if err != nil {
		return "", expires_at, err
	}
	return tokenStr, expires_at, nil

}

func (j *JWTManager) VerifyToken(tokenString Token) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return j.SeacretKey, nil
	})

	if err != nil {
		log.Println("failed to parse token, err : ", err)
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token ")
	}

	return nil
}
