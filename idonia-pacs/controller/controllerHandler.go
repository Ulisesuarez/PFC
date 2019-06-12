package controller

import (
	"bitbucket.org/inehealth/api/api"
	"bitbucket.org/inehealth/api/util"
	"bitbucket.org/inehealth/idonia-pacs/configuration"
	"database/sql"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

//GenericModule ...
type ControllerHandler struct {
	Conf   *configuration.Configuration
	DB     *sql.DB
	Logger *logrus.Logger
}

//GetEndpoints ...
func (ch ControllerHandler) GetEndpointList() api.EndpointList {

	return api.EndpointList{
		&api.Endpoint{Path: "tuso", Methods: []string{http.MethodGet, http.MethodPost}}: ch.test,
		&api.Endpoint{Path: "statuso", Methods: []string{http.MethodGet}}:               ch.test,
	}

}

func (ch ControllerHandler) test(w http.ResponseWriter, r *http.Request) {
	util.JSONResponse(
		struct {
			Hola string    `json:"version"`
			Que  string    `json:"compilation"`
			Tal  time.Time `json:"serverTime"`
		}{
			Hola: "Holaa",
			Que:  "k tal",
			Tal:  time.Now(),
		},
		w,
		http.StatusOK,
	)
}
