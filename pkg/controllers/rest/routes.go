package rest

import (
	"github.com/Bigyin1/GoMobileBackend/pkg/crypter"
	"github.com/gorilla/mux"
)

func BuildRestApiRouter(cryptService *crypter.Service) *mux.Router {
	r := mux.NewRouter()
	fileRes := filesResource{cryptService: cryptService}
	r.HandleFunc("/file", fileRes.Post).Methods("POST")
	r.HandleFunc("/file/{fid}", fileRes.Get).Methods("GET")
	return r
}
