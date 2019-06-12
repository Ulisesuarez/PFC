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
	"time"
)

type PostStudyReq struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ParentID    *int    `json:"parent_id"`
	Color       *string `json:"color"`
}

type PostStudyRes struct {
	StudyID int `json:"study_id"`
}

func PostStudy(body *PostStudyReq, auth *idonia.Auth) (res *PostStudyRes, err error) {
	u, err := url.Parse(idonia.APIHost + "/study")
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

type GetStudyReq struct {
	FileID string
}

type GetStudyRes struct {
	FileID       int32      `json:"id_file" model:"file_id"`
	IsDeleted    bool       `json:"deleted"`
	DeletedAt    *time.Time `json:"deleted_at"`
	LastReportID *int32     `json:"last_report_id"`
}

func GetStudy(req *GetStudyReq, auth *idonia.Auth) (res *GetStudyRes, err error) {
	u, err := url.Parse(idonia.APIHost + "/study/" + req.FileID)
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
