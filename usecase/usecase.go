package usecase

import (
	"TokenServer/service/db"
	"TokenServer/service/token"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"time"
)

type ErrorsBody struct {
	Header string // "Contents_type/json"
	Body   string // "Errors: ..."
}
type UseCase struct {
	DB    db.DbAccess
	Token token.Token
	Smpt  SmtpConfig
}

type SmtpConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func (u *UseCase) CreateSession(guid, UserIp string) (string, string, error) {
	refreshToken, id := u.Token.CreateRefreshToken(UserIp, guid)

	err := u.DB.AddToken(guid, id, u.hashSHA256(refreshToken))
	if err != nil {
		return "", "", err
	}
	accessToken := u.Token.CreateAccessToken(UserIp, guid, id)
	RefreshBase64 := base64.StdEncoding.EncodeToString([]byte(refreshToken))
	return accessToken, RefreshBase64, nil
}

func (u *UseCase) RefreshSession(Refresh, UserIp string) (string, string, error) {
	RefreshByte, err := base64.StdEncoding.DecodeString(Refresh)
	if err != nil {
		return "", "", err
	}

	Refresh = string(RefreshByte)

	claimsRefresh, err := u.Token.Parse(Refresh)
	if err != nil {
		return "", "", err
	}

	guid := claimsRefresh.UserID
	if UserIp != claimsRefresh.IP {
		u.sendEmail(guid)
		err = errors.New("email warning")
	}

	OldRefreshBcrypt, err := u.DB.GetRefreshToken(claimsRefresh.ID)
	if OldRefreshBcrypt == "" {
		return "", "", errors.New("SetRefreshToken token not exist")
	}

	if err = u.Token.CheckTokens(OldRefreshBcrypt, u.hashSHA256(Refresh)); err != nil {
		return "", "", fmt.Errorf("Mismatch refresh token: %s", err)
	}

	NewRefreshToken, id := u.Token.CreateRefreshToken(UserIp, guid)

	err = u.DB.SetRefreshToken(claimsRefresh.ID, u.hashSHA256(NewRefreshToken))
	if err != nil {
		return "", "", err
	}
	RefreshBase64 := base64.StdEncoding.EncodeToString([]byte(NewRefreshToken))
	NewAccessToken := u.Token.CreateAccessToken(UserIp, guid, id)

	return NewAccessToken, RefreshBase64, err
}

func (u *UseCase) sendEmail(guid string) error {
	//TODO: РЕализовать функцию до конца
	//	log.Println("Sending email: ", u.DB.GetEmail(guid))
	server := mail.NewSMTPClient()

	server.Host = u.Smpt.Host
	server.Port = u.Smpt.Port
	server.Username = u.Smpt.Username
	server.Password = u.Smpt.Password
	server.Encryption = mail.EncryptionSSLTLS

	server.KeepAlive = true

	server.ConnectTimeout = 10 * time.Second

	server.SendTimeout = 10 * time.Second

	server.TLSConfig = &tls.Config{}

	smtpClient, err := server.Connect()

	if err != nil {
		log.Fatal(err)
	}

	email := mail.NewMSG()
	email.SetFrom(u.Smpt.Username).
		AddTo(u.DB.GetEmail(guid)).
		SetSubject("Email warning")

	email.SetBody(mail.TextPlain, "email warning")

	if email.Error != nil {
		log.Fatal(email.Error)
	}

	err = email.Send(smtpClient)
	if err != nil {
		return fmt.Errorf("Error sending email: %s", err)
	}
	return nil
}

func (u *UseCase) hashSHA256(str string) []byte {
	hash := sha256.Sum256([]byte(str))
	return hash[:]
}
