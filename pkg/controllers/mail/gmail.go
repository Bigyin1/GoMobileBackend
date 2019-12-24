package mail

import (
	"bytes"
	"context"
	"encoding/base64"
	"io/ioutil"
	"log"
	"time"

	"github.com/Bigyin1/GoMobileBackend/pkg/crypter"
	"github.com/palantir/stacktrace"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type GmailController struct {
	gmailService   *gmail.Service
	cryptService   *crypter.Service
	historyID      uint64
	pollingPeriod  int
	uploadSubject  string
	historyRequest *gmail.UsersHistoryListCall
	gmailAddr      string
	mailTmplPath   string
}

func (gc *GmailController) sendOutputMapping(output crypter.Mapping, to string) {
	_, err := gc.gmailService.Users.Messages.Send("me",
		renderOutputMessage(output, gc.mailTmplPath, to, gc.gmailAddr, "Mapping")).Do()
	if err != nil {
		log.Println("error during email send:", err)
		return
	}
	log.Println("Successfully sent mapping to", to)
}

func (gc *GmailController) getFilePartData(part *gmail.MessagePart, mid string) ([]byte, error) {
	if part.Body.AttachmentId == "" {
		fileData, err := base64.URLEncoding.DecodeString(part.Body.Data)
		if err != nil {
			return nil, stacktrace.Propagate(err, "Failed to decode part.Body.Data")
		}
		return fileData, nil
	}
	att, err := gc.gmailService.Users.Messages.Attachments.Get("me", mid, part.Body.AttachmentId).Do()
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to get attachment")
	}
	fileData, err := base64.URLEncoding.DecodeString(att.Data)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to decode attachment.Data")
	}
	return fileData, nil
}

func (gc *GmailController) processMessage(message *gmail.Message) {
	log.Println("Start processing email message:", logMessage(message))
	inputFiles := make(crypter.InputFiles)
	for _, part := range message.Payload.Parts {
		if part.Filename != "" {
			fileData, err := gc.getFilePartData(part, message.Id)
			if err != nil {
				log.Printf("Failed to get file %s %s", part.Filename, err)
				continue
			}
			if len(fileData) == 0 {
				log.Printf("Got file with zero length: %s", part.Filename)
				continue
			}
			dataReader := bytes.NewReader(fileData)
			log.Println(part.Filename)
			inputFiles[part.Filename] = dataReader
		}
	}
	mapping := gc.cryptService.EncryptAndSaveFiles(inputFiles)
	sendTo := getMessageHeader(message, "From")
	go gc.sendOutputMapping(mapping, sendTo)
}

func (gc *GmailController) updateHistoryID(newHistoryID uint64) {
	gc.historyID = newHistoryID
}

func (gc *GmailController) getHistory() (*gmail.ListHistoryResponse, error) {
	h, err := gc.historyRequest.StartHistoryId(gc.historyID).Do()
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to get history with id: %d", gc.historyID)
	}
	return h, nil
}

func (gc *GmailController) processHistory(history *gmail.ListHistoryResponse) uint64 {
	lastHistoryId := gc.historyID
	for i, h := range history.History {
		if i == len(history.History)-1 {
			lastHistoryId = h.Id
		}
		for _, m := range h.MessagesAdded {
			log.Println("Start getting new email message")
			m, err := gc.gmailService.Users.Messages.Get("me", m.Message.Id).Do()
			if err != nil {
				log.Println("Failed to get message", err)
				continue
			}
			log.Println("Got new email message:", logMessage(m))
			if isUploadSubject(m, gc.uploadSubject) {
				go gc.processMessage(m)
			}
		}
	}
	return lastHistoryId

}

func (gc *GmailController) StartPolling() {
	ticker := time.NewTicker(time.Duration(gc.pollingPeriod) * time.Second)
	for range ticker.C {
		history, err := gc.getHistory()
		if err != nil {
			log.Println("get history request failed", err)
			continue
		}
		newHistoryId := gc.processHistory(history)
		if newHistoryId != gc.historyID {
			gc.updateHistoryID(newHistoryId)
		}
	}
}

func NewGmailController(tokenPath, credsPath, uploadSubject, gmailAddr,
	mailTmplPath string,
	pollPeriod int,
	crServ *crypter.Service) *GmailController {
	b, err := ioutil.ReadFile(credsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	token, err := tokenFromFile(tokenPath)
	if err != nil {
		log.Fatalln("Failed to read gmail token", err)
	}
	gmailService, err := gmail.NewService(context.TODO(), option.WithTokenSource(config.TokenSource(context.TODO(), token)))
	if err != nil {
		log.Fatalln("Failed to init gmail service", err)
	}
	historyId, err := getInitialHistoryId(gmailService)
	if err != nil {
		log.Fatalln("Failed to get initial history id", err)
	}
	historyReq := gmailService.Users.History.List("me").LabelId("INBOX").HistoryTypes("messageAdded")
	return &GmailController{gmailService: gmailService,
		cryptService:   crServ,
		uploadSubject:  uploadSubject,
		gmailAddr:      gmailAddr,
		historyRequest: historyReq,
		historyID:      historyId,
		pollingPeriod:  pollPeriod,
		mailTmplPath:   mailTmplPath,
	}
}
