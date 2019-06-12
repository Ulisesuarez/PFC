package controller

import (
	"bitbucket.org/inehealth/idonia-pacs/repository"
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia"
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia-auth"
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia-core"
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia-share"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"

	"bitbucket.org/inehealth/api/api"
	"bitbucket.org/inehealth/idonia-pacs/configuration"
	"bitbucket.org/inehealth/idonia-pacs/model"
	"bitbucket.org/inehealth/idonia-pacs/service/auth"
	"bitbucket.org/inehealth/idonia-pacs/service/task"
	"bitbucket.org/inehealth/idonia-pacs/util"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"regexp"
)

//TaskHandler ...
type TaskHandler struct {
	Conf   *configuration.Configuration `inject:"configuration"`
	DB     *sql.DB                      `inject:"db"`
	Logger *logrus.Logger               `inject:"logger"`
	Engine *task.Engine                 `inject:"engine"`
}

func parseOutput(value string, additionalFields map[string]string) (parsedValue string) {

	parsedValue = value

	rp := regexp.MustCompile(`\${([a-zA-Z_0-9]+)}`)
	parsedValue = rp.ReplaceAllStringFunc(parsedValue, func(s string) string {
		s = strings.TrimSuffix(strings.TrimPrefix(s, "${"), "}")
		if val, ok := additionalFields[s]; ok {
			return val
		}
		return parsedValue
	})

	return parsedValue
}

//GetStaticList
func (ch TaskHandler) GetStaticList() api.StaticList {
	return api.StaticList{}
}

//GetEndpoints ...
func (ch TaskHandler) GetEndpointList() api.EndpointList {

	return api.EndpointList{
		&api.Endpoint{Path: "tasks", Methods: []string{http.MethodGet}}:                                          ch.GetTasks,
		&api.Endpoint{Path: "task", Methods: []string{http.MethodPost}}:                                          ch.AddTask,
		&api.Endpoint{Path: "task/{id:.+}/frontend/public-url/{file_id:.+}", Methods: []string{http.MethodPost}}: ch.GetMagicLink,
	}

}

func (th TaskHandler) AddTask(w http.ResponseWriter, r *http.Request) {
	authService := auth.AuthenticationService(r)
	if authService == nil || !authService.IsValid() {
		util.JSONResponse(
			model.ErrUserInvalidBearer,
			w,
			http.StatusUnauthorized,
		)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.JSONResponse(
			model.ErrServer,
			w,
			http.StatusInternalServerError,
		)
	}
	var TaskRQ = model.AddTaskRQ{
		SOPObjects: []model.SOPObject{},
	}
	err = json.Unmarshal(body, &TaskRQ)
	if err != nil {
		util.JSONResponse(
			model.ErrServer,
			w,
			http.StatusInternalServerError,
		)
	}
	aTask := task.Task{}
	aTask.StudyUID = TaskRQ.StudyUID
	aTask.File = []byte(TaskRQ.File)
	aTask.AddAditionalFields(TaskRQ)
	fmt.Println("QUE PASASAAAAA", idonia.Token)
	var session *repository.Session
	if idonia.Token == "" {
		session, err = repository.GetSession(th.DB, authService.GetToken())
		if err != nil {
			fmt.Println(err)
			th.Logger.Error(err.Error())
		}
		idonia.Token = session.IdoniaToken
	}
	if session != nil {
		aTask.Authorization = &idonia.Auth{
			Token:     session.IdoniaToken,
			AccountID: th.Conf.Idonia.ApiId,
		}
	} else {
		aTask.Authorization = &idonia.Auth{
			Token:     idonia.Token,
			AccountID: th.Conf.Idonia.ApiId,
		}
	}
	fmt.Println(idonia.Token)

	aTask.DICOMFiles = []task.DICOMFile{}
	aTask.UUID = uuid.New().String()
	fmt.Println(fmt.Sprintf("%+v", TaskRQ))
	fmt.Println(fmt.Sprintf("%+v", TaskRQ.SOPObjects))
	th.Logger.Info(TaskRQ.SOPObjects)
	for _, dcm := range TaskRQ.SOPObjects {
		th.Logger.Info(dcm.StudyUID)
		var dicom = task.DICOMFile{
			UUID:              uuid.New().String(),
			IsPending:         true,
			Status:            task.New,
			StudyUID:          dcm.StudyUID,
			SeriesInstanceUID: dcm.SeriesInstanceUID,
			SOPInstanceUID:    dcm.SOPInstanceUID,
		}
		aTask.DICOMFiles = append(aTask.DICOMFiles, dicom)
	}
	aTask.Status = task.New
	aTask.Steps = []string{}
	uploadToIDonia := false
	for _, step := range TaskRQ.Steps {
		if step == "uploadToIdonia" {
			uploadToIDonia = true
			continue
		}
		aTask.Steps = append(aTask.Steps, step)
	}
	if uploadToIDonia {
		aTask.AddUploadAllToIdoniaSteps()
	}
	err = th.Engine.Add(&aTask)
	if err != nil {
		util.JSONResponse(
			model.ErrServer,
			w,
			http.StatusInternalServerError,
		)
	}
	util.JSONResponse(
		"OK",
		w,
		http.StatusOK,
	)
}

func (th TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	authService := auth.AuthenticationService(r)
	if authService == nil || !authService.IsValid() {
		util.JSONResponse(
			model.ErrUserInvalidBearer,
			w,
			http.StatusUnauthorized,
		)
		return
	}

	page, count, _ := util.GetPaginatorFromRq(r, task.Task{})

	tasks, total, err := repository.GetTasks(th.DB, authService.GetIdoniaAuth().AccountID, uint32(count), uint32((page-1)*count))
	if err != nil {
		errRet := model.ErrServerDB
		errRet.AdditionalMessage = err.Error()
		util.JSONResponse(
			errRet,
			w,
			http.StatusInternalServerError,
		)
		return
	}

	tasksRsp := model.TasksRS{}

	for _, t := range tasks {

		tasksAdditionalFields, _ := repository.GetTaskAdditionalFields(th.DB, t.UUID)
		additionalFields := make(map[string]string)
		for _, v := range tasksAdditionalFields {
			additionalFields[v.Key] = v.Value
		}

		additionalFieldsByte, err := json.Marshal(additionalFields)
		if err != nil {
			th.Logger.Println(err)
		}
		var steps []string
		err = json.Unmarshal([]byte(t.Steps), &steps)
		if err != nil {
			fmt.Println(err)
		}
		stepName := ""

		if int(t.Step) < len(steps) {
			stepName = steps[t.Step]
		}

		taskRsp := model.TaskRS{
			UUID:            t.UUID,
			AditionalFields: string(additionalFieldsByte),
			Steps: model.StepRS{
				CurrentStepID: stepName,
				CurrentStep:   t.Step,
				TotalSteps:    uint32(len(steps)),
			},
			Status:      t.Status,
			Error:       t.Error,
			IsReported:  t.IsReported,
			ContainerID: t.ContainerID,
			StudyUID:    t.StudyUID,
			CreatedAt:   t.CreatedAt,
		}

		tasksRsp = append(tasksRsp, &taskRsp)

	}

	util.JSONResponse(
		model.PaginatedRS{
			Items: tasksRsp,
			Pagination: model.Pagination{
				Total: total,
				Page:  page,
				Count: count,
			},
		},
		w,
		http.StatusOK,
	)
}

func (ch TaskHandler) GetMagicLink(w http.ResponseWriter, r *http.Request) {
	authService := auth.AuthenticationService(r)

	if authService == nil || !authService.IsValid() {
		util.JSONResponse(
			model.ErrUserInvalidBearer,
			w,
			http.StatusUnauthorized,
		)
		return
	}
	taskId, _ := mux.Vars(r)["id"]
	aTask, err := repository.GetTask(ch.DB, taskId)
	if err != nil {
		util.JSONResponse(
			model.ErrServerDB,
			w,
			http.StatusInternalServerError,
		)
		return
	}
	fileID, _ := mux.Vars(r)["file_id"]
	var authI idonia.Auth
	err = json.Unmarshal([]byte(aTask.Authorization), &authI)
	response, err := idonia_core.GetFile_FileIDShared(fileID, &authI)
	if err != nil {
		println(err.Error())
	}
	if response.PublicURL.PublicURLID == 0 {
		//TODO println("0!!!!")
	}
	if err != nil || response.PublicURL.PublicURLID == 0 {
		type databaseError struct {
			code              int
			message           string
			additionalMessage string
		}
		var dbe databaseError
		if err != nil {
			ch.Logger.Error(err.Error())
			err = json.Unmarshal([]byte(err.Error()), &dbe)
			if dbe.additionalMessage != "sql: no rows in result set" {
				err = errors.New("Internal Server Error")
			}
		}
		if err == nil {
			ch.Logger.Error("sql: no rows in result set")
			req := idonia_share.PostPublicURLReq{}

			defaultNamespace := "nonamespacedefined"

			accRes, _ := idonia_auth.GetAccount(&authI)

			if accRes != nil && accRes.Account.Username != "" {
				defaultNamespace = accRes.Account.Username
			}

			req.Namespace = &defaultNamespace
			req.FileID, err = strconv.Atoi(fileID)
			if err != nil {
				util.JSONResponse(
					err,
					w,
					http.StatusInternalServerError,
				)
				return
			}
			res, err := idonia_share.PostPublicURL(&req, &authI)
			if err != nil {
				util.JSONResponse(
					err,
					w,
					http.StatusInternalServerError,
				)
				return
			}
			util.JSONResponse(
				res,
				w,
				http.StatusOK,
			)
			return
		}
	}
	if err != nil {
		util.JSONResponse(
			err,
			w,
			http.StatusInternalServerError,
		)
		return
	}
	util.JSONResponse(
		response.PublicURL,
		w,
		http.StatusOK,
	)
}
