package mail

import (
	"context"
	"encoding/base64"
	"github.com/Bigyin1/GoMobileBackend/pkg/crypter"
	"github.com/palantir/stacktrace"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"time"
)

type GmailController struct {
	gmailService *gmail.Service
	cryptService *crypter.Service
	historyID    uint64
	pollingPeriod int
	uploadSubject string
	historyRequest *gmail.UsersHistoryListCall
	gmailAddr string
}

func (gc *GmailController) sendOutputMapping(output crypter.Mapping, to string) {
	_, err := gc.gmailService.Users.Messages.Send("me",
		getOutputMessage(output, to, gc.gmailAddr, "Mapping")).Do()
	if err != nil {
		log.Println("error during email send:", err.Error())
		return
	}
	log.Println("Successfully sent mapping to", to)
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
	inputFiles := make(map[string][]byte)
	for _, part := range message.Payload.Parts {
		if part.Filename != "" {
			fileData, err := gc.getFilePartData(part, message.Id)
			if err != nil {
				log.Printf("Failed to get file %s %s", part.Filename, err.Error())
				continue
			}
			if len(fileData) == 0 {
				log.Printf("Got file with zero length: %s", part.Filename)
				continue
			}
			inputFiles[part.Filename] = fileData
		}
	}
	mapping := gc.cryptService.EncryptAndSaveFiles(inputFiles)
	sendTo := getMessageHeader(message, "From")
	go gc.sendOutputMapping(mapping, sendTo)
}

func (gc *GmailController) processHistory(history *gmail.ListHistoryResponse) uint64 {
	lastHistoryId := gc.historyID
	for i, h := range history.History {
		if i == len(history.History) - 1 {
			lastHistoryId = h.Id
		}
		for _, m := range h.MessagesAdded {
			m, err := gc.gmailService.Users.Messages.Get("me", m.Message.Id).Do()
			if err != nil {
				log.Println("Failed to get message", err.Error())
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
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		history, err := gc.getHistory()
		if err != nil {
			log.Println("get history request failed", err.Error())
			continue
		}
		newHistoryId := gc.processHistory(history)
		gc.updateHistoryID(newHistoryId)
	}
}

func NewGmailController(tokenPath, credsPath, uploadSubject, gmailAddr string, crServ *crypter.Service) *GmailController {
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
		log.Fatalln("Failed to read gmail token", err.Error())
	}
	gmailService, err := gmail.NewService(context.TODO(), option.WithTokenSource(config.TokenSource(context.TODO(), token)))
	if err != nil {
		log.Fatalln("Failed to init gmail service", err)
	}
	historyReq := gmailService.Users.History.List("me").LabelId("INBOX").HistoryTypes("messageAdded")
	return &GmailController{gmailService: gmailService,
		cryptService:crServ,
		uploadSubject:uploadSubject,
		gmailAddr:gmailAddr,
		historyRequest:historyReq,
		historyID:1928}
}
