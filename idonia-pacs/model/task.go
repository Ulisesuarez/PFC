package model

import (
	"time"
)

type AddTaskRQ struct {
	StudyUID      string       `json:"studyUID"`
	StudyName     string       `json:"studyName"`
	ContainerName string       `json:"containerName"`
	SOPObjects    []SOPObject  `json:"sopObjects"`
	Steps         []string     `json:"steps"`
	File          string       `json:"file"`
	ShareRQ       *ShareRQ     `json:"shareRQ, omitempty"`
	TransferRQ    *TransferRQ  `json:"transferRQ, omitempty"`
	MagicLinkRQ   *PublicURLRQ `json:"magicLinkRQ, omitempty"`
}

type SOPObject struct {
	StudyUID          string `json:"studyUID"`
	SeriesInstanceUID string `json:"seriesInstanceUID"`
	SOPInstanceUID    string `json:"sopInstanceUID"`
}

type TaskRS struct {
	UUID            string    `json:"uuid"`
	StudyUID        string    `json:"studyUID"`
	Steps           StepRS    `json:"steps"`
	AditionalFields string    `json:"aditionalFields"`
	Status          string    `json:"status"`
	Error           error     `json:"error"`
	CreatedAt       time.Time `json:"createdAt"`
	ContainerID     uint32    `json:"containerID"`
	StudyID         string    `json:"studyID"`
	IsReported      bool      `json:"isReported"`
}

type TasksRS []*TaskRS

type StepRS struct {
	CurrentStepID string `json:"currentStepId"`
	CurrentStep   uint32 `json:"currentStep"`
	TotalSteps    uint32 `json:"totalSteps"`
}
