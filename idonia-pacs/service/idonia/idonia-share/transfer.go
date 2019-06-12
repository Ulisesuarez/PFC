package idonia_share

import (
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

type PostTransferFile_FileIDReq struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type PostTransferFile_FileIDRes struct {
	Token string `json:"token"`
}

func PostTransferFile_FileID(body *PostTransferFile_FileIDReq, fileID int, auth *idonia.Auth) (res *PostTransferFile_FileIDRes, err error) {
	u, err := url.Parse(idonia.APIHost + fmt.Sprintf("/transfer/file/%d", fileID))
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
	if response.StatusCode != http.StatusOK {
		err = errors.New("Unexpected HTTP Status: " + strconv.Itoa(response.StatusCode) + ", Server said: " + string(bodyString))
		return
	}
	err = json.Unmarshal(bodyString, &res)
	return
}
