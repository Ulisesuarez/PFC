package controller

import (
	idonia_auth "bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia-auth"
	"crypto/sha256"
	"encoding/base64"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"

	"database/sql"
	"encoding/json"
	"time"

	"bitbucket.org/inehealth/api/api"
	"bitbucket.org/inehealth/idonia-pacs/configuration"
	"bitbucket.org/inehealth/idonia-pacs/model"
	"bitbucket.org/inehealth/idonia-pacs/repository"
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia"
	"bitbucket.org/inehealth/idonia-pacs/util"
	"github.com/google/uuid"
)

//AuthHandler ...
type AuthHandler struct {
	Conf   *configuration.Configuration `inject:"configuration"`
	DB     *sql.DB                      `inject:"db"`
	Logger *logrus.Logger               `inject:"logger"`
}

//GetEndpoints ...
func (ah AuthHandler) GetEndpointList() api.EndpointList {

	return api.EndpointList{
		&api.Endpoint{Path: "auth/login", Methods: []string{http.MethodPost}}:  ah.Login,
		&api.Endpoint{Path: "auth/logout", Methods: []string{http.MethodPost}}: ah.Logout,
	}

}

//GetStaticList
func (ah AuthHandler) GetStaticList() api.StaticList {
	return api.StaticList{}
}

//Login ...
func (ah AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	//	authService := auth.AuthenticationService(r)

	var loginRQ *model.LoginRQ

	bodyString, _ := util.ReadBodyFromRequest(r, loginRQ)
	_ = json.Unmarshal(bodyString, &loginRQ)

	accountID := idonia.AccountID
	token := idonia.Token
	apiKey := idonia.ApiKey
	apiSecret := idonia.ApiSecret

	/*	if  authService!= nil && authService.IsValid() {
		util.JSONResponse(
			model.LoginRS{Token: authService.GetToken()},
			w,
			http.StatusOK,
		)
		return
	}*/

	if ah.Conf.Login.Type != "None" && (loginRQ == nil || loginRQ.Username == "" || loginRQ.Password == "") {
		ah.Logger.Println(ah.Conf.Login.Type)
		ah.Logger.Println(loginRQ.Username)
		util.JSONResponse(
			model.ErrUserInvalidCredentials,
			w,
			http.StatusUnauthorized,
		)
		return
	}

	if ah.Conf.Login != nil && ah.Conf.Login.Type == "Config" {
		found := false
		for _, v := range ah.Conf.Login.Users {
			ah.Logger.Info(loginRQ.Username)
			ah.Logger.Info(v.Username)
			ah.Logger.Info(loginRQ.Password)
			ah.Logger.Info(v.Password)
			if v.Username == loginRQ.Username && v.Password == loginRQ.Password {
				found = true
				break
			}
		}
		ah.Logger.Info(found)
		if !found {
			util.JSONResponse(
				model.ErrUserInvalidCredentials,
				w,
				http.StatusUnauthorized,
			)
			return
		}

	}
	accountID = 0
	token = ""
	userHasher := sha256.New()
	userHasher.Write([]byte(strings.ToLower(ah.Conf.Idonia.APiUsername)))
	idoniaEmail := base64.StdEncoding.EncodeToString(userHasher.Sum(nil))
	passHasher := sha256.New()
	passHasher.Write([]byte(ah.Conf.Idonia.ApiPass))
	idoniaPassword := base64.StdEncoding.EncodeToString(passHasher.Sum(nil))
	res, err := idonia_auth.PostAccountLogin(&idonia_auth.PostAccountLoginReq{
		Email:    idoniaEmail,
		Password: idoniaPassword,
	}, nil, ah.Logger)
	if err != nil {
		ah.Logger.Error(err.Error())
		util.JSONResponse(
			model.ErrUserInvalidCredentials,
			w,
			http.StatusUnauthorized,
		)
		return
	}
	accountID = res.AccountID
	idonia.Token = res.Token
	token = res.Token
	uuid := uuid.New().String()

	err = repository.InsertSession(ah.DB, &repository.Session{
		UUID:            uuid,
		ExpireDate:      time.Now().Add(1 * time.Hour),
		IdoniaAccountID: accountID,
		IdoniaToken:     token,
		IdoniaAPIKey:    apiKey,
		IdoniaAPISecret: apiSecret,
	})
	if err != nil {
		retErr := model.ErrServer
		retErr.AdditionalMessage = err.Error()
		util.JSONResponse(
			retErr,
			w,
			http.StatusInternalServerError,
		)
		return
	}

	util.JSONResponse(
		model.LoginRS{Token: uuid},
		w,
		http.StatusOK,
	)

}

//Login ...

func (ah AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

	util.JSONResponse(
		"OK",
		w,
		http.StatusOK,
	)

}
