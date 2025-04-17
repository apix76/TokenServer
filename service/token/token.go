package token

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
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
		// FIXME: refresh token никогда не истекает?
		//"exp":      jwt.NewNumericDate(time.Now().Add(time.Duration(t.ExpTimeRefresh) * time.Minute)),
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
		// FIXME: access token никогда не истекает?
		//	"exp":      jwt.NewNumericDate(time.Now().Add(time.Duration(t.ExpTimeAccess) * time.Minute)),
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

// FIXME: Подумать о том, чтобы опустить bcrypt внутрь DBAccess, сейчас это размазано по всем слоям и какая-то каша выходит.
func (t *Token) Check(hashtoken, Refresh string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashtoken), t.HashSHA256(Refresh))
	return err
}

func (t *Token) Parse(token string) (Claims, error) {
	tokenParse, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.Key), nil
	})
	if err != nil {
		return Claims{}, err
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

func (t *Token) HashSHA256(str string) []byte {
	hash := sha256.Sum256([]byte(str))
	return hash[:]
}
