package controller

import (
	"bytes"
	"net/http"
	"time"

	"bitbucket.org/inehealth/api/api"
	"bitbucket.org/inehealth/api/util"
)

//GenericModule ...
type GenericModule struct {
	Version     string
	Compilation string
}

//GetEndpoints ...
func (gh GenericModule) GetEndpointList() api.EndpointList {

	return api.EndpointList{
		&api.Endpoint{Path: "crossdomain.xml", Methods: []string{http.MethodGet, http.MethodPost}}: gh.GetCrossDomainPolicy,
		&api.Endpoint{Path: "status", Methods: []string{http.MethodGet}}:                           gh.GetAPIStatus,
	}

}

//GetCrossDomainPolicy Controller
func (gh GenericModule) GetCrossDomainPolicy(w http.ResponseWriter, r *http.Request) {
	var cdBuffer bytes.Buffer

	cdBuffer.WriteString("<?xml version=\"1.0\"?>")
	cdBuffer.WriteString("<!DOCTYPE cross-domain-policy SYSTEM \"http://www.macromedia.com/xml/dtds/cross-domain-policy.dtd\">")
	cdBuffer.WriteString("<cross-domain-policy>")
	cdBuffer.WriteString("<site-control permitted-cross-domain-policies=\"all\"/>")
	cdBuffer.WriteString("<allow-access-from domain=\"*\" to-ports=\"*\" secure=\"false\"/>")
	cdBuffer.WriteString("<allow-http-request-headers-from domain=\"*\" headers=\"*\"/>")
	cdBuffer.WriteString("</cross-domain-policy>")

	w.Header().Set("Content-Type", "application/xml")
	w.Write(cdBuffer.Bytes())
}

//GetAPIStatus ...
func (gh GenericModule) GetAPIStatus(w http.ResponseWriter, r *http.Request) {
	util.JSONResponse(
		struct {
			Version     string    `json:"version"`
			Compilation string    `json:"compilation"`
			ServerTime  time.Time `json:"serverTime"`
		}{
			Version:     gh.Version,
			Compilation: gh.Compilation,
			ServerTime:  time.Now(),
		},
		w,
		http.StatusOK,
	)
}
