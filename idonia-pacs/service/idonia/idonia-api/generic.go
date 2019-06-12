package idonia_api

import (
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type GetVersionRes struct {
	Version    string `json:"version"`
	ServerTime string `json:"server_time"`
}

func GetVersion(auth *idonia.Auth) (res *GetVersionRes, err error) {
	u, err := url.Parse(idonia.APIHost + "/version")
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
		err = errors.New("Unexpected HTTP Status: " + strconv.Itoa(response.StatusCode) + ", Server said: " + string(bodyString))
		return
	}
	err = json.Unmarshal(bodyString, &res)
	return
}
