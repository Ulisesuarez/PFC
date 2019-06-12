package controller

import (
	"bitbucket.org/inehealth/api/api"
	"bitbucket.org/inehealth/idonia-pacs/configuration"
	"bitbucket.org/inehealth/idonia-pacs/model"
	"bitbucket.org/inehealth/idonia-pacs/repository"
	"bitbucket.org/inehealth/idonia-pacs/service/auth"
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia"
	"bitbucket.org/inehealth/idonia-pacs/service/study"
	"bitbucket.org/inehealth/idonia-pacs/util"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

//TaskHandler ...
type StudyHandler struct {
	Conf   *configuration.Configuration `inject:"configuration"`
	DB     *sql.DB                      `inject:"db"`
	Logger *logrus.Logger               `inject:"logger"`
}

//GetStaticList
func (sh StudyHandler) GetStaticList() api.StaticList {
	return api.StaticList{}
}

//GetEndpoints ...
func (sh StudyHandler) GetEndpointList() api.EndpointList {

	return api.EndpointList{
		&api.Endpoint{Path: "studies", Methods: []string{http.MethodPost}}: sh.GetStudies,
	}

}

func (sh StudyHandler) GetStudies(w http.ResponseWriter, r *http.Request) {
	authService := auth.AuthenticationService(r)
	if authService == nil || !authService.IsValid() {

		util.JSONResponse(
			model.ErrUserInvalidBearer,
			w,
			http.StatusUnauthorized,
		)
		return
	}
	decoder := json.NewDecoder(r.Body)

	var studiesUIDs []string
	err := decoder.Decode(&studiesUIDs)
	if err != nil {
		panic(err)
	}
	var session *repository.Session
	var auth *idonia.Auth
	if idonia.Token == "" {
		session, err = repository.GetSession(sh.DB, authService.GetToken())
		if err != nil {
			fmt.Println(err)
		}
		idonia.Token = session.IdoniaToken
	}
	if session != nil {
		auth = &idonia.Auth{
			Token:     session.IdoniaToken,
			AccountID: sh.Conf.Idonia.ApiId,
		}
	} else {
		auth = &idonia.Auth{
			Token:     idonia.Token,
			AccountID: sh.Conf.Idonia.ApiId,
		}
	}
	fmt.Println(idonia.Token)
	studiesInIdonia, err := study.RetriveAllStudies(sh.DB, sh.Logger, sh.Conf, auth)
	var studies []repository.Study
	for _, aStudy := range studiesInIdonia {
		for _, studyUIDAtPacs := range studiesUIDs {
			if aStudy.StudyUID == studyUIDAtPacs {
				studies = append(studies, *aStudy)
			}
		}

	}
	if len(studies) == 0 {
		util.JSONResponse(
			"No studies at the time",
			w,
			http.StatusNoContent,
		)
	}
	util.JSONResponse(
		studies,
		w,
		http.StatusOK,
	)
}
