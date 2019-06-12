package task

import (
	"bitbucket.org/inehealth/idonia-pacs/configuration"
	"bitbucket.org/inehealth/idonia-pacs/model"
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia"
	"strconv"
)

type Status string

const (
	Error      Status = "error"
	New               = "new"
	Retrieving        = "retrieving"
	Retrieved         = "retrieved"
	Sending           = "sending"
	Sent              = "sent"
	Completed         = "completed"
	Completing        = "completing"
)

func NewStatus(string string) (status Status) {
	return Status(string)
}

type File struct {
	UUID     string      `json:"uuid"`
	Name     string      `json:"name"`
	Path     string      `json:"-"`
	IsReport bool        `json:"isReport"`
	Status   Status      `json:"status"`
	Error    model.Error `json:"error"`
}

type ActionState struct {
	UUID            string      `json:"uuid"`
	IdoniaAccountID uint32      `json:"-"`
	IdoniaToken     string      `json:"-"`
	Step            uint32      `json:"step"`
	ActionID        string      `json:"action_id"`
	Status          Status      `json:"status"`
	Error           model.Error `json:"error"`
}

type Task struct {
	UUID             string            `json:"uuid"`
	Step             uint32            `json:"step"`
	Steps            []string          `json:"steps"`
	Status           Status            `json:"status"`
	Error            model.Error       `json:"error"`
	StudyUID         string            `json:"studyUID"`
	DICOMFiles       []DICOMFile       `json:"dicomFiles"`
	File             []byte            `json:"file"`
	AdditionalFields map[string]string `json:"additionalFields"`
	Authorization    *idonia.Auth      `json:"authorization"`
	StudyID          string            `json:"studyID"`
	ContainerID      uint32            `json:"containerID"`
	IsReported       bool              `json:"isReported"`
}

type DICOMFile struct {
	UUID              string      `json:"uuid"`
	Bytes             []byte      `json:"-"`
	StudyUID          string      `json:"studyUID"`
	ContainerID       uint32      `json:"containerID"`
	SeriesInstanceUID string      `json:"seriesInstanceUID"`
	SOPInstanceUID    string      `json:"sopInstanceUID"`
	Status            Status      `json:"status"`
	Error             model.Error `json:"error"`
	IsPending         bool        `json:"isPending"`
}

func (task *Task) RunNextStep(config configuration.Configuration) (noMoreSteps bool, err error) {
	if int(task.Step) >= len(task.Steps) {
		return true, nil
	}
	stepDescription := task.Steps[task.Step]
	step, ok := Steps[stepDescription]
	if !ok {
		err = model.ErrConfStepTypeMistyped
		task.Error = err.(model.Error)
		task.Status = Error
		return
	}
	err = step(task, config)
	if err != nil {
		return
	}

	if int(task.Step) == len(task.Steps) {
		task.Status = Sent
		return
	}

	task.Step++
	return
}

func (task *Task) HasStep(step string) bool {
	for _, s := range task.Steps {
		if step == s {
			return true
		}
	}
	return false
}

func (task *Task) AddUploadAllToIdoniaSteps() {
	steps := []string{}
	steps = append(steps, "CreateStudy")
	steps = append(steps, "CreateDICOMContainer")
	steps = append(steps, "UploadAllDICOM")
	if !task.HasStep("CreateStudy") {
		task.Steps = append(steps, task.Steps...)
	}
}

func (task *Task) markPendingDICOM(oldDicom []DICOMFile) {
	for i, taskDicom := range task.DICOMFiles {
		found := false
		for _, dicom := range oldDicom {
			if dicom.SOPInstanceUID == taskDicom.SOPInstanceUID {
				found = true
			}
		}
		if !found {
			task.DICOMFiles[i].IsPending = true
		}
	}
}

func (task *Task) AddUploadPendingStep() {
	steps := []string{}
	steps = append(steps, "UploadPendingDICOM")
}

func (task *Task) AddAditionalFields(rq model.AddTaskRQ) {
	task.AdditionalFields = make(map[string]string)
	task.AdditionalFields["StudyName"] = rq.StudyName
	task.AdditionalFields["ContainerName"] = rq.ContainerName
	if rq.TransferRQ != nil {
		task.AdditionalFields["Phone"] = rq.TransferRQ.Phone
		task.AdditionalFields["Email"] = rq.TransferRQ.Email
	}
	if rq.ShareRQ != nil {
		if rq.ShareRQ.Account > 0 {
			accountID := strconv.Itoa(rq.ShareRQ.Account)
			task.AdditionalFields["AccountID"] = accountID
		}
		if rq.ShareRQ.Group > 0 {
			groupID := strconv.Itoa(rq.ShareRQ.Group)
			task.AdditionalFields["GroupID"] = groupID
		}
		if rq.ShareRQ.Guests != "" {
			task.AdditionalFields["Emial"] = rq.ShareRQ.Guests
		}
		task.AdditionalFields["Permission"] = rq.ShareRQ.Permission
		task.AdditionalFields["Message"] = rq.ShareRQ.Message
		if rq.ShareRQ.ExpiredAt != nil {
			task.AdditionalFields["ExpiredAt"] = rq.ShareRQ.ExpiredAt.String()
		}
	}
	if rq.MagicLinkRQ != nil {
		task.AdditionalFields["Namespace"] = rq.MagicLinkRQ.Namespace
	}
}
