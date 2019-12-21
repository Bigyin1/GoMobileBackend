package rest

import (
	"github.com/Bigyin1/GoMobileBackend/pkg/crypter"
	"github.com/gorilla/mux"
)

func BuildRestApiRouter(cryptService *crypter.Service, debug bool) *mux.Router {
	r := mux.NewRouter()
	fileRes := filesResource{cryptService: cryptService}
	r.HandleFunc("/file", errorWrapperMiddleware(fileRes.Post, debug))
	r.HandleFunc("/file/{fid}", errorWrapperMiddleware(fileRes.Get, debug)).Methods("GET")
	return r
}
