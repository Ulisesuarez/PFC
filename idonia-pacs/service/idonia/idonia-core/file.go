package idonia_core

import (
	"bitbucket.org/inehealth/idonia-core/model"
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type PostParent_ParentIDChild_ChildIDReq struct{}

type PostParent_ParentIDChild_ChildIDRes struct{}

func PostParent_ParentIDChild_ChildID(body *PostParent_ParentIDChild_ChildIDReq, parentID int, childID int, auth *idonia.Auth) (res *PostParent_ParentIDChild_ChildIDRes, err error) {
	u, err := url.Parse(idonia.APIHost + fmt.Sprintf("/parent/%d/child/%d", parentID, childID))
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

func GetFile_FileIDShared(fileID string, auth *idonia.Auth) (res *GetFile_FileIDSharedRes, err error) {
	u, err := url.Parse(idonia.APIHost + "/file/" + fileID + "/shared")
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
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusConflict {
		err = errors.New("Unexpected HTTP Status: " + strconv.Itoa(response.StatusCode) + ", Server said: " + string(bodyString))
		return
	}
	err = json.Unmarshal(bodyString, &res)
	return
}

type GetFile_FileIDSharedRes struct {
	Accounts  []model.AccountPermission `json:"accounts"`
	Groups    []model.Group             `json:"groups"`
	Guests    []model.GuestPermission   `json:"guests"`
	PublicURL struct {
		PublicURLID int    `json:"public_url_id" model:"public_url_id"`
		Namespace   string `json:"namespace" model:"namespace"`
		SecureID    string `json:"secure_id" model:"secure_id"`
	} `json:"public_url"`
}
