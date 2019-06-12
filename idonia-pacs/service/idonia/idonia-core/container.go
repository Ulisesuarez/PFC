package idonia_core

import (
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

type PostContainerReq struct {
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id"`
}

type PostContainerRes struct {
	FileID      int `json:"file_id"`
	ContainerID int `json:"container_id"`
}

func PostContainer(body *PostContainerReq, auth *idonia.Auth) (res *PostContainerRes, err error) {
	u, err := url.Parse(idonia.APIHost + "/containers")
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

type PostContainers_ContainerIDDICOMReq struct {
	UUID *string `json:"uuid"`
}

type PostContainers_ContainerIDDICOMRes struct {
	TaskID string `json:"id_task"`
}

func PostContainers_ContainerIDDICOM(body *PostContainers_ContainerIDDICOMReq, ContainerID int, auth *idonia.Auth) (res *PostContainers_ContainerIDDICOMRes, err error) {
	if body.UUID == nil {
		return nil, errors.New("Bad Param")
	}

	b := new(bytes.Buffer)
	writer := multipart.NewWriter(b)
	defer writer.Close()

	writer.WriteField("uuid", *body.UUID)
	writer.Close()

	u, err := url.Parse(idonia.APIHost + fmt.Sprintf("/containers/%d/dicom", ContainerID))
	if err != nil {
		return
	}

	response, err := idonia.PostWithContentType(u, b, auth, writer.FormDataContentType())

	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return
	}
	var bodyString []byte
	bodyString, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusCreated {
		err = errors.New(string(bodyString))
		return
	}
	err = json.Unmarshal(bodyString, &res)
	return
}

type GetContainerReq struct {
	ContainerID *uint32
}

type GetContainerRes struct {
	ID            int  `json:"id"`
	FileID        int  `json:"file_id"`
	StudyID       *int `json:"study_id"`
	ReportID      *int `json:"report_id"`
	DicomPatients []struct {
		DicomStudies []struct {
			ID          int `json:"id"`
			DicomSeries []struct {
				ID             int `json:"id"`
				InstancesCount int `json:"instances_count"`
			} `json:"dicom_series"`
		} `json:"dicom_studies"`
	} `json:"dicom_patients"`
	MediaInstances interface{} `json:"media_instances"`
}

func GetContainer(body *GetContainerReq, auth *idonia.Auth) (res *GetContainerRes, err error) {
	if body.ContainerID == nil {
		return nil, errors.New("Bad Param")
	}

	u, err := url.Parse(idonia.APIHost + fmt.Sprintf("/containers/%d/", *body.ContainerID))
	if err != nil {
		return
	}

	response, err := idonia.Get(u, auth)

	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return
	}
	var bodyString []byte
	bodyString, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		err = errors.New(string(bodyString))
		return
	}
	err = json.Unmarshal(bodyString, &res)
	return
}
