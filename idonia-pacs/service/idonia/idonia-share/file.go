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
	"time"
)

type PutShareFile_FileIDReq struct {
	Accounts   []int      `json:"accounts"`
	Groups     []int      `json:"groups"`
	Guests     []string   `json:"guests"`
	ExpiredAt  *time.Time `json:"expired_at"`
	Permission string     `json:"permission"`
	Message    string     `json:"message,omitempty"`
}

type PutShareFile_FileIDRes struct{}

func PutShareFile_FileID(body *PutShareFile_FileIDReq, fileID int, auth *idonia.Auth) (res *PutShareFile_FileIDRes, err error) {
	u, err := url.Parse(idonia.APIHost + fmt.Sprintf("/share/file/%d", fileID))
	if err != nil {
		return
	}
	encodedBody, err := json.Marshal(body)
	if err != nil {
		return
	}
	response, err := idonia.Put(u, bytes.NewReader(encodedBody), auth)
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
	//err = json.Unmarshal(bodyString, &res)
	return &PutShareFile_FileIDRes{}, nil
}
