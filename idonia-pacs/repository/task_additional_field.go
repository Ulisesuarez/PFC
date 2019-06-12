package repository

import "database/sql"

type TaskAdditionalField struct {
	TaskUUID string `json:"-"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

func InsertTaskAdditionalField(db *sql.DB, taskAdditionalField *TaskAdditionalField) (err error) {
	var taskUUID, key string
	res := db.QueryRow(
		"SELECT task_uuid, key FROM task_additional_field WHERE task_uuid = $1 AND key = $2",
		taskAdditionalField.TaskUUID, taskAdditionalField.Key,
	)
	if err = res.Scan(&taskUUID, &key); err != nil {
		if err == sql.ErrNoRows {
			_, err = db.Exec(
				"INSERT INTO task_additional_field (task_uuid, key, value) VALUES ($1, $2, $3)",
				taskAdditionalField.TaskUUID, taskAdditionalField.Key, taskAdditionalField.Value,
			)
			return
		}
		return
	}
	_, err = db.Exec(
		"UPDATE task_additional_field SET value = $3  WHERE task_uuid = $1 AND key = $2",
		taskAdditionalField.TaskUUID, taskAdditionalField.Key, taskAdditionalField.Value,
	)
	return
}

func GetTaskAdditionalFields(db *sql.DB, taskUUID string) (taskAdditionalFields []*TaskAdditionalField, err error) {
	rows, err := db.Query("SELECT task_uuid, key, value FROM task_additional_field WHERE task_uuid = $1", taskUUID)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return
	}
	for rows.Next() {
		taskAdditionalField := &TaskAdditionalField{}
		err = rows.Scan(&taskAdditionalField.TaskUUID, &taskAdditionalField.Key, &taskAdditionalField.Value)
		if err != nil {
			return
		}
		taskAdditionalFields = append(taskAdditionalFields, taskAdditionalField)
	}
	return
}

func DeleteTaskAdditionalField(db *sql.DB, taskUUID string, key string) (err error) {
	_, err = db.Exec("DELETE FROM task_additional_field WHERE task_uuid = $1 AND key = $2", taskUUID, key)
	return
}

func DeleteTaskAdditionalFields(db *sql.DB, taskUUID string) (err error) {
	_, err = db.Exec("DELETE FROM task_additional_field WHERE task_uuid = $1", taskUUID)
	return
}
