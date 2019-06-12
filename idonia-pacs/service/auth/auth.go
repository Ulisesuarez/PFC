package auth

import (
	"bitbucket.org/inehealth/idonia-pacs/repository"
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Auth struct {
	Token           string
	expireDatetime  time.Time
	idoniaAccountID uint32
	idoniaToken     string
	apiKey          string
	apiSecret       string
}

func AuthenticationService(r *http.Request) *Auth {
	fmt.Println(r.Context().Value("authorization"))
	auth, ok := r.Context().Value("authorization").(*Auth)
	if ok {
		return auth
	}
	return nil
}

func (auth *Auth) Init(r *http.Request, db *sql.DB) (err error) {
	authHeader := r.Header.Get("Authorization")
	token := ""
	if len(authHeader) > 0 && strings.Contains(authHeader, "Bearer") {
		token = strings.Trim(authHeader[len("Bearer"):], " ")
	}
	session, err := repository.GetSession(db, token)
	if err != nil && err != sql.ErrNoRows {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	if session == nil {
		return
	}
	if session.ExpireDate.Before(time.Now()) {
		err = repository.DeleteSession(db, token)
		return nil
	}

	auth.Token = token
	auth.expireDatetime = session.ExpireDate
	auth.idoniaAccountID = session.IdoniaAccountID
	auth.idoniaToken = session.IdoniaToken
	auth.apiKey = session.IdoniaAPIKey
	auth.apiSecret = session.IdoniaAPISecret
	return
}

func (auth *Auth) GetToken() (token string) {
	return auth.Token
}

func (auth *Auth) IsValid() (ok bool) {
	return auth.GetToken() != ""
}

func (auth *Auth) GetIdoniaAuth() (idoniaAuth *idonia.Auth) {
	return &idonia.Auth{
		AccountID: auth.idoniaAccountID,
		Token:     auth.idoniaToken,
		APIKey:    auth.apiKey,
		APISecret: auth.apiSecret,
	}
}
