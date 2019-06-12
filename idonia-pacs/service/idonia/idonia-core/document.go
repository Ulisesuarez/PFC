package idonia_core

import (
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type GetPolicyRes struct {
	Policy      string `json:"policy"`
	Signature   string `json:"signature"`
	UUID        string `json:"uuid"`
	URL         string `json:"url"`
	Bucket      string `json:"bucket"`
	ClientEmail string `json:"client_email"`
}

func GetPolicy(auth *idonia.Auth) (res *GetPolicyRes, err error) {
	u, err := url.Parse(idonia.APIHost + "/document/policy")
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
	if response.StatusCode != http.StatusCreated {
		err = errors.New("Unexpected HTTP Status: " + strconv.Itoa(response.StatusCode) + ", Server said: " + string(bodyString))
		return
	}
	err = json.Unmarshal(bodyString, &res)
	return
}

type PostDocumentReq struct {
	FileName *string `json:"file_name"`
	UUID     *string `json:"uuid"`
	ParentID *int    `json:"parent_id"`
	IsReport *bool   `json:"description"`
}

type PostDocumentRes struct {
	FileID int    `json:"id_file"`
	TaskID string `json:"id_task"`
}

func PostDocument(body *PostDocumentReq, auth *idonia.Auth) (res *PostDocumentRes, err error) {
	parameters := url.Values{}
	if body.FileName != nil {
		parameters.Add("file_name", *body.FileName)
	}
	if body.UUID != nil {
		parameters.Add("uuid", *body.UUID)
	}
	if body.ParentID != nil {
		parameters.Add("parent_id", fmt.Sprintf("%d", *body.ParentID))
	}
	if body.IsReport != nil {
		parameters.Add("is_report", fmt.Sprintf("%t", *body.IsReport))
	}
	u, err := url.Parse(idonia.APIHost + fmt.Sprintf("/document"))
	u.RawQuery = parameters.Encode()
	if err != nil {
		return
	}

	response, err := idonia.Post(u, nil, auth)
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

type GetDocumentReq struct {
	FileID string
}

type GetDocumentRes struct {
	IsDeleted bool `json:"deleted"`
	IsReport  bool `json:"is_report"`
}

func GetDocument(req *GetDocumentReq, auth *idonia.Auth) (res *GetDocumentRes, err error) {
	u, err := url.Parse(idonia.APIHost + "/document/" + req.FileID)
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
