package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func encodeJSON(data interface{}) (string, error) {
	dataBytes, err := json.Marshal(data)
	return string(dataBytes), err
}

func writeResponse(response interface{}, status int, w http.ResponseWriter) {
	responseData, err := encodeJSON(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = fmt.Fprintln(w, responseData)
}

func writeFileResponse(file []byte, w http.ResponseWriter) {

	FileHeader := make([]byte, 512)
	r := bytes.NewReader(file)
	_, err := r.Read(FileHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	FileContentType := http.DetectContentType(FileHeader)

	FileSize := strconv.Itoa(len(file))

	//w.Header().Set("Content-Disposition", "attachment; filename="+Filename) TODO: store real file name somewhere
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	r.Seek(0, 0)
	_, err = io.Copy(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
