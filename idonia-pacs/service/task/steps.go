package task

import (
	"bitbucket.org/inehealth/idonia-pacs/configuration"
	"bitbucket.org/inehealth/idonia-pacs/model"
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/gcs"
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia-auth"
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia-core"
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia-share"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var UserAgent = "Idonia PACS"

type Step = func(task *Task, config configuration.Configuration) (err error)

var Steps = map[string]Step{
	"CreateStudy":          CreateStudy,
	"CreateDICOMContainer": CreateDICOMContainer,
	"TransferFile":         TransferFile,
	"ShareFileWithAccount": ShareFileWithAccount,
	"ShareFileWithGroup":   ShareFileWithGroup,
	"ShareFileWithEmail":   ShareFileWithEmail,
	"GetMagicLink":         GetMagicLink,
	"UploadAllDICOM":       UploadAllDICOM,
	"UploadPendingDICOM":   UploadPendingDICOM,
	"ReportStudy":          ReportStudy,
}

const (
	CreateStudyID = 1 + iota
	CreateDICOMContainerID
	UploadAllDICOMID
)

var IdoniaSteps = map[uint32]string{
	CreateStudyID:          "CreateStudy",
	CreateDICOMContainerID: "CreateDICOMContainer",
	UploadAllDICOMID:       "UploadAllDICOM",
}

func CreateStudy(task *Task, config configuration.Configuration) (err error) {
	req := idonia_core.PostStudyReq{}
	req.Name = &task.StudyUID
	if val, ok := task.AdditionalFields["StudyName"]; ok {
		req.Name = &val
	}
	if val, ok := task.AdditionalFields["Description"]; ok {
		req.Description = &val
	}
	if val, ok := task.AdditionalFields["ParentID"]; ok {
		num, _ := strconv.ParseInt(val, 10, 64)
		numInt := int(num)
		req.ParentID = &numInt
	}
	if val, ok := task.AdditionalFields["Color"]; ok {
		req.Color = &val
	}

	res, err := idonia_core.PostStudy(&req, task.Authorization)

	if err != nil {
		task.Error = model.ErrIdonia
		task.Error.AdditionalMessage = err.Error()
		task.Status = Error
		return
	}
	task.StudyID = strconv.Itoa(res.StudyID)
	task.AdditionalFields["StudyID"] = strconv.Itoa(res.StudyID)

	return
}

func CreateDICOMContainer(task *Task, config configuration.Configuration) (err error) {
	if len(task.DICOMFiles) == 0 {
		return
	}
	req := idonia_core.PostContainerReq{}
	req.Name = "Imágenes de"

	if val, ok := task.AdditionalFields["ContainerName"]; ok {
		req.Name = val
	}
	if val, ok := task.AdditionalFields["StudyID"]; ok {
		num, _ := strconv.ParseInt(val, 10, 64)
		numInt := int(num)
		req.ParentID = &numInt
	}
	if req.ParentID == nil {
		num, _ := strconv.Atoi(task.StudyID)
		req.ParentID = &num
	}
	fmt.Println("SE LLAMA ESSTE PASO CONTAINERO QUE", req.ParentID)
	res, err := idonia_core.PostContainer(&req, task.Authorization)

	if err != nil {
		task.Error = model.ErrIdonia
		task.Error.AdditionalMessage = err.Error()
		task.Status = Error
		return
	}
	task.ContainerID = uint32(res.ContainerID)
	task.AdditionalFields["FileID"] = strconv.Itoa(res.FileID)
	task.AdditionalFields["ContainerID"] = strconv.Itoa(res.ContainerID)
	return
}

func TransferFile(task *Task, config configuration.Configuration) (err error) {
	req := idonia_share.PostTransferFile_FileIDReq{}

	var FileID int

	if task.StudyID != "" {
		FileID, _ = strconv.Atoi(task.StudyID)

	}
	if FileID == 0 {
		if val, ok := task.AdditionalFields["StudyID"]; ok {
			FileID, err = strconv.Atoi(val)
			if err != nil || FileID == 0 {
				return errors.New("FileID cant be equal 0")
			}
		}
	}
	if val, ok := task.AdditionalFields["Email"]; ok {
		req.Email = val
	}
	if val, ok := task.AdditionalFields["Phone"]; ok {
		req.Phone = val
	}
	fmt.Println(FileID)
	res, err := idonia_share.PostTransferFile_FileID(&req, FileID, task.Authorization)

	if err != nil {
		task.Error = model.ErrIdonia
		task.Error.AdditionalMessage = err.Error()
		task.Status = Error
		return
	}
	task.AdditionalFields["Token"] = res.Token

	return
}

func ShareFileWithAccount(task *Task, config configuration.Configuration) (err error) {
	req := idonia_share.PutShareFile_FileIDReq{}

	var FileID int
	if task.StudyID != "" {
		FileID, _ = strconv.Atoi(task.StudyID)

	}
	if FileID == 0 {
		if val, ok := task.AdditionalFields["StudyID"]; ok {
			FileID, err = strconv.Atoi(val)
			if err != nil || FileID == 0 {
				return errors.New("FileID cant be equal 0")
			}
		}
	}
	if val, ok := task.AdditionalFields["AccountID"]; ok {
		AccountID, _ := strconv.ParseInt(val, 10, 64)
		req.Accounts = append(req.Accounts, int(AccountID))
	}
	if val, ok := task.AdditionalFields["Permission"]; ok {
		req.Permission = val
	}
	if val, ok := task.AdditionalFields["ExpiredAt"]; ok {
		expiredAtString := val
		expiredAt, _ := time.Parse(time.RFC3339, expiredAtString)
		req.ExpiredAt = &expiredAt
	}
	if val, ok := task.AdditionalFields["Message"]; ok {
		req.Message = val
	}

	_, err = idonia_share.PutShareFile_FileID(&req, int(FileID), task.Authorization)

	if err != nil {
		task.Error = model.ErrIdonia
		task.Error.AdditionalMessage = err.Error()
		task.Status = Error
		return
	}

	return
}

func ShareFileWithGroup(task *Task, config configuration.Configuration) (err error) {
	req := idonia_share.PutShareFile_FileIDReq{}

	var FileID int
	if task.StudyID != "" {
		FileID, _ = strconv.Atoi(task.StudyID)

	}
	if FileID == 0 {
		if val, ok := task.AdditionalFields["StudyID"]; ok {
			FileID, err = strconv.Atoi(val)
			if err != nil || FileID == 0 {
				return errors.New("FileID cant be equal 0")
			}
		}
	}
	if val, ok := task.AdditionalFields["GroupID"]; ok {
		GroupID, _ := strconv.ParseInt(val, 10, 64)
		req.Groups = append(req.Groups, int(GroupID))
	}
	if val, ok := task.AdditionalFields["Permission"]; ok {
		req.Permission = val
	}
	if val, ok := task.AdditionalFields["ExpiredAt"]; ok {
		expiredAtString := val
		expiredAt, _ := time.Parse(time.RFC3339, expiredAtString)
		req.ExpiredAt = &expiredAt
	}
	if val, ok := task.AdditionalFields["Message"]; ok {
		req.Message = val
	}

	_, err = idonia_share.PutShareFile_FileID(&req, int(FileID), task.Authorization)

	if err != nil {
		task.Error = model.ErrIdonia
		task.Error.AdditionalMessage = err.Error()
		task.Status = Error
		return
	}

	return
}

func ShareFileWithEmail(task *Task, config configuration.Configuration) (err error) {
	req := idonia_share.PutShareFile_FileIDReq{}

	var FileID int
	if task.StudyID != "" {
		FileID, _ = strconv.Atoi(task.StudyID)

	}
	if FileID == 0 {
		if val, ok := task.AdditionalFields["StudyID"]; ok {
			FileID, err = strconv.Atoi(val)
			if err != nil || FileID == 0 {
				return errors.New("FileID cant be equal 0")
			}
		}
	}

	if val, ok := task.AdditionalFields["Email"]; ok {
		req.Guests = append(req.Guests, val)
	}
	if val, ok := task.AdditionalFields["Permission"]; ok {
		req.Permission = val
	}
	if val, ok := task.AdditionalFields["ExpiredAt"]; ok {
		expiredAtString := val
		expiredAt, _ := time.Parse(time.RFC3339, expiredAtString)
		req.ExpiredAt = &expiredAt
	}
	if val, ok := task.AdditionalFields["Message"]; ok {
		req.Message = val
	}
	_, err = idonia_share.PutShareFile_FileID(&req, FileID, task.Authorization)

	if err != nil {
		task.Error = model.ErrIdonia
		task.Error.AdditionalMessage = err.Error()
		task.Status = Error
		return
	}

	return
}

func GetMagicLink(task *Task, config configuration.Configuration) (err error) {
	req := idonia_share.PostPublicURLReq{}

	defaultNamespace := config.Idonia.Namespace

	accRes, _ := idonia_auth.GetAccount(task.Authorization)

	if accRes != nil && accRes.Account.Username != "" {
		defaultNamespace = accRes.Account.Username
	}
	var id int
	req.Namespace = &defaultNamespace
	if task.StudyID != "" {
		id, _ = strconv.Atoi(task.StudyID)

	}
	if id == 0 {
		if val, ok := task.AdditionalFields["StudyID"]; ok {
			id, err = strconv.Atoi(val)
			if err != nil || id == 0 {
				return errors.New("ID cant be equal 0")
			}
		}
	}
	req.FileID = id

	if val, ok := task.AdditionalFields["Namespace"]; ok {
		reg, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			log.Fatal(err)
		}
		filtered := reg.ReplaceAllString(val, "")
		req.Namespace = &filtered
	}
	fmt.Println(fmt.Sprintf("%+v"), req)
	res, err := idonia_share.PostPublicURL(&req, task.Authorization)

	if err != nil {
		task.Error = model.ErrIdonia
		task.Error.AdditionalMessage = err.Error()
		task.Status = Error
		return
	}
	task.AdditionalFields["PublicURL"]= config.Idonia.ApiHost + "/v/" + res.Namespace +"#"+ res.SecureID
	task.AdditionalFields["PublicURLID"] = strconv.Itoa(res.PublicURLID)
	task.AdditionalFields["Namespace"] = res.Namespace
	task.AdditionalFields["SecureID"] = res.SecureID
	return
}

func UploadAllDICOM(task *Task, config configuration.Configuration) (err error) {
	fmt.Println("UPLOADALLDICOM", fmt.Sprintf("%+v", task.DICOMFiles), task.ContainerID)
	for i, dicom := range task.DICOMFiles {
		dicomURL := GetDICOMURL(task.DICOMFiles[i], config)
		response, err := http.Get(dicomURL)
		if err != nil {
			fmt.Println(fmt.Sprintf("%+v", response), dicomURL, err.Error())
			continue
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(fmt.Sprintf("%+v", response), dicomURL, err.Error())
			return err
		}
		fmt.Println(dicomURL)
		task.DICOMFiles[i].Bytes = body
		task.DICOMFiles[i].Status = Retrieved
		policy, err := idonia_core.GetPolicy(task.Authorization)
		if err != nil {
			task.Error = model.ErrIdonia
			task.Error.AdditionalMessage = err.Error()
			task.Status = Error
			return err
		}
		err = gcs.UploadFile(body, dicom.SOPInstanceUID, policy)
		if err != nil {
			task.Error = model.ErrIdonia
			task.Error.AdditionalMessage = err.Error()
			task.Status = Error
			return err
		}

		_, err = idonia_core.PostContainers_ContainerIDDICOM(&idonia_core.PostContainers_ContainerIDDICOMReq{
			UUID: &policy.UUID,
		}, int(task.ContainerID), task.Authorization)
		if err != nil {
			task.Error = model.ErrIdonia
			task.Error.AdditionalMessage = err.Error()
			task.Status = Error
			return err
		}
		task.DICOMFiles[i].Status = Completed
		task.DICOMFiles[i].IsPending = false
	}

	return
}

func UploadPendingDICOM(task *Task, config configuration.Configuration) (err error) {

	for i, dicom := range task.DICOMFiles {
		if dicom.IsPending {
			dicomURL := GetDICOMURL(task.DICOMFiles[i], config)
			response, err := http.Get(dicomURL)
			if err != nil {
				continue
			}
			body, err := ioutil.ReadAll(response.Body)
			task.DICOMFiles[i].Bytes = body
			task.DICOMFiles[i].Status = Retrieved
			policy, err := idonia_core.GetPolicy(task.Authorization)
			if err != nil {
				task.Error = model.ErrIdonia
				task.Error.AdditionalMessage = err.Error()
				task.Status = Error
				return err
			}
			err = gcs.UploadFile(body, dicom.SOPInstanceUID, policy)
			if err != nil {
				task.Error = model.ErrIdonia
				task.Error.AdditionalMessage = err.Error()
				task.Status = Error
				return err
			}

			_, err = idonia_core.PostContainers_ContainerIDDICOM(&idonia_core.PostContainers_ContainerIDDICOMReq{
				UUID: &policy.UUID,
			}, int(task.ContainerID), task.Authorization)
			if err != nil {
				task.Error = model.ErrIdonia
				task.Error.AdditionalMessage = err.Error()
				task.Status = Error
				return err
			}
			task.DICOMFiles[i].Status = Completed
			task.DICOMFiles[i].IsPending = false
		}

	}

	return
}

func ReportStudy(task *Task, config configuration.Configuration) (err error) {
	var req idonia_core.PostDocumentReq
	defaultName := "Informe"
	req.FileName = &defaultName

	if val, ok := task.AdditionalFields["ReportName"]; ok {
		req.FileName = &val
	}
	if task.StudyID != "" {
		id, err := strconv.Atoi(task.StudyID)
		if err != nil {
			if val, ok := task.AdditionalFields["StudyID"]; ok {
				id, err = strconv.Atoi(val)
				if err != nil {
					return err
				}
			}
		}
		req.ParentID = &id
	}
	policy, err := idonia_core.GetPolicy(task.Authorization)
	if err != nil {
		task.Error = model.ErrIdonia
		task.Error.AdditionalMessage = err.Error()
		task.Status = Error
		return err
	}

	err = gcs.UploadFile(task.File, *req.FileName, policy)
	if err != nil {
		task.Error = model.ErrIdonia
		task.Error.AdditionalMessage = err.Error()
		task.Status = Error
		return err
	}
	req.UUID = &policy.UUID
	true := true
	req.IsReport = &true
	_, err = idonia_core.PostDocument(&req, task.Authorization)
	if err != nil {
		task.Error = model.ErrIdonia
		task.Error.AdditionalMessage = err.Error()
		task.Status = Error
		return err
	}
	task.IsReported = true
	return
}
