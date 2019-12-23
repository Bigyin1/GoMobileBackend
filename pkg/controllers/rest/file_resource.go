package rest

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/Bigyin1/GoMobileBackend/pkg/crypter"
	"github.com/gorilla/mux"
	"github.com/palantir/stacktrace"
)

type ResourceHandler func(w http.ResponseWriter, r *http.Request) error

type filesResource struct {
	cryptService *crypter.Service
}

func (fr *filesResource) processMultipart(reader *multipart.Reader) (crypter.InputFiles, error) {
	files := make(crypter.InputFiles)
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, stacktrace.PropagateWithCode(err, ErrMultipartBadFormat, "multipart bad format")
		}

		name := part.FileName()
		if name == "" {
			name = part.FormName()
		}
		files[name] = part
		log.Println(part.FileName())
	}
	return files, nil
}

func (fr *filesResource) Post(w http.ResponseWriter, r *http.Request) error {
	reader, err := r.MultipartReader()
	if err == http.ErrNotMultipart {
		return stacktrace.NewMessageWithCode(ErrNotMultipart, "accept only multipart/form-data Content-Type")
	}
	if err != nil {
		return stacktrace.PropagateWithCode(err, ErrMultipartProcessing, "failed to get multipart reader")
	}

	files, err := fr.processMultipart(reader)
	if err != nil {
		return stacktrace.Propagate(err, "processMultipart failed")
	}

	mapping := fr.cryptService.EncryptAndSaveFiles(files)
	writeResponse(mapping, http.StatusOK, w)
	return nil
}

func (fr *filesResource) Get(w http.ResponseWriter, r *http.Request) error {
	fid := mux.Vars(r)["fid"]
	key := r.FormValue("key")

	err := fr.cryptService.DecryptFileAndGet(fid, key, w)
	if err != nil {
		return stacktrace.Propagate(err, "failed to decrypt file %s with key %s", fid, key)
	}
	//writeFileResponse(file, w)
	return nil
}
