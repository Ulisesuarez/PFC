package repository

import (
	"bitbucket.org/inehealth/idonia-pacs/model"
	"database/sql"
	"encoding/json"
)

type TaskDicom struct {
	UUID              string      `json:"uuid"`
	TaskUUID          string      `json:"-"`
	SeriesInstanceUID string      `json:"serieInstanceUID"`
	StudyUID          string      `json:"studyUID"`
	ContainerID       uint32      `json:"containerID"`
	SOPInstanceUID    string      `json:"sopInstanceUID"`
	Status            string      `json:"status"`
	Error             model.Error `json:"error"`
	IsPending         bool        `json:"isPending"`
}

func InsertTaskDicom(db *sql.DB, taskDICOM *TaskDicom) (err error) {
	_, err = db.Exec(
		"INSERT INTO task_dicom ("+
			"uuid, "+
			"task_uuid,"+
			"series_uid,"+
			"study_uid,"+
			"container_id,"+
			"sop_uid,"+
			"status,"+
			"error,"+
			"is_pending)"+
			" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		taskDICOM.UUID,
		taskDICOM.TaskUUID,
		taskDICOM.SeriesInstanceUID,
		taskDICOM.StudyUID,
		taskDICOM.ContainerID,
		taskDICOM.SOPInstanceUID,
		taskDICOM.Status,
		taskDICOM.Error.String(),
		taskDICOM.IsPending,
	)
	return
}

func GetTaskDicom(db *sql.DB, uuid string) (taskDICOM *TaskDicom, err error) {
	taskDICOM = &TaskDicom{}
	var errString string
	res := db.QueryRow(
		"SELECT "+
			"uuid, "+
			"task_uuid,"+
			"series_uid,"+
			"study_uid,"+
			"container_id,"+
			"sop_uid,"+
			"status,"+
			"error,"+
			"is_pending"+
			" FROM task_dicom WHERE uuid = $1",
		uuid,
	)
	err = res.Scan(
		&taskDICOM.UUID,
		&taskDICOM.TaskUUID,
		&taskDICOM.SeriesInstanceUID,
		&taskDICOM.StudyUID,
		&taskDICOM.ContainerID,
		&taskDICOM.SOPInstanceUID,
		&taskDICOM.Status,
		&errString,
		&taskDICOM.IsPending)
	json.Unmarshal([]byte(errString), &taskDICOM.Error)
	return
}

func GetTaskDicoms(db *sql.DB, taskUUID string) (tasksDICOM []*TaskDicom, err error) {
	rows, err := db.Query(
		"SELECT uuid,task_uuid,series_uid,study_uid,container_id,sop_uid,status,error,is_pending FROM task_dicom WHERE task_uuid = $1", taskUUID)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return
	}
	for rows.Next() {
		taskDICOM := &TaskDicom{}
		var errString string
		err = rows.Scan(
			&taskDICOM.UUID,
			&taskDICOM.TaskUUID,
			&taskDICOM.SeriesInstanceUID,
			&taskDICOM.StudyUID,
			&taskDICOM.ContainerID,
			&taskDICOM.SOPInstanceUID,
			&taskDICOM.Status,
			&errString,
			&taskDICOM.IsPending)
		if err != nil {
			return
		}
		json.Unmarshal([]byte(errString), &taskDICOM.Error)
		tasksDICOM = append(tasksDICOM, taskDICOM)
	}
	return
}

func GetTaskDicomWithStatus(db *sql.DB, taskUUID string, status string) (tasksDICOM []*TaskDicom, err error) {
	rows, err := db.Query(`
		SELECT
			uuid,  
			task_uuid, 
			series_uid, 
			study_uid, 
			container_id
			sop_uid, 
			status, 
			error,
		    is_pending   
		FROM task_dicom WHERE task_uuid = $1 AND status = $2`, taskUUID, status)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return
	}
	for rows.Next() {
		taskDICOM := &TaskDicom{}
		var errString string
		err = rows.Scan(
			&taskDICOM.UUID,
			&taskDICOM.TaskUUID,
			&taskDICOM.SeriesInstanceUID,
			&taskDICOM.StudyUID,
			&taskDICOM.ContainerID,
			&taskDICOM.SOPInstanceUID,
			&taskDICOM.Status,
			&errString,
			&taskDICOM.IsPending)
		if err != nil {
			return
		}
		json.Unmarshal([]byte(errString), &taskDICOM.Error)
		tasksDICOM = append(tasksDICOM, taskDICOM)
	}
	return
}

func GetTaskDicomByStudyID(db *sql.DB, studyUID string) (tasksDICOM []*TaskDicom, err error) {
	rows, err := db.Query(`
		SELECT
			uuid,  
			task_uuid, 
			serie_uid, 
			study_uid, 
			container_id
			sop_instance_uid, 
			status, 
			error 
		FROM task_dicom WHERE study_uid = $1 AND status = 'completed'`, studyUID)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return
	}
	for rows.Next() {
		taskDICOM := &TaskDicom{}
		var errString string
		err = rows.Scan(&taskDICOM.UUID, &taskDICOM.TaskUUID, &taskDICOM.SeriesInstanceUID, &taskDICOM.StudyUID, &taskDICOM.ContainerID, &taskDICOM.SOPInstanceUID, &taskDICOM.Status, &errString)
		if err != nil {
			return
		}
		json.Unmarshal([]byte(errString), &taskDICOM.Error)
		tasksDICOM = append(tasksDICOM, taskDICOM)
	}
	return
}
func EditTaskDicom(db *sql.DB, uuid string, taskDICOM *TaskDicom) (err error) {
	_, err = db.Exec(
		`
		UPDATE task_dicom 
		SET uuid = $1,
		    task_uuid = $2, 
		    series_uid = $3,
		    study_uid= $4,
		    container_id = $5,
		    sop_uid = $6,
		    status = $7,
		    error = $8, 
		    is_pending = $9 WHERE uuid = $10`,
		taskDICOM.UUID,
		taskDICOM.TaskUUID,
		taskDICOM.SeriesInstanceUID,
		taskDICOM.StudyUID,
		taskDICOM.ContainerID,
		taskDICOM.SOPInstanceUID,
		taskDICOM.Status,
		taskDICOM.Error.String(),
		taskDICOM.IsPending,
		taskDICOM.UUID)
	return
}

func DeleteTaskDicom(db *sql.DB, uuid string) (err error) {
	_, err = db.Exec("DELETE FROM task_dicom WHERE uuid = $1", uuid)
	return
}
