package app

import (
	"github.com/Bigyin1/GoMobileBackend/config"
	"github.com/Bigyin1/GoMobileBackend/pkg/controllers/rest"
	"github.com/Bigyin1/GoMobileBackend/pkg/crypter"
	"github.com/Bigyin1/GoMobileBackend/pkg/infrastructure"
	"log"
	"net/http"
	"strconv"
	"time"
)

type App struct {
	restServer *http.Server
}

func (app *App) StartApp() {
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
	cryptService := crypter.NewCrypterService(fileRepo, conf.FileURIPrefix)
	restRouter := rest.BuildRestApiRouter(cryptService)
	restServer := &http.Server{
		Handler:      restRouter,
		Addr:         ":" + strconv.Itoa(conf.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return &App{restServer: restServer}
}
