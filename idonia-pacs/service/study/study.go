package study

import (
	"bitbucket.org/inehealth/idonia-pacs/configuration"
	"bitbucket.org/inehealth/idonia-pacs/repository"
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia"
	idonia_core "bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia-core"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
)

func RetriveAllStudies(db *sql.DB, logger *logrus.Logger, config *configuration.Configuration, auth *idonia.Auth) (studies []*repository.Study, err error) {

	studies, err = repository.GetAllStudies(db)
	if err != nil {
		return
	}

	for _, study := range studies {
		studyres, err := idonia_core.GetStudy(&idonia_core.GetStudyReq{
			FileID: study.StudyID,
		}, auth)
		if err != nil {
			logger.Errorln(err.Error())
			//err = repository.DeleteStudy(db, study.StudyUID)
		}
		if studyres == nil || studyres.IsDeleted {
			logger.Info(fmt.Sprintf("%+v", studyres))
			//err = repository.DeleteStudy(db, study.StudyUID)
			if err != nil {
				fmt.Println(err.Error())
			}
			continue
		}
		logger.Info(fmt.Sprintf("%+v", studyres))
		if studyres.LastReportID != nil {
			document, err := idonia_core.GetDocument(&idonia_core.GetDocumentReq{
				FileID: strconv.Itoa(int(*studyres.LastReportID)),
			}, auth)
			if err == nil && document.IsReport && !document.IsDeleted {
				study.ReportID = uint32(*studyres.LastReportID)
				study.IsReported = true
			} else {
				if err != nil {
					logger.Error(err.Error())
				}
				study.ReportID = 0
				study.IsReported = false
			}

		} else {
			study.ReportID = 0
			study.IsReported = false
		}
		logger.Info(fmt.Sprintf("%+v", study))
		container, err := idonia_core.GetContainer(&idonia_core.GetContainerReq{
			ContainerID: &study.ContainerID,
		}, auth)
		if err != nil {
			logger.Errorln(err.Error())
			//err = repository.DeleteStudy(db, study.StudyUID)
		}
		if container == nil {
			logger.Info(fmt.Sprintf("%+v", container))
			//err = repository.DeleteStudy(db, study.StudyUID)
			if err != nil {
				fmt.Println(err.Error())
			}
			continue
		}
		logger.Info(fmt.Sprintf("%+v", container))
		err = repository.UpdateStudy(db, *study)
		if err != nil {
			logger.Errorln(err.Error())
		}
	}
	studies, err = repository.GetAllStudies(db)
	return
}
