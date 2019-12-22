package rest

import (
	"github.com/Bigyin1/GoMobileBackend/pkg/crypter"
	"github.com/Bigyin1/GoMobileBackend/pkg/infrastructure"
	"github.com/palantir/stacktrace"
	"log"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request: %s %s\n", r.RequestURI, r.Method)
		start := time.Now()
		next.ServeHTTP(w, r)
		since := time.Since(start)
		log.Printf("request done: %s %s in %s \n", r.RequestURI, r.Method, since)
	}
}

type respError struct {
	Error string `json:"error,omitempty"`
}

func errorWrapperMiddleware(next ResourceHandler, isDebug bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}

		if r.Method == "OPTIONS" {
			return
		}
		err := next(w, r)

		if err == nil {
			return
		}
		if isDebug {
			log.Println("request failed with error", err.Error())
		} else {
			log.Println("request failed with error", stacktrace.RootCause(err).Error())
		}

		errResp := respError{Error: stacktrace.RootCause(err).Error()}

		switch stacktrace.GetCode(err) {
		case ErrMultipartProcessing:
			writeResponse(errResp, http.StatusInternalServerError, w)
		case ErrNotMultipart:
			writeResponse(errResp, http.StatusBadRequest, w)
		case ErrMultipartBadFormat:
			writeResponse(errResp, http.StatusBadRequest, w)
		case crypter.ErrUnexpected:
			writeResponse(errResp, http.StatusInternalServerError, w)
		case crypter.ErrWrongKey:
			writeResponse(errResp, http.StatusBadRequest, w)
		case infrastructure.ErrFileNotFound:
			writeResponse(errResp, http.StatusBadRequest, w)
		case infrastructure.ErrUnexpected:
			writeResponse(errResp, http.StatusInternalServerError, w)
		}
	}
}
