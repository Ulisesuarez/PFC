package task

import (
	"bitbucket.org/inehealth/idonia-pacs/configuration"
	"net/url"
)

type DICOMInstance struct {
	StudyUID struct {
		Value []string `json:"Value"`
	} `json:"0020000D"`
	SerieInstanceUID struct {
		Value []string `json:"Value"`
	} `json:"0020000E"`
	SOPInstanceUID struct {
		Value []string `json:"Value"`
	} `json:"00080018"`
}

func GetDicomFilesFromInstance(instances []DICOMInstance) (dicomFiles []DICOMFile) {
	for _, instance := range instances {
		dicom := DICOMFile{}
		dicom.StudyUID = instance.StudyUID.Value[0]
		dicom.SeriesInstanceUID = instance.SerieInstanceUID.Value[0]
		dicom.SOPInstanceUID = instance.SOPInstanceUID.Value[0]
		dicomFiles = append(dicomFiles, dicom)
	}
	return
}

func GetDICOMURL(dicomFile DICOMFile, conf configuration.Configuration) string {
	dicomURL := url.URL{}
	dicomURL.Scheme = "http"
	dicomURL.Host = conf.PACS.Host
	dicomURL.Path = conf.PACS.Path + conf.PACS.AET + "/wado"
	dicomURL.ForceQuery = true
	query := url.Values{}
	query.Set("requestType", "WADO")
	query.Set("studyUID", dicomFile.StudyUID)
	query.Set("seriesUID", dicomFile.SeriesInstanceUID)
	query.Set("objectUID", dicomFile.SOPInstanceUID)

	dicomURL.RawQuery = query.Encode()
	return dicomURL.String() + "&contentType=application/dicom"
}

func GetStudyInstancesURL(studyUID string, conf configuration.Configuration) string {
	dicomURL := url.URL{}
	dicomURL.Scheme = "http"
	dicomURL.Host = conf.PACS.Host
	dicomURL.Path = conf.PACS.Path + conf.PACS.AET + "/rs/studies/" + studyUID + "/metadata"

	return dicomURL.String()
}
