package idonia_core

import (
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type PostDirectoryReq struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ParentID    *int    `json:"parent_id"`
}

type PostDirectoryRes struct {
	DirectoryID int `json:"directory_id"`
}

func PostDirectory(body *PostDirectoryReq, auth *idonia.Auth) (res *PostDirectoryRes, err error) {
	u, err := url.Parse(idonia.APIHost + "/directory")
	if err != nil {
		return
	}
	encodedBody, err := json.Marshal(body)
	if err != nil {
		return
	}
	response, err := idonia.Post(u, bytes.NewReader(encodedBody), auth)
	if err != nil {
		return
	}
	defer response.Body.Close()
	var bodyString []byte
	bodyString, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusCreated {
		err = errors.New("Unexpected HTTP Status: " + strconv.Itoa(response.StatusCode) + ", Server said: " + string(bodyString))
		return
	}
	err = json.Unmarshal(bodyString, &res)
	return
}
