package rest

import (
	"github.com/Bigyin1/GoMobileBackend/pkg/crypter"
	"github.com/gorilla/mux"
)

func BuildRestApiRouter(cryptService *crypter.Service) *mux.Router {
	r := mux.NewRouter()
	fileRes := filesResource{cryptService: cryptService}
	r.HandleFunc("/file", errorWrapperMiddleware(fileRes.Post, true))
	r.HandleFunc("/file/{fid}", errorWrapperMiddleware(fileRes.Get, true)).Methods("GET")
	return r
}
