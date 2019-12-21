package rest

import (
	"github.com/Bigyin1/GoMobileBackend/pkg/crypter"
	"github.com/gorilla/mux"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

type filesResource struct {
	cryptService *crypter.Service
}

func (fr *filesResource) processMultipart(reader *multipart.Reader) (map[string][]byte, error) {
	files := make(map[string][]byte)
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		name := part.FormName()
		if part.FileName() == "" {
			name = part.FormName()
		}
		var buf []byte
		_, err = io.ReadFull(part, buf)
		if err != nil {
			return nil, err
		}
		log.Printf("Got data with name: %s\n", name)
		files[name] = buf
	}
	return files, nil
}

func (fr *filesResource) Post(w http.ResponseWriter, r *http.Request) {
	log.Printf("request: %s %s\n", r.RequestURI, r.Method)
	reader, err := r.MultipartReader()
	if err == http.ErrNotMultipart {
		log.Printf("Got simple request type: %s", r.Header.Get("Content-Type"))
		return
	}
	if err != nil {
		log.Println(getErrorStr(r, err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Got multipart request")
	files, err := fr.processMultipart(reader)
	if err != nil {
		log.Printf("Failed request: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	mapping := fr.cryptService.EncryptAndSaveFiles(files)
	writeMappingResponse(mapping, http.StatusOK, w)
}

func (fr *filesResource) Get(w http.ResponseWriter, r *http.Request) {
	log.Printf("request: %s %s\n", r.RequestURI, r.Method)

	fid := mux.Vars(r)["fid"]
	key := r.FormValue("key")
	log.Printf("Quering file %s with key %s", fid, key)
	file, err := fr.cryptService.DecryptFile(fid, key)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	writeFileResponse(file, w)
}
