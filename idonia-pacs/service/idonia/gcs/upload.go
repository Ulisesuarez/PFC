package gcs

import (
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia-core"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

func UploadFile(file []byte, name string, policy *idonia_core.GetPolicyRes) (err error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	defer writer.Close()
	fmt.Println(fmt.Sprintf("%+v", policy))
	urlParts := strings.Split(policy.URL, "/")
	writer.WriteField("GoogleAccessId", policy.ClientEmail)
	writer.WriteField("key", "in/"+policy.UUID)
	writer.WriteField("bucket", urlParts[len(urlParts)-1])
	writer.WriteField("policy", policy.Policy)
	writer.WriteField("signature", policy.Signature)
	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		return err
	}
	r := bytes.NewReader(file)
	written, err := io.Copy(part, r)
	fmt.Println(written, err)
	// This writer must be closed now in order to add the boundary closing tag to the multipart
	writer.Close()

	httpClient := &http.Client{}

	req, err := http.NewRequest("POST", policy.URL, body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	response, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	var bodyString []byte
	bodyString, err = ioutil.ReadAll(response.Body)
	if response.StatusCode != http.StatusNoContent {
		err = errors.New("Unexpected HTTP Status: " + strconv.Itoa(response.StatusCode) + ", Server said: " + string(bodyString))
		return
	}

	return
}
