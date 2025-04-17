package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"time"
)

type Token struct {
	Key            string
	ExpTimeAccess  int
	ExpTimeRefresh int
}

type Claims struct {
	ID     string
	UserID string
	IP     string
}

func (t *Token) CreateRefreshToken(ip, guid string) (string, string) {
	idSessionInt := rand.Int63()
	idSession := fmt.Sprintf("%v", idSessionInt)

	RefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"exp":      jwt.NewNumericDate(time.Now().Add(time.Duration(t.ExpTimeRefresh) * time.Minute)),
		"jti":      idSession,
		"guid":     guid,
		"ipClient": ip,
	})
	token, err := RefreshToken.SignedString([]byte(t.Key))
	if err != nil {
		log.Fatal(err)
	}
	return token, idSession
}

func (t *Token) CreateAccessToken(ip, guid, idSession string) string {
	AccessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"exp":      jwt.NewNumericDate(time.Now().Add(time.Duration(t.ExpTimeAccess) * time.Minute)),
		"jti":      idSession,
		"guid":     guid,
		"ipClient": ip,
	})
	token, err := AccessToken.SignedString([]byte(t.Key))
	if err != nil {
		log.Fatal(err)
	}
	return token
}

func (t *Token) CheckTokens(hashtoken string, RefreshToken []byte) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashtoken), RefreshToken)
	return err
}

func (t *Token) Parse(token string) (Claims, error) {
	tokenParse, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.Key), nil
	})
	if err != nil {
		return Claims{}, fmt.Errorf("token is invalide: %s", err)
	}

	claims, ok := tokenParse.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, errors.New("Empty claims")
	}
	claim := Claims{
		ID:     claims["jti"].(string),
		UserID: claims["guid"].(string),
		IP:     claims["ipClient"].(string),
	}
	return claim, err
}
