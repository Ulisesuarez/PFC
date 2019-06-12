package main

import (
	"bitbucket.org/inehealth/api/api"
	commonMiddleware "bitbucket.org/inehealth/api/middleware"
	apiservices "bitbucket.org/inehealth/api/services"
	"bitbucket.org/inehealth/idonia-pacs/configuration"
	"bitbucket.org/inehealth/idonia-pacs/controller"
	"bitbucket.org/inehealth/idonia-pacs/middleware"
	"bitbucket.org/inehealth/idonia-pacs/repository"
	"bitbucket.org/inehealth/idonia-pacs/service/db"
	"bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia"
	idonia_auth "bitbucket.org/inehealth/idonia-pacs/service/idonia/idonia-auth"
	"bitbucket.org/inehealth/idonia-pacs/service/task"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"
)

const defaultConfigFilePath = "./config.json"

var IdoniaPacsVersion string
var IdoniaPacsCompilation string
var IdoniaImplementationClassUID = "1.2.826.0.1.3680043.10.154.1.1"

//go:generate go get -u github.com/jteeuwen/go-bindata/...
//go:generate git submodule init
//go:generate git submodule update --recursive --remote
//go:generate npm ci --prefix ./idonia-pacs-ui
//go:generate npm run build --prefix ./idonia-pacs-ui
//go:generate go-bindata -pkg ui -o ui/ui.go -prefix "idonia-pacs-ui/dist/" idonia-pacs-ui/dist/...
//go:generate echo "UI Built! hash is: $(git submodule status)"

func welcomeBanner(config *configuration.Configuration) {
	fmt.Printf("IDONIA API HOST: %s\n", config.API.REST.HostName)
	fmt.Printf("IDONIA API KEY: %s\n", config.API.REST.SSLPublicKey)
	fmt.Printf("IDONIA PACS API: http%s://%s:%d%s\n", func() string {
		if len(config.API.REST.SSLPrivateKey) > 0 {
			return "s"
		} else {
			return ""
		}
	}(), config.API.REST.Host, config.API.REST.Port, config.API.REST.Prefix)
	fmt.Printf("IDONIA PACS UI: http%s://%s:%d\n", func() string {
		if len(config.API.REST.SSLPrivateKey) > 0 {
			return "s"
		} else {
			return ""
		}
	}(), config.API.REST.Host, config.API.REST.Port)
	fmt.Printf("Everything is ready!\n")
}

func main() {
	fmt.Printf("Welcome to Idonia PACS!\n")
	fmt.Printf("© INEHEALTH TEAM S.L.\n")
	fmt.Printf("Initialising, this may take a few minutes...\n")

	versionFlag := flag.Bool("version", false, "Show version")
	configFlag := flag.String("configFile", defaultConfigFilePath, "Path where config file is found")
	flag.Parse()

	if versionFlag != nil && *versionFlag {
		println("Version: ", IdoniaPacsVersion)
		println("Compilation: ", IdoniaPacsCompilation)
		os.Exit(0)
	}

	config, err := configuration.NewConfiguration(*configFlag)
	if err != nil {
		println("[ config ] Config file not found or not readable")
		println("[ config ] ", err.Error())
		os.Exit(1)
	}

	loggingLevel := logrus.InfoLevel

	switch strings.ToLower(config.Logger.Level) {
	case "debug":
		loggingLevel = logrus.DebugLevel
	default:
		fallthrough
	case "info":
		loggingLevel = logrus.InfoLevel
	case "warn":
		loggingLevel = logrus.WarnLevel
	case "error":
		loggingLevel = logrus.ErrorLevel
	case "fatal":
		loggingLevel = logrus.FatalLevel
	case "panic":
		loggingLevel = logrus.PanicLevel
	}

	logger := logrus.New()
	logger.Level = loggingLevel

	loggerQuit := make(chan struct{})

	if config.Logger.Folder != "" {
		err = os.MkdirAll(config.Logger.Folder, 0770)
		if err != nil {
			panic("Log File can't be created")
		}
		currentDate := time.Now().Format("2006-01-02")
		logFile, err := os.OpenFile(path.Join(config.Logger.Folder, "idonia-pacs-"+currentDate+".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0775)
		if err != nil {
			panic("Log File can't be created or opened")
		}

		ticker := time.NewTicker(1 * time.Minute)
		go func() {
			for {
				select {
				case <-ticker.C:
					tmpDate := time.Now().Format("2006-01-02")
					if currentDate != tmpDate {
						currentDate = tmpDate

						newLogFile, err := os.OpenFile(path.Join(config.Logger.Folder, "idonia-pacs-"+currentDate+".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0775)
						if err != nil {
							break
						}

						logger.Out = newLogFile
						err = logFile.Close()

						logFile = newLogFile
					}
				case <-loggerQuit:
					err = logFile.Close()
					return
				}
			}
		}()

		logger.Out = logFile
	}

	loggerOutput := logger.Writer()

	log.SetOutput(loggerOutput)
	api.LoggerWriter = loggerOutput

	//Create a new server with preconfiguration

	run(config, logger)
	loggerOutput.Close()
	os.Exit(0)
}

func run(config *configuration.Configuration, logger *logrus.Logger) {
	data := url.Values{}
	data.Set("logicalName", "arc")
	data.Set("aeTitle", "DCM4CHEE")
	data.Set("hostName", "arc")
	data.Set("port", "11112")
	data.Set("retrieve", "WADO")
	data.Set("wadoContext", "dcm4chee-arc/aets/DCM4CHEE/wado")
	data.Set("wadoPort", "8080")
	data.Set("imageType", "JPEG")
	data.Set("previews", "true")
	data.Set("todo", "ADD")

	response, err := http.Post("http://0.0.0.0:3000/ServerConfig.do", "application/x-www-form-urlencoded", strings.NewReader((data.Encode())))

	server := api.NewAPI(config.API)
	server.Logger = logger
	server.Logger.Info(fmt.Printf("%+v", response))
	if err != nil {
		server.Logger.Error(err.Error())
		err = nil
	}
	if response != nil {
		server.Logger.Info(fmt.Printf("%+v", response.Request))
	}
	services := generateServices(config, logger)
	server.RegisterService(services...)
	server.RegisterMiddleware(generateMiddlewares(config, logger)...)
	server.RegisterModule(
		&controller.GenericModule{
			Version:     IdoniaPacsVersion,
			Compilation: IdoniaPacsCompilation,
		},
		&controller.ControllerHandler{},
		&controller.AuthHandler{},
		&controller.StudyHandler{},
		&controller.TaskHandler{})
	server.RegisterRepository(
		repository.Study{},
		repository.Task{},
		repository.TaskDicom{},
		repository.TaskAdditionalField{},
		repository.Session{},
	)

	err = server.Start()
	if err != nil {
		panic(err)
	}
	var engineEntry interface{}
	for _, service := range services {
		if service.Name == "engine" {
			engineEntry = service.Value
		}
	}
	if engineEntry == nil {
		panic("Unexpected task engine error")
	}
	engine, ok := engineEntry.(*task.Engine)
	if !ok {
		panic("Unexpected task engine error")
	}

	go engine.Start(4) //TODO: Read from config
	logger.Infof("[init] Task engine running")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	engine.Stop()
	server.Stop()
}

func generateServices(config *configuration.Configuration, logger *logrus.Logger) (services []*api.Service) {
	if config.Idonia != nil {
		idonia.AccountID = config.Idonia.ApiId
		idonia.APIHost = config.Idonia.ApiHost
		//idonia.ApiKey = config.Idonia.ApiKey
		//idonia.ApiSecret = config.Idonia.ApiSecret
		if strings.HasPrefix(idonia.ApiKey, "K1") {
			idonia_auth.ApiKeyLogin()
		}
	}

	database, err := db.NewDatabase(config.DB.Location)

	if err != nil {
		panic(err)
	}
	err = repository.DeleteAllSessions(database)
	if err != nil {
		panic(err)
	}
	credentialsFile := ""

	if credentialsFile = config.GCloudCredentials.ApplicationCredentials; credentialsFile == "" {
		credentialsFile = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	}

	if credentialsFile == "" {
		panic("No credentials file defined")
	}

	fileDataBucketID := ""

	if fileDataBucketID = config.FileDataStorage.Bucket; fileDataBucketID == "" {
		panic("No file data bucket defined")
	}
	fmt.Println(credentialsFile)
	gcloudSecurity, err := apiservices.NewGCSSecurity(credentialsFile, fileDataBucketID)
	if err != nil {
		panic(err.Error())
	}
	var engine *task.Engine

	// endregion

	// region Task Engine

	engine = task.NewEngine(database, logger, config)

	services = append(services,
		&api.Service{Name: "db", Value: database},
		&api.Service{Name: "configuration", Value: config},
		&api.Service{Name: "logger", Value: logger},
		&api.Service{Name: "gcsjj-security", Value: gcloudSecurity},
		&api.Service{Name: "engine", Value: engine},
	)
	return
}

func generateMiddlewares(config *configuration.Configuration, logger *logrus.Logger) (middlewares []api.Middleware) {
	middlewares = append(middlewares,
		middleware.Authorization{},
		commonMiddleware.PanicRecover{},
		commonMiddleware.RequestCleanUp{},
		commonMiddleware.CORS{},
		commonMiddleware.InitCorrelation{CorrelationKey: "X-Idonia-Correlation-ID", RequestKey: "X-Idonia-Request-ID"},
		commonMiddleware.LoggerRQ{},
	)
	return
}
