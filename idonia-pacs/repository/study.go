package repository

import (
	"database/sql"
)

type Study struct {
	UUID        string `json:"uuid"`
	StudyUID    string `json:"studyUID"`
	IsReported  bool   `json:"isReported"`
	ReportID    uint32 `json:"reportID"`
	StudyID     string `json:"studyID"`
	ContainerID uint32 `json:"containerID"`
	DICOMFiles  string `json:"dicomFiles"`
	MagicLink  string  `json:"magicLink"`
}

func GetStudy(db *sql.DB, studyUID string) (study *Study, err error) {
	study = &Study{}
	res := db.QueryRow("SELECT uuid, study_uid, is_reported, report_id, study_id,container_id, dicom_files, magic_link FROM study WHERE study_uid = $1", studyUID)
	err = res.Scan(
		&study.UUID,
		&study.StudyUID,
		&study.IsReported,
		&study.ReportID,
		&study.StudyID,
		&study.ContainerID,
		&study.DICOMFiles,
		&study.MagicLink,
	)
	if err != nil {
		return nil, err
	}

	return
}
func InsertStudy(db *sql.DB, study Study) (err error) {

	_, err = db.Exec("INSERT INTO study (uuid,is_reported,report_id,study_id,container_id,dicom_files,study_uid, magic_link) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
		study.UUID,
		study.IsReported,
		study.ReportID,
		study.StudyID,
		study.ContainerID,
		study.DICOMFiles,
		study.StudyUID,
		study.MagicLink)
	return
}
func UpdateStudy(db *sql.DB, study Study) (err error) {

	_, err = db.Exec("UPDATE study SET uuid = $1,is_reported=$2,report_id =$3,study_id=$4,container_id=$5,dicom_files=$6,magic_link=$7 WHERE study_uid=$8",
		study.UUID,
		study.IsReported,
		study.ReportID,
		study.StudyID,
		study.ContainerID,
		study.DICOMFiles,
		study.MagicLink,
		study.StudyUID,
		)
	return
}

func DeleteStudy(db *sql.DB, studyUID string) (err error) {
	_, err = db.Exec("DELETE FROM study WHERE study_uid = $1", studyUID)
	return
}

func GetAllStudies(db *sql.DB) (studies []*Study, err error) {
	studies = []*Study{}
	rows, err := db.Query("SELECT uuid, study_uid, is_reported, report_id, study_id, container_id, dicom_files, magic_link FROM study")
	if err != nil {
		return nil, err
	}
	if rows != nil {
		defer rows.Close()
	}
	for rows.Next() {
		var study = Study{}
		err = rows.Scan(
			&study.UUID,
			&study.StudyUID,
			&study.IsReported,
			&study.ReportID,
			&study.StudyID,
			&study.ContainerID,
			&study.DICOMFiles,
			&study.MagicLink,
		)
		if err != nil {
			return
		}
		studies = append(studies, &study)
	}

	return
}
