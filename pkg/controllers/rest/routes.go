package rest

import (
	"github.com/Bigyin1/GoMobileBackend/pkg/crypter"
	"github.com/gorilla/mux"
)

func BuildRestApiRouter(cryptService *crypter.Service, debug bool) *mux.Router {
	r := mux.NewRouter()
	fileRes := filesResource{cryptService: cryptService}
	r.HandleFunc("/upload", loggingMiddleware(errorWrapperMiddleware(fileRes.Post, debug))).
		Methods("POST", "OPTIONS")
	r.HandleFunc("/file/{fid}", loggingMiddleware(errorWrapperMiddleware(fileRes.Get, debug))).
		Methods("GET", "OPTIONS")
	return r
}
