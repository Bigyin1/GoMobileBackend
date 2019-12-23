package rest

import (
	"bytes"
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

func (fr *filesResource) processMultipart(reader *multipart.Reader) (<-chan crypter.Mapping, crypter.Mapping, int) {
	outChan := make(chan crypter.Mapping)
	var num int

	unprocessed := crypter.NewFilesMapping()
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Failed to read next multipart part", err)
			continue
			//return nil, num, stacktrace.PropagateWithCode(err, ErrMultipartBadFormat, "multipart bad format")
		}
		name := part.FileName()
		if name == "" {
			name = part.FormName()
		}

		var b []byte
		buf := bytes.NewBuffer(b)
		_, err = io.Copy(buf, part)
		if err != nil {
			unprocessed.AddError(name, "Failed to read", "")
			log.Println("failed to read file", name)
			continue
			//return nil, num, stacktrace.PropagateWithCode(err, ErrMultipartProcessing, "failed to read next part")
		}
		file := make(crypter.InputFiles)
		file[name] = buf.Bytes()
		log.Println(part.FileName())

		num++
		fr.cryptService.EncryptAndSaveFilesAsync(file, outChan)
	}
	return outChan, unprocessed, num
}

func (fr *filesResource) Post(w http.ResponseWriter, r *http.Request) error {
	reader, err := r.MultipartReader()
	if err == http.ErrNotMultipart {
		return stacktrace.NewMessageWithCode(ErrNotMultipart, "accept only multipart/form-data Content-Type")
	}
	if err != nil {
		return stacktrace.PropagateWithCode(err, ErrMultipartProcessing, "failed to get multipart reader")
	}

	outChan, unprocessed, num := fr.processMultipart(reader)
	log.Println(num)
	mapping := fr.cryptService.WaitForAllFilesAndMergeMappings(outChan, num)
	mapping.MergeWith(unprocessed)

	writeResponse(mapping, http.StatusOK, w)
	return nil
}

func (fr *filesResource) Get(w http.ResponseWriter, r *http.Request) error {
	fid := mux.Vars(r)["fid"]
	key := r.FormValue("key")

	file, err := fr.cryptService.DecryptFileAndGet(fid, key)
	if err != nil {
		return stacktrace.Propagate(err, "failed to decrypt file %s with key %s", fid, key)
	}
	writeFileResponse(file, w)
	return nil
}
