package app

import (
	"github.com/Bigyin1/GoMobileBackend/config"
	"github.com/Bigyin1/GoMobileBackend/pkg/controllers/mail"
	"github.com/Bigyin1/GoMobileBackend/pkg/controllers/rest"
	"github.com/Bigyin1/GoMobileBackend/pkg/crypter"
	"github.com/Bigyin1/GoMobileBackend/pkg/infrastructure"
	"log"
	"net/http"
	"strconv"
	"time"
)

type App struct {
	restServer      *http.Server
	gmailController *mail.GmailController
}

func (app *App) StartApp() {
	go app.gmailController.StartPolling()
	err := app.restServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Cann't start serving rest, check port num %s", app.restServer.Addr)
	}
}

func InitApp() *App {
	config.LoadEnvironment()
	conf := config.LoadConfig(config.ConfPath)
	log.Printf("Current configuration: %s\n", conf.AsString())

	fileRepo := infrastructure.NewInFsFileStorage(conf.StoragePath)
	cryptService := crypter.NewCrypterService(fileRepo, conf.FileURIPrefix, crypter.GetRandomEncrKey)
	restRouter := rest.BuildRestApiRouter(cryptService, conf.Debug)
	restServer := &http.Server{
		Handler:      restRouter,
		Addr:         ":" + strconv.Itoa(conf.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	gmailCtrl := mail.NewGmailController(conf.GmailTokenPath,
		conf.GmailCredsPath,
		conf.UploadSubject,
		conf.GmailAddr,
		conf.HistoryIdPath,
		conf.MailTmplPath,
		conf.PollingPeriod,
		cryptService)

	return &App{restServer: restServer, gmailController: gmailCtrl}
}
