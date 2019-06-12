package configuration

import (
	"bitbucket.org/inehealth/api/configuration"
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	API *configuration.API `json:"api"`
	DB  struct {
		Location string
	} `json:"db"`
	Logger struct {
		Level  string
		Folder string
	} `json:"logger"`
	PACS              *PACSConfiguration               `json:"pacs"`
	Login             *Login                           `json:"login"`
	GCloudCredentials *configuration.GCloudCredentials `json:"gcloudCredentials"`
	FileDataStorage   *configuration.GCloudStorage     `json:"fileDataStorage"`
	IMWPubSub         *configuration.GCloudPubSub      `json:"imwPubSub"`
	Idonia            *IdoniaConfiguration             `json:"idonia"`
}
type User struct {
	Username string
	Password string
}

type Login struct {
	Type              string
	ForgotPasswordURL string
	Users             []User
}
type IdoniaConfiguration struct {
	Namespace	string
	APiUsername string
	ApiPass     string
	ApiId       uint32
	ApiKey      string
	ApiSecret   string
	ApiHost     string
}

type PACSConfiguration struct {
	Host string
	Path string
	AET  string
}

//NewConfiguration ...
func NewConfiguration(filePath string) (*Configuration, error) {

	var file *os.File
	defer file.Close()

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	config := Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	if config.PACS == nil {
		config.PACS = &PACSConfiguration{
			Host: "localhost:8080",
			Path: "/dcm4chee-arc/aets/",
			AET:  "DCM4CHEE",
		}
	}
	if config.API == nil {
		config.API = &configuration.API{}
	}
	if config.API.REST == nil {
		config.API.REST = &configuration.REST{}
	}
	config.API.REST.Prefix = "/api"
	config.API.REST.HostName = ""
	config.API.REST.StrictSlash = true
	if config.FileDataStorage == nil {
		config.FileDataStorage = &configuration.GCloudStorage{}
	}

	config.FileDataStorage.Bucket = "whatever"
	if config.DB.Location == "" {
		config.DB.Location = "./idonia-pacs.db"
	}
	if config.Logger.Level == "" {
		config.Logger.Level = "info"
	}
	if config.Login == nil {
		fmt.Println("login not found")
		config.Login = &Login{}
		config.Login.Type = "None"
	}
	if config.Idonia.ApiHost == "" {
		config.Idonia.ApiHost = "https://api-staging.idonia.com/api"
	}

	return &config, err
}
