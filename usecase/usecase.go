package usecase

import (
	"TokenServer/service/db"
	"TokenServer/service/token"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type ErrorsBody struct {
	Header string // "Contents_type/json"
	Body   string // "Errors: ..."
}
type UseCase struct {
	DB     db.DbAccess
	Token  token.Token
	UserIp string
}

func (u *UseCase) CreateSession(guid string) (string, string, error) {
	refreshToken, id := u.Token.CreateRefreshToken(u.UserIp, guid)
	RefreshBcrypt, err := bcrypt.GenerateFromPassword(u.Token.HashSHA256(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	err = u.DB.Add(guid, string(RefreshBcrypt), id)
	if err != nil {
		return "", "", err
	}
	accessToken := u.Token.CreateAccessToken(u.UserIp, guid, id)
	RefreshBase64 := base64.StdEncoding.EncodeToString([]byte(refreshToken))
	return accessToken, RefreshBase64, nil
}

func (u *UseCase) RefreshSession(Access, Refresh string) (string, string, error) {
	RefreshByte, err := base64.StdEncoding.DecodeString(Refresh)
	if err != nil {
		return "", "", err
	}

	Refresh = string(RefreshByte)

	claimsRefresh, err := u.Token.Parse(Refresh)
	if err != nil {
		return "", "", err
	}
	claimsAccess, err := u.Token.Parse(Access)
	if err != nil {
		return "", "", err
	}
	guid := claimsRefresh.UserID
	if u.UserIp != claimsRefresh.IP {
		u.SendEmail(guid)
		err = errors.New("email warning")
	}

	OldRefreshBcrypt, err := u.DB.Check(claimsRefresh.ID)
	if OldRefreshBcrypt == "" {
		err = errors.New("Refresh token not exict")
		return "", "", err
	}

	if err = u.Token.Check(OldRefreshBcrypt, Refresh); err != nil {
		return "", "", err
	}

	if claimsAccess.UserID != claimsRefresh.UserID {
		return "", "", errors.New("Token mismatch")
	}

	NewRefreshToken, id := u.Token.CreateRefreshToken(u.UserIp, guid)
	NewRefreshBcrypt, err := bcrypt.GenerateFromPassword(u.Token.HashSHA256(NewRefreshToken), bcrypt.DefaultCost)
	if err != nil {
		// FIXME: Здесь и везде желательно заворачивать ошибки, добавляя к ним пояснения см. fmt.Errorf
		err = fmt.Errorf("пизданулась генераци\ бкрипта: %w", err)
		err = fmt.Errorf("Ошибка 2: ", err)
		return "", "",
	}
	// FIXME: В базе нет индекса по токену, будет долгий поиск по базе понему, мб использовать ID?
	err = u.DB.Refresh(string(NewRefreshBcrypt), claimsRefresh.ID)
	if err != nil {
		return "", "", err
	}
	RefreshBase64 := base64.StdEncoding.EncodeToString([]byte(NewRefreshToken))
	NewAccessToken := u.Token.CreateAccessToken(u.UserIp, guid, id)

	return NewAccessToken, RefreshBase64, err
}

// FIXME: Это точно публичная функция юзкейса?
func (u *UseCase) SendEmail(guid string) error {
	log.Println("Sending email: ", u.DB.GetEmail(guid))
	return nil
}
