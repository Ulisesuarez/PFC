package task

import (
	"bitbucket.org/inehealth/idonia-pacs/configuration"
	"bitbucket.org/inehealth/idonia-pacs/model"
	"bitbucket.org/inehealth/idonia-pacs/repository"
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia-core"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gopkg.in/oleiade/lane.v1"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Engine struct {
	db           *sql.DB
	queue        *lane.Queue
	stopThreads  map[int]chan bool
	stop         chan bool
	logger       *logrus.Logger
	config       *configuration.Configuration
	studies      map[string]string
	CreateTask   *sync.Mutex
	TasksWaiting map[string]chan interface{}
}

func NewEngine(db *sql.DB, logger *logrus.Logger, config *configuration.Configuration) (engine *Engine) {
	studies := make(map[string]string)

	return &Engine{
		db:           db,
		queue:        lane.NewQueue(),
		stopThreads:  make(map[int]chan bool),
		stop:         make(chan bool),
		logger:       logger,
		config:       config,
		studies:      studies,
		CreateTask:   &sync.Mutex{},
		TasksWaiting: make(map[string]chan interface{}),
	}
}

func (engine *Engine) Add(task *Task) (err error) {
	err = insertTaskToDB(engine.db, engine.logger, task)
	if err != nil {
		engine.logger.Debugf("[ engine ] Add: engine.insertTaskToDB(task) err %v", err)
		return
	}
	engine.queue.Enqueue(task)
	if engine.logger != nil {
		engine.logger.WithFields(logrus.Fields{
			"UUID":         task.UUID,
			"IdoniaAPIKey": task.Authorization.APIKey,
		}).Infof("[ engine ] New task queued")
	}
	return
}

func (engine *Engine) Start(threads int) {
	pendingTasks, err := readPendingTasksFromDB(engine.db, engine.logger)
	if err != nil {
		engine.logger.Debugf("[ engine ] Start: engine.readPendingTasksFromDB() err %v", err)
		panic(err)
	}

	for _, v := range pendingTasks {
		engine.queue.Enqueue(v)
		if engine.logger != nil {
			engine.logger.WithFields(logrus.Fields{
				"UUID":         v.UUID,
				"IdoniaAPIKey": v.Authorization.APIKey,
			}).Infof("[ engine ] New task queued")
		}
	}

	for thread := 0; thread < threads; thread++ {
		go func(thread int) {
			for {
				select {
				case <-engine.stopThreads[thread]:
					if engine.logger != nil {
						if engine.logger != nil {
							engine.logger.WithFields(logrus.Fields{
								"Thread": thread,
							}).Debugf("[ engine ] Stopping thread")
						}
					}
					break
				default:
					if taskElement := engine.queue.Dequeue(); taskElement != nil {
						task := taskElement.(*Task)
						if err = engine.executeTask(task); err != nil {
							engine.logger.WithError(err).WithFields(logrus.Fields{
								"Step": task.Step,
							}).Error("[ engine ] Error executing task!")
						}
					} else {
						time.Sleep(5 * time.Second)
					}
				}
			}
		}(thread)
	}

	<-engine.stop
	for thread := 0; thread < threads; thread++ {
		engine.stopThreads[thread] <- true
	}
}

func (engine *Engine) Stop() {
	if engine.logger != nil {
		engine.logger.Debugf("[ engine ] Stopping task engine")
	}
	engine.stop <- true
}

func insertTaskToDB(db *sql.DB, logger *logrus.Logger, task *Task) (err error) {
	taskUUID := task.UUID
	if len(taskUUID) == 0 {
		taskUUID = uuid.New().String()
	}
	for _, v := range task.DICOMFiles {

		err = repository.InsertTaskDicom(db, &repository.TaskDicom{
			UUID:              v.UUID,
			TaskUUID:          taskUUID,
			SeriesInstanceUID: v.SeriesInstanceUID,
			StudyUID:          v.StudyUID,
			ContainerID:       0,
			IsPending:         v.IsPending,
			SOPInstanceUID:    v.SOPInstanceUID,
			Status:            string(New),
		})
		if err != nil {
			logger.Debugf("[ engine ] insertTaskToDB: repository.InsertTaskDicom(engine.db, &repository.TaskDicom...) err %v", err)
			return err
		}
	}
	for k, v := range task.AdditionalFields {
		err = repository.InsertTaskAdditionalField(db, &repository.TaskAdditionalField{
			TaskUUID: taskUUID,
			Key:      k,
			Value:    v,
		})
		if err != nil {
			logger.Debugf("[ engine ] insertTaskToDB: repository.InsertTaskAdditionalField(engine.db, &repository.TaskAdditionalField...) err %v", err)
			return
		}
	}
	authorization, err := json.Marshal(task.Authorization)
	if err != nil {
		logger.Debugf("[ engine ] insertTaskToDB: json.Marshal(task.Destination) err %v", err)
		return
	}
	steps, err := json.Marshal(task.Steps)
	if err != nil {
		logger.Debugf("[ engine ] insertTaskToDB: json.Marshal(task.Destination) err %v", err)
		return
	}
	err = repository.InsertTask(db, &repository.Task{
		UUID:          taskUUID,
		Authorization: string(authorization),
		Steps:         string(steps),
		Status:        string(New),
		Error:         model.Error{},
		StudyID:       task.StudyID,
		StudyUID:      task.StudyUID,
		ContainerID:   task.ContainerID,
	})
	if err != nil {
		logger.Debugf("[ engine ] insertTaskToDB: repository.InsertTask(engine.db, &repository.Task... err %v", err)
	}
	return

}

func updateTaskInDB(db *sql.DB, logger *logrus.Logger, task *Task) (err error) {
	steps, err := json.Marshal(task.Steps)
	if err != nil {
		logger.Debugf("[ engine ] updateTaskInDB: json.Marshal(task.Destination) err %v", err)
		return
	}
	auth, err := json.Marshal(task.Authorization)
	if err != nil {
		logger.Debugf("[ engine ] updateTaskInDB: json.Marshal(v.Description) err %v", err)
		return err
	}
	err = repository.EditTask(db, task.UUID, &repository.Task{
		UUID:          task.UUID,
		Authorization: string(auth),
		Steps:         string(steps),
		StudyID:       task.StudyID,
		IsReported:    task.IsReported,
		Step:          task.Step,
		Status:        string(task.Status),
		Error:         task.Error,
		StudyUID:      task.StudyUID,
		ContainerID:   task.ContainerID,
	})
	if err != nil {
		logger.Debugf("[ engine ] updateTaskInDB: repository.EditTask(engine.db, task.UUID, &repository.Task... err %v", err)
		return
	}

	for _, v := range task.DICOMFiles { // UPSERT
		_, err := repository.GetTaskDicom(db, v.UUID)
		if err != nil {
			err = nil
			err = repository.InsertTaskDicom(db, &repository.TaskDicom{
				UUID:              v.UUID,
				TaskUUID:          task.UUID,
				SeriesInstanceUID: v.SeriesInstanceUID,
				StudyUID:          v.StudyUID,
				SOPInstanceUID:    v.SOPInstanceUID,
				Status:            string(v.Status),
				Error:             v.Error,
			})
		} else {
			err = repository.EditTaskDicom(db, v.UUID, &repository.TaskDicom{
				UUID:              v.UUID,
				TaskUUID:          task.UUID,
				SeriesInstanceUID: v.SeriesInstanceUID,
				StudyUID:          v.StudyUID,
				SOPInstanceUID:    v.SOPInstanceUID,
				Status:            string(v.Status),
				Error:             v.Error,
			})
		}
		if err != nil {
			logger.Debugf("[ engine ] updateTaskInDB: repository.EditTaskDicom(engine.db, v.UUID, &repository.TaskDicom... err %v", err)
			return err
		}
	}

	err = repository.DeleteTaskAdditionalFields(db, task.UUID)
	if err != nil {
		logger.Debugf("[ engine ] updateTaskInDB: repository.DeleteTaskAdditionalFields(engine.db, task.UUID) err %v", err)
		return
	}
	for k, v := range task.AdditionalFields {
		err = repository.InsertTaskAdditionalField(db, &repository.TaskAdditionalField{
			TaskUUID: task.UUID,
			Key:      k,
			Value:    v,
		})
		if err != nil {
			logger.Debugf("[ engine ] updateTaskInDB: InsertTaskAdditionalField(engine.db, &repository.TaskAdditionalField... err %v", err)
			return
		}
	}
	return
}

func readTaskFromDB(db *sql.DB, logger *logrus.Logger, uuid string) (task *Task, err error) {
	fmt.Println(uuid)
	simpleTask, err := repository.GetTask(db, uuid)
	if err != nil {
		fmt.Println(1)
		logger.Debugf("[ engine ] readTaskFromDB: repository.GetTask(engine.db, uuid) err %v", err)
		return
	}
	task = &Task{
		UUID:        simpleTask.UUID,
		Status:      Status(simpleTask.Status),
		Step:        simpleTask.Step,
		Error:       simpleTask.Error,
		StudyUID:    simpleTask.StudyUID,
		StudyID:     simpleTask.StudyID,
		IsReported:  simpleTask.IsReported,
		ContainerID: simpleTask.ContainerID,
	}
	err = json.Unmarshal([]byte(simpleTask.Steps), &task.Steps)
	if err != nil {
		fmt.Println(2)
		logger.Debugf("[ engine ] readTaskFromDB: json.Unmarshal([]byte(simpleTask.Destination), &task.Destination) err %v", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(simpleTask.Authorization), &task.Authorization)
	if err != nil {

		fmt.Println(3)
		logger.Debugf("[ engine ] readTaskFromDB: json.Unmarshal([]byte(v.Description), &description) err %v", err)
		return nil, err
	}
	task.AdditionalFields = make(map[string]string)
	DICOMFiles, err := repository.GetTaskDicoms(db, simpleTask.UUID)
	if err != nil {
		fmt.Println(4)
		logger.Debugf("[ engine ] readTaskFromDB: repository.GetTaskDicoms(engine.db, simpleTask.UUID) err %v", err)
		return nil, err
	}
	for _, v := range DICOMFiles {

		task.DICOMFiles = append(task.DICOMFiles, DICOMFile{
			UUID:              v.UUID,
			SeriesInstanceUID: v.SeriesInstanceUID,
			StudyUID:          v.StudyUID,
			ContainerID:       v.ContainerID,
			SOPInstanceUID:    v.SOPInstanceUID,
			Status:            Status(v.Status),
			Error:             v.Error,
		})
	}

	AdditionalFields, err := repository.GetTaskAdditionalFields(db, simpleTask.UUID)
	if err != nil {
		fmt.Println(5)
		logger.Debugf("[ engine ] readTaskFromDB: repository.GetTaskAdditionalFields(engine.db, simpleTask.UUID) err %v", err)
		return nil, err
	}
	for _, v := range AdditionalFields {
		task.AdditionalFields[v.Key] = v.Value
	}

	return
}

func readPendingTasksFromDB(db *sql.DB, logger *logrus.Logger) (tasks []*Task, err error) {
	repoTasks, err := repository.GetPendingTasks(db)
	if err != nil {
		fmt.Println("pending")
		logger.Debugf("[ engine ] readPendingTasksFromDB: repository.GetPendingTasks(engine.db) err %v", err)
		return nil, err
	}
	for _, v := range repoTasks {
		task, err := readTaskFromDB(db, logger, v.UUID)
		if err != nil {
			fmt.Println("read task")
			logger.Debugf("[ engine ] readPendingTasksFromDB: engine.readTaskFromDB(v.UUID) err %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return
}

func readAffectedTasksByUpdateFromDB(db *sql.DB, logger *logrus.Logger, studyUID string) (tasks []*Task, err error) {
	repoTasks, err := repository.GetTasksByStudyUID(db, studyUID)
	if err != nil {
		logger.Debugf("[ engine ] readAffectedTasksByUpdateFromDB: repository.GetTasksByStudyUID(engine.db) err %v", err)
		return nil, err
	}
	for _, v := range repoTasks {
		task, err := readTaskFromDB(db, logger, v.UUID)
		if err != nil {
			logger.Debugf("[ engine ] readAffectedTasksByUpdateFromDB: engine.readTaskFromDB(v.UUID) err %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return
}

func (engine *Engine) updateDICOMFiles(task *Task) (err error) {
	//TODO RESYNC FOR NEW changes

	return

}

func (engine *Engine) retrieveDICOM(task *Task) (err error) {
	if len(task.StudyUID) > 0 {
		var InstancesInPacs []DICOMInstance

		url := GetStudyInstancesURL(task.StudyUID, *engine.config)
		response, err := http.Get(url)
		if response == nil {
			if err != nil {
				fmt.Println("Error 8")
				return err
			}
			fmt.Println("Error 9")
			return err
		}
		fmt.Println(response.Body)
		fmt.Println(url)
		body, err := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(body, &InstancesInPacs)
		if err != nil {
			fmt.Println("Error 10")
			return err
		}
		var DICOMsInPacs []DICOMFile
		DICOMsInPacs = GetDicomFilesFromInstance(InstancesInPacs)
		study, err := repository.GetStudy(engine.db, task.StudyUID)
		if err != nil {
			if err != sql.ErrNoRows {
				fmt.Println("Error 1")
				return err
			}
			task.AddUploadAllToIdoniaSteps()
			task.DICOMFiles = DICOMsInPacs
			err = updateTaskInDB(engine.db, engine.logger, task)
			fmt.Println("Error 2")
			return err
		}
		container, err := idonia_core.GetContainer(&idonia_core.GetContainerReq{
			ContainerID: &study.ContainerID}, task.Authorization)
		if err != nil {
			// suponer que se ha eliminado el estudio
			engine.logger.Error(err.Error())
			task.AddUploadAllToIdoniaSteps()
			task.DICOMFiles = DICOMsInPacs
			err = repository.DeleteStudy(engine.db, task.StudyUID)
			if err != nil {
				fmt.Println("Error 3")
				return err
			}
			err = updateTaskInDB(engine.db, engine.logger, task)
			if err != nil {
				fmt.Println("Error 4")
				return err
			}
			fmt.Println("Error 5")
			return err
		}
		var studyDICOMs []DICOMFile
		err = json.Unmarshal([]byte(study.DICOMFiles), &studyDICOMs)
		if err != nil {
			if err != nil {
				fmt.Println("Error 6")
				return err
			}
			if err != nil {
				fmt.Println("Error 7")
				return err
			}
			return err
		}

		DICOMsInIdonia := GetNumberOfDicomInstances(container)
		if DICOMsInIdonia != len(studyDICOMs) {
			//TODO
		}

		task.StudyID = study.StudyID

		task.ContainerID = study.ContainerID
		task.IsReported = study.IsReported
		if len(DICOMsInPacs) != len(studyDICOMs) {

			task.DICOMFiles = DICOMsInPacs
			task.markPendingDICOM(studyDICOMs)
			task.AddUploadPendingStep()
		}

		//call get container if numImages are == to completeFiles
		// SE puede repetir upload task.Steps.type == reupload entonces ignoramos el ignorar
		// sino saltamos todo esto

		//engine.updateDICOMFiles(task)
	}

	err = updateTaskInDB(engine.db, engine.logger, task)
	fmt.Println("Error 11")
	return
}

func GetNumberOfDicomInstances(container *idonia_core.GetContainerRes) (numberOfDICOM int) {
	for _, patient := range container.DicomPatients {
		for _, study := range patient.DicomStudies {
			for _, serie := range study.DicomSeries {
				numberOfDICOM = numberOfDICOM + serie.InstancesCount
			}
		}
	}
	return
}

func (engine *Engine) executeTask(task *Task) (err error) {
	for {
		if engine.logger != nil {
			engine.logger.WithFields(logrus.Fields{
				"UUID":         task.UUID,
				"IdoniaAPIKey": task.Authorization.APIKey,
				"Status":       task.Status,
			}).Infof("[ engine ] Task Processing")
		}
		if err != nil {
			engine.logger.Debugf("[ engine ] executeTask: err %v", err)
			break
		}
		switch task.Status {
		default:
			fallthrough
		case New:
			task.Status = Retrieving
			err = updateTaskInDB(engine.db, engine.logger, task)
		case Retrieving:
			err = engine.retrieveDICOM(task)
			if err != nil {
				break
			}
			task.Status = Retrieved
			err = updateTaskInDB(engine.db, engine.logger, task)
		case Retrieved:
			task.Status = Sending
			err = updateTaskInDB(engine.db, engine.logger, task)
		case Sending:
			var noMoreSteps bool
			if noMoreSteps, err = task.RunNextStep(*engine.config); err != nil || noMoreSteps {
				if err != nil {
					break
				}
				if noMoreSteps {
					task.Status = Sent
				}
			}
			err = updateTaskInDB(engine.db, engine.logger, task)
		case Sent:
			task.Status = Completed
			err = updateTaskInDB(engine.db, engine.logger, task)
		case Error:
			if err == nil {
				return task.Error
			}
			return
		case Completed:
			err = saveStudyInDB(engine.db, engine.logger, task)
			if err != nil {
				fmt.Println(err.Error())
				task.Error = model.Error{Code: 500, Message: err.Error(), AdditionalMessage: "Cant save study as completed"}
				return task.Error
			}
			err = updateTaskInDB(engine.db, engine.logger, task)
			return
		}
	}
	if err != nil {
		task.Status = Error
		task.Error = model.ErrConf
		task.Error.AdditionalMessage = err.Error()
		erro := updateTaskInDB(engine.db, engine.logger, task)
		if erro != nil {
			engine.logger.Debugf("[ engine ] executeTask: Error during updateTask (status=Error): err %v", erro)
		}
		engine.logger.WithFields(logrus.Fields{
			"UUID":         task.UUID,
			"IdoniaAPIKey": task.Authorization.APIKey,
			"Status":       task.Status,
		}).WithError(err).Error("[ engine ] Task Error")
	}

	return
}

func saveStudyInDB(db *sql.DB, logger *logrus.Logger, task *Task) error {
	if task.ContainerID == 0 {
		fmt.Println(fmt.Sprintf("%+v", task))
		if cID, ok := task.AdditionalFields["ContainerID"]; ok {
			containerId, err := strconv.Atoi(cID)
			if err != nil {
				return errors.New("ContainerID cant be empty")
			}
			task.ContainerID = uint32(containerId)
			if task.ContainerID == 0 {
				return errors.New("ContainerID cant be empty")
			}
		}
	}
	if task.StudyID == "" {
		fmt.Println(fmt.Sprintf("%+v", task))
		if studyID, ok := task.AdditionalFields["StudyID"]; ok {
			task.StudyID = studyID
			if task.StudyID == "" {
				return errors.New("StudyID cant be empty")
			}
		}
		return errors.New("StudyID cant be empty")
	}
	study, err := repository.GetStudy(db, task.StudyUID)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		logger.Info("Study no exists in DB:", task.StudyUID)
	}
	dicom, err := json.Marshal(task.DICOMFiles)
	if err != nil {
		return err
	}
	aStudy := &repository.Study{
		UUID:        uuid.New().String(),
		StudyUID:    task.StudyUID,
		IsReported:  task.IsReported,
		ContainerID: task.ContainerID,
		StudyID:     task.StudyID,
		DICOMFiles:  string(dicom),
	}
	if url,ok:=task.AdditionalFields["PublicURL"]; ok{
		aStudy.MagicLink = url
	}

	fmt.Println(study)
	if study == nil {
		err = repository.InsertStudy(db, *aStudy)
		if err != nil {
			logger.Error("Error inserting study:", err.Error())
		}
	} else {
		err = repository.UpdateStudy(db, *aStudy)
		if err != nil {
			logger.Error("Error updating study:", err.Error())
		}
	}

	return err
}

/*func RunContainerUpdate(db *sql.DB, logger *logrus.Logger, task *Task, datasets map[string]*dicom.DataSet) (updated bool, err error) {
	var newDICOMFiles []DICOMFile

	for path, dataset := range datasets {
		found := false
		for _, dicom := range task.DICOMFiles {
			if dicom.Path == path {
				found = true
			}
		}
		if !found {
			namestring := ""
			patientNamestring := ""
			patientIdstring := ""
			modalitystring := ""
			modalitiesInStudystring := ""
			institutionNamestring := ""
			studyInstanceUIDstring := ""
			studyIDstring := ""
			studyDatestring := ""
			studyTimestring := ""
			studyDescriptionstring := ""
			bodyPartExaminedstring := ""
			name, err := dataset.FindElementByTag(dicomtag.SOPInstanceUID)
			if err == nil {
				namestring, _ = name.GetString()
			}
			patientName, err := dataset.FindElementByTag(dicomtag.PatientName)
			if err == nil {
				patientNamestring, _ = patientName.GetString()
			}
			patientId, err := dataset.FindElementByTag(dicomtag.PatientID)
			if err == nil {
				patientIdstring, _ = patientId.GetString()
			}
			modality, err := dataset.FindElementByTag(dicomtag.Modality)
			if err == nil {
				modalitystring, _ = modality.GetString()
			}
			modalitiesInStudy, err := dataset.FindElementByTag(dicomtag.ModalitiesInStudy)
			if err == nil {
				modalitiesInStudystring, _ = modalitiesInStudy.GetString()
			}
			institutionName, err := dataset.FindElementByTag(dicomtag.InstitutionName)
			if err == nil {
				institutionNamestring, _ = institutionName.GetString()
			}
			studyInstanceUID, err := dataset.FindElementByTag(dicomtag.StudyInstanceUID)
			if err == nil {
				studyInstanceUIDstring, _ = studyInstanceUID.GetString()
			}
			studyID, err := dataset.FindElementByTag(dicomtag.StudyID)
			if err == nil {
				studyIDstring, _ = studyID.GetString()
			}
			studyDate, err := dataset.FindElementByTag(dicomtag.StudyDate)
			if err == nil {
				studyDatestring, _ = studyDate.GetString()
			}
			studyTime, err := dataset.FindElementByTag(dicomtag.StudyTime)
			if err == nil {
				studyTimestring, _ = studyTime.GetString()
			}
			studyDescription, err := dataset.FindElementByTag(dicomtag.StudyDescription)
			if err == nil {
				studyDescriptionstring, _ = studyDescription.GetString()
			}
			bodyPartExamined, err := dataset.FindElementByTag(dicomtag.BodyPartExamined)
			if err == nil {
				bodyPartExaminedstring, _ = bodyPartExamined.GetString()
			}

			newDICOMFiles = append(newDICOMFiles, DICOMFile{
				UUID: uuid.New().String(),
				Name: namestring + ".dcm",
				Description: DICOMDescription{
					PatientName:       patientNamestring,
					PatientId:         patientIdstring,
					Modality:          modalitystring,
					ModalitiesInStudy: modalitiesInStudystring,
					InstitutionName:   institutionNamestring,
					StudyInstanceUID:  studyInstanceUIDstring,
					StudyID:           studyIDstring,
					StudyDate:         studyDatestring,
					StudyTime:         studyTimestring,
					StudyDescription:  studyDescriptionstring,
					BodyPartExamined:  bodyPartExaminedstring,
				},
				Path:           path,
				StudyUID:       studyInstanceUIDstring,
				SOPInstanceUID: namestring,
				Status:         New,
			})
		}
	}

	task.DICOMFiles = append(task.DICOMFiles, newDICOMFiles...)

	for _, dicom := range newDICOMFiles {
		policy, err := idonia_core.GetPolicy(&idonia.Auth{
			AccountID: task.IdoniaAccountID,
			Token:     task.IdoniaToken,
		})
		if err != nil {
			task.Error = model.ErrIdonia
			task.Error.AdditionalMessage = err.Error()
			task.Status = Error
			return false, err
		}
		err = gcs.UploadFile(dicom.Path, policy)
		if err != nil {
			task.Error = model.ErrIdonia
			task.Error.AdditionalMessage = err.Error()
			task.Status = Error
			return false, err
		}

		_, err = idonia_core.PostContainers_ContainerIDDICOM(&idonia_core.PostContainers_ContainerIDDICOMReq{
			UUID: &policy.UUID,
		}, int(task.ContainerID), &idonia.Auth{
			AccountID: task.IdoniaAccountID,
			Token:     task.IdoniaToken,
		})
		if err != nil {
			task.Error = model.ErrIdonia
			task.Error.AdditionalMessage = err.Error()
			task.Status = Error
			return false, err
		}
	}

	updateTaskInDB(db, logger, task)
	return len(newDICOMFiles) > 0, nil
}
*/
/*func (engine *Engine) UpdatedStudy(studyUID string) (affected int) {
	tasks, _ := readAffectedTasksByUpdateFromDB(engine.db, engine.logger, studyUID)
	affectedTasks := 0
	wg := sync.WaitGroup{}
	for _, task := range tasks {
		wg.Add(1)
		go func(task *Task) {
			defer wg.Done()
			for {
				task, _ := readTaskFromDB(engine.db, engine.logger, task.UUID)
				switch task.Status {
				case Error:
					return
				case Completed:
					paths, _ := engine.dicomManager.GetStudyFiles(studyUID)
					files := make(map[string]*dicom.DataSet)
					for _, path := range paths {
						files[path], _ = engine.dicomManager.GetDCMDataset(path)
					}
					updated, err := RunContainerUpdate(engine.db, engine.logger, task, files)
					if err == nil && updated {
						affectedTasks++
					}
					return
				case Retrieving:
					if engine.TasksWaiting[task.UUID] != nil {
						close(engine.TasksWaiting[task.UUID])
					}
					affectedTasks++
					return
				default:
					//Wait for the task to finish
					time.Sleep(10 * time.Second)
				}
			}
		}(task)
	}
	wg.Wait()
	return affectedTasks
}*/

func (engine *Engine) ErrorWhileRetrieving(studyUID string, err error) {
	tasks, _ := readAffectedTasksByUpdateFromDB(engine.db, engine.logger, studyUID)
	for _, task := range tasks {
		if task.Status == Retrieving && engine.TasksWaiting[task.UUID] != nil {
			engine.TasksWaiting[task.UUID] <- err
		}
	}
	return
}

func (engine *Engine) GetCompletedFiles(taskUUID string) (dicomFiles []DICOMFile, err error) {
	dicomTask, err := repository.GetTaskDicomWithStatus(engine.db, taskUUID, Completed)
	for _, dicom := range dicomTask {
		dicomFile := DICOMFile{
			UUID:              dicom.UUID,
			StudyUID:          dicom.StudyUID,
			SeriesInstanceUID: dicom.SeriesInstanceUID,
			ContainerID:       dicom.ContainerID,
			SOPInstanceUID:    dicom.SOPInstanceUID,
			Status:            Status(dicom.Status),
			Error:             dicom.Error,
		}
		dicomFiles = append(dicomFiles, dicomFile)
	}
	return
}
