package repository

import (
	"bitbucket.org/inehealth/idonia-pacs/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Task struct {
	UUID          string      `json:"uuid"`
	Step          uint32      `json:"step"`
	Steps         string      `json:"steps"`
	Status        string      `json:"status"`
	Error         model.Error `json:"error"`
	Authorization string      `json:"authorization"`
	StudyUID      string      `json:"studyUID"`
	IsReported    bool        `json:"isReported"`
	StudyID       string      `json:"studyID"`
	ContainerID   uint32      `json:"-"`
	InternalData  string      `json:"-"`
	CreatedAt     time.Time   `json:"createdAt"`
}

func InsertTask(db *sql.DB, task *Task) (err error) {
	_, err = db.Exec(
		`
			INSERT INTO task (
				uuid, 
				created_date, 
				auth,
				steps, 
				step, 
				status, 
				error,  
				study_uid, 
				container_id,
				study_id,
			    is_reported,              
				internal_data
			) VALUES (
				$1, 
				$2, 
				$3, 
				$4, 
				$5,
				$6, 
				$7, 
				$8, 
				$9, 
				$10, 
				$11, 
				$12
			)`,
		task.UUID,
		time.Now(),
		task.Authorization,
		task.Steps,
		task.Step,
		task.Status,
		task.Error.String(),
		task.StudyUID,
		task.ContainerID, task.StudyID, task.IsReported, task.InternalData,
	)
	return
}

func GetTask(db *sql.DB, uuid string) (task *Task, err error) {
	task = &Task{}
	var errString string
	res := db.QueryRow("SELECT uuid,auth,steps,step,status,error,study_id,is_reported,study_uid,container_id,internal_data,created_date FROM task WHERE uuid=$1", uuid)
	var date string
	err = res.Scan(
		&task.UUID,
		&task.Authorization,
		&task.Steps,
		&task.Step,
		&task.Status,
		&errString,
		&task.StudyID,
		&task.IsReported,
		&task.StudyUID,
		&task.ContainerID,
		&task.InternalData,
		&date,
	)
	if date != "" {
		task.CreatedAt, err = time.Parse("2006-01-02T15:04:05.999999-07:00", date)
	}

	json.Unmarshal([]byte(errString), &task.Error)
	return
}

func GetIncompletedTasksUUIDsByStudyUID(db *sql.DB, studyUID string, sopInstanceUID string) (tasks []string, err error) {
	rows, err := db.Query(`
		SELECT 
			uuid 
		FROM 
			task, task_dicom
		WHERE 
			task.uuid = task_dicom.task_uuid AND 
			task_dicom.study_uid = $1 AND task_dicom.sop_instance_uid != $2
		`,
		studyUID, sopInstanceUID)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var uuid string
		err = rows.Scan(&uuid)
		tasks = append(tasks, uuid)
	}
	return tasks, nil
}

func GetPendingTasks(db *sql.DB) (tasks []*Task, err error) {
	query := `
		SELECT 
			uuid,
		    auth, 
			steps,
		    step,
			status, 
			error,
			study_id,
			is_reported,
			study_uid, 
			container_id, 
			internal_data,
		    created_date   
		FROM task WHERE status != 'completed'  ORDER BY created_date DESC
	`
	rows, err := db.Query(query)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return
	}
	for rows.Next() {
		task := &Task{}
		var errString string
		var date string
		err = rows.Scan(
			&task.UUID,
			&task.Authorization,
			&task.Steps,
			&task.Step,
			&task.Status,
			&errString,
			&task.StudyID,
			&task.IsReported,
			&task.StudyUID,
			&task.ContainerID,
			&task.InternalData,
			&date,
		)
		if err != nil {
			return
		}
		task.CreatedAt, err = time.Parse("2006-01-02T15:04:05.999999-07:00", date)
		json.Unmarshal([]byte(errString), &task.Error)
		tasks = append(tasks, task)
	}
	return
}

func GetTasksByStudyUID(db *sql.DB, studyUID string) (tasks []*Task, err error) {
	query := `
		SELECT 
			uuid,
		    auth, 
			steps,
		    step,
			status, 
			error,
			study_id,
			is_reported,
			study_uid, 
			container_id, 
			internal_data 
		FROM task WHERE study_uid == $1 ORDER BY created_date DESC
	`
	rows, err := db.Query(query, studyUID)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return
	}
	for rows.Next() {
		task := &Task{}
		var errString string
		err = rows.Scan(
			&task.UUID,
			&task.Authorization,
			&task.Steps,
			&task.Step,
			&task.Status,
			&errString,
			&task.StudyID,
			&task.IsReported,
			&task.StudyUID,
			&task.ContainerID,
			&task.InternalData,
		)
		if err != nil {
			return
		}
		json.Unmarshal([]byte(errString), &task.Error)
		tasks = append(tasks, task)
	}
	return
}

func GetTasks(db *sql.DB, idoniaAccountID uint32, limit uint32, offset uint32) (tasks []*Task, total uint32, err error) {
	query := `
		SELECT 
			uuid,
		    auth, 
			steps,
		    step,
			status, 
			error,
			study_id,
			is_reported,
			study_uid, 
			container_id, 
			internal_data, 
			(SELECT count() FROM task) AS total
		FROM task
		ORDER BY created_date DESC
	`
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d ", limit)
		if offset > 0 {
			query += fmt.Sprintf(" OFFSET %d ", offset)
		}
	}

	rows, err := db.Query(query)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return
	}

	for rows.Next() {
		task := &Task{}
		var errString string
		err = rows.Scan(
			&task.UUID,
			&task.Authorization,
			&task.Steps,
			&task.Step,
			&task.Status,
			&errString,
			&task.StudyID,
			&task.IsReported,
			&task.StudyUID,
			&task.ContainerID,
			&task.InternalData,
			&total)
		if err != nil {
			return
		}
		json.Unmarshal([]byte(errString), &task.Error)
		tasks = append(tasks, task)
	}
	return
}

func EditTask(db *sql.DB, uuid string, task *Task) (err error) {
	_, err = db.Exec(
		`
			UPDATE task 
			SET 
				uuid = $1, 
				auth = $2, 
				steps = $3,  
				step = $4, 
				status = $5, 
				error = $6, 
			    study_id = $7,
			    is_reported =$8,
				study_uid = $9, 
				container_id = $10, 
				internal_data = $11 
			WHERE uuid = $12
			`,
		task.UUID,
		task.Authorization,
		task.Steps,
		task.Step,
		task.Status,
		task.Error.String(),
		task.StudyID,
		task.IsReported,
		task.StudyUID,
		task.ContainerID,
		task.InternalData,
		uuid,
	)
	return
}

func DeleteTask(db *sql.DB, uuid string) (err error) {
	_, err = db.Exec("DELETE FROM task WHERE uuid = $1", uuid)
	return
}
