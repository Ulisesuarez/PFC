package idonia_auth

import (
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
)

func ApiKeyLogin() (err error) {
	username := idonia.ApiKey
	password := idonia.ApiSecret
	login := &PostAccountLoginReq{
		Username: username[2:],
		Password: password[2:],
	}
	res, err := PostAccountLogin(login, nil, nil)
	if err != nil {
		return
	}
	idonia.Token = res.Token
	idonia.AccountID = res.AccountID
	return nil
}

type PostAccountLoginReq struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
	Code2FA  string `json:"code_2fa,omitempty"`
}

type PostAccountLoginRes struct {
	Token     string `json:"token"`
	AccountID uint32 `json:"account_id"`
}

func PostAccountLogin(body *PostAccountLoginReq, auth *idonia.Auth, logger *logrus.Logger) (res *PostAccountLoginRes, err error) {
	u, err := url.Parse(idonia.APIHost + "/account/login")
	if err != nil {
		logger.Error("1", idonia.APIHost+"/account/login")
		return
	}
	encodedBody, err := json.Marshal(body)
	if err != nil {
		logger.Error("2")
		return
	}
	response, err := idonia.Post(u, bytes.NewReader(encodedBody), auth)
	if err != nil {
		logger.Error("3", err.Error(), response, idonia.APIHost+"/account/login")
		return
	}
	defer response.Body.Close()
	var bodyString []byte
	bodyString, err = ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("4", err.Error())
		return
	}
	if response.StatusCode != http.StatusOK {
		logger.Error("5", fmt.Sprintf("%+v, %+v ", response, response.Request))
		err = errors.New(fmt.Sprintf("[%d] Server said: %s", response.StatusCode, string(bodyString)))
		return
	}
	err = json.Unmarshal(bodyString, &res)
	if res == nil {
		logger.Error("6", response.StatusCode)
		return nil, errors.New(fmt.Sprintf("[%d] Server said: %s", response.StatusCode, string(bodyString)))
	}
	return
}

type GetAccountRes struct {
	Account struct {
		Username string `json:"username"`
	} `json:"account"`
}

func GetAccount(auth *idonia.Auth) (res *GetAccountRes, err error) {
	u, err := url.Parse(idonia.APIHost + "/account")
	if err != nil {
		return
	}
	response, err := idonia.Get(u, auth)
	if err != nil {
		return
	}
	defer response.Body.Close()
	var bodyString []byte
	bodyString, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("[%d] Server said: %s", response.StatusCode, string(bodyString)))
		return
	}
	err = json.Unmarshal(bodyString, &res)
	if res == nil {
		return nil, errors.New(fmt.Sprintf("[%d] Server said: %s", response.StatusCode, string(bodyString)))
	}
	return
}
