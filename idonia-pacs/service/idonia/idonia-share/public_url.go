package idonia_share

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

type PostPublicURLReq struct {
	Namespace *string `json:"namespace,omitempty"`
	FileID    int     `json:"file_id"`
}

type PostPublicURLRes struct {
	PublicURLID int    `json:"public_url_id"`
	Namespace   string `json:"namespace"`
	SecureID    string `json:"secure_id"`
}

func PostPublicURL(body *PostPublicURLReq, auth *idonia.Auth) (res *PostPublicURLRes, err error) {
	u, err := url.Parse(idonia.APIHost + "/public-url")
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
	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusConflict {
		err = errors.New("Unexpected HTTP Status: " + strconv.Itoa(response.StatusCode) + ", Server said: " + string(bodyString))
		return
	}
	err = json.Unmarshal(bodyString, &res)
	return
}
