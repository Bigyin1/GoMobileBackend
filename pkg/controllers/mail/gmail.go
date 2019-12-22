package main

import (
	"context"
	"fmt"
	"github.com/Bigyin1/GoMobileBackend/pkg/crypter"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
)

type GmailController struct {
	gmailService *gmail.Service
	cryptService *crypter.Service
	historyID    int64
}

func (gc *GmailController) isUploadSubject(headers []*gmail.MessagePartHeader) bool {
	for _, h := range headers {
		if h.Name == "Subject" && h.Value == "upload" {
			return true
		} else {
			return false
		}
	}
	return false
}

func (gc *GmailController) updateHistoryID(newHistoryID int64) {
	gc.historyID = newHistoryID
}

func NewGmailController(tokenPath, credsPath string) *GmailController {
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
	return &GmailController{gmailService: gmailService}
}

func main() {
	ctrl := NewGmailController("./token/token.json", "./token/credentials.json")
	user := "me"
	r, _ := ctrl.gmailService.Users.History.List(user).StartHistoryId(1928).LabelId("INBOX").HistoryTypes("messageAdded").Do()
	//if err != nil {
	//	log.Fatalf("Unable to retrieve labels: %v", err)
	//}
	//if len(r.Messages) == 0 {
	//	fmt.Println("No labels found.")
	//	return
	//}
	//fmt.Println("M:")
	//for _, m := range r.Messages {
	//	m, _ := ctrl.gmailService.Users.Messages.Get(user, m.Id).Do()
	//	//fmt.Println(m.Snippet)
	//	fmt.Printf("%+v\n", *m)
	//	//for _, p := range m.Payload.Parts {
	//	//	res, _ := base64.URLEncoding.DecodeString(p.Body.Data)
	//	//	fmt.Println(string(res))
	//	//}
	//}
	//fmt.Printf("%+v\n", r)
	for _, h := range r.History {
		for _, m := range h.MessagesAdded {
			m, _ := ctrl.gmailService.Users.Messages.Get(user, m.Message.Id).Do()
			for _, h := range m.Payload.Headers {
				fmt.Println(h.Value)
			}
		}
	}
	//m, _ := ctrl.gmailService.Users.Messages.Get(user, m.).Do()
	//fmt.Printf("%+v\n", *m)
	//fmt.Println(m.HistoryId)
	//for _, part := range m.Payload.Parts {
	//	if part.Filename != "" {
	//		fileData , _ := base64.URLEncoding.DecodeString(part.Body.Data)
	//		f, _ := os.Create(part.Filename)
	//		f.Write(fileData)
	//		f.Close()
	//	}
	//}

}
