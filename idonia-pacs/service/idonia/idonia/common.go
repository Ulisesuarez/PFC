package idonia

import (
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ApiKey string
var ApiSecret string
var UserAgent = "Idonia Connect"
var APIHost string
var AccountID uint32
var Token string

var httpClient *http.Client

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type Auth struct {
	AccountID uint32 `json:"account_id"`
	Token     string `json:"token"`
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

func Get(u *url.URL, auth *Auth) (response *http.Response, err error) {
	return request("GET", u, nil, auth, "")
}

func Put(u *url.URL, b io.Reader, auth *Auth) (response *http.Response, err error) {
	return request("PUT", u, b, auth, "")
}

func Post(u *url.URL, b io.Reader, auth *Auth) (response *http.Response, err error) {
	return request("POST", u, b, auth, "")
}

func PostWithContentType(u *url.URL, b io.Reader, auth *Auth, contentType string) (response *http.Response, err error) {
	return request("POST", u, b, auth, contentType)
}

func Delete(u *url.URL, b io.Reader, auth *Auth) (response *http.Response, err error) {
	return request("DELETE", u, nil, auth, "")
}

func request(verb string, u *url.URL, b io.Reader, auth *Auth, contentType string) (response *http.Response, err error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if contentType == "" {
		contentType = "application/json"
	}
	req, err := http.NewRequest(verb, u.String(), b)
	if err != nil {
		fmt.Println("1")
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-type", contentType)

	fmt.Println("TOKEN Y CUENTA", Token, AccountID)
	if auth != nil {
		fmt.Println("TOKEN Y CUENTA2", auth.Token, auth.AccountID)
	} else {
		fmt.Println("NO")
	}
	if strings.Contains(u.String(), "login") {

	} else if auth != nil && auth.Token != "" && auth.AccountID > 0 {
		fmt.Println("authtoken", auth, auth.Token)
		req.Header.Set("Authorization", fmt.Sprintf("IDONIA token=\"%s\",account_id=\"%d\",socket_id=\"\"", auth.Token, auth.AccountID))

	} else if auth != nil && strings.HasPrefix(auth.APIKey, "K2") && strings.HasPrefix(auth.APISecret, "S2") {
		fmt.Println("authkey", auth, auth.APIKey)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": auth.APIKey,
			"iat": time.Now().Add(-5 * time.Minute).Unix(),
			"exp": time.Now().Add(5 * time.Minute).Unix(),
		})
		tokenString, err := token.SignedString(auth.APISecret[2:])
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	} else if Token != "" && AccountID > 0 {
		fmt.Println("token floating", auth, auth.APIKey)
		req.Header.Set("Authorization", fmt.Sprintf("IDONIA token=\"%s\",account_id=\"%d\",socket_id=\"\"", Token, AccountID))

	} else if strings.HasPrefix(ApiKey, "K2") && strings.HasPrefix(ApiSecret, "S2") {
		fmt.Println("HOLAAA", ApiKey, ApiSecret)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": ApiKey,
			"iat": time.Now().Add(-5 * time.Minute).Unix(),
			"exp": time.Now().Add(5 * time.Minute).Unix(),
		})
		apiSecret, err := base64.URLEncoding.DecodeString(ApiSecret[2:])
		if err != nil {
			fmt.Println("2")
			return nil, err
		}
		tokenString, err := token.SignedString(apiSecret)
		if err != nil {
			fmt.Println("3")
			return nil, err
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	}
	if err != nil {
		return nil, err
	}

	response, err = httpClient.Do(req)
	if err != nil {
		fmt.Println("4")
		return nil, err
	}
	fmt.Println("5")
	return
}
