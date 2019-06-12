package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"bitbucket.org/inehealth/idonia-pacs/model"
)

func ReadBodyFromRequest(r *http.Request, target interface{}) (bodyString []byte, err error) {
	defer r.Body.Close()
	bodyString, err = ioutil.ReadAll(r.Body)
	return
}

//GetDataFromRequest ...
func GetDataFromRequest(r *http.Request, data interface{}) (e *model.Error) {

	if r.Body == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if errDecode := decoder.Decode(data); errDecode != nil {
		err := model.ErrUserInvalidData
		e = &err
		e.AdditionalMessage = errDecode.Error()
	}

	return
}
