package mail

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/Bigyin1/GoMobileBackend/pkg/crypter"
	"github.com/palantir/stacktrace"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func renderMappingTmpl(mapping crypter.Mapping, tmplPath string) string {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Println(err)
		return ""
	}
	var buff bytes.Buffer
	err = tmpl.Execute(&buff, mapping.GetMapping())
	if err != nil {
		log.Println(err)
		return ""
	}
	return buff.String()
}

func renderOutputMessage(mapping crypter.Mapping, tmplPath, to, from, subject string) *gmail.Message {
	m := gmail.Message{}
	messageStr := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html;\r\n\r\n%s",
		from, to, subject,
		renderMappingTmpl(mapping, tmplPath)))
	//log.Println(string(messageStr))
	m.Raw = base64.URLEncoding.EncodeToString(messageStr)
	return &m
}

func getMessageHeader(m *gmail.Message, header string) string {
	for _, h := range m.Payload.Headers {
		if h.Name == header {
			return h.Value
		}
	}
	return ""
}

func getInitialHistoryId(gmailService *gmail.Service) (uint64, error) {
	m, err := gmailService.Users.Messages.List("me").LabelIds("INBOX").Do() //sync
	if err != nil {
		return 0, stacktrace.Propagate(err, "Failed to get messages to sync")
	}
	lastMessageid := m.Messages[0].Id
	lm, err := gmailService.Users.Messages.Get("me", lastMessageid).Format("metadata").Do()
	if err != nil {
		return 0, stacktrace.Propagate(err, "Failed to get message to sync")
	}
	return lm.HistoryId, nil
}

func isUploadSubject(m *gmail.Message, s string) bool {
	subject := getMessageHeader(m, "Subject")
	if subject == s {
		return true
	}
	return false
}

func logMessage(m *gmail.Message) string {
	return fmt.Sprintf("From: %s Subject: %s",
		getMessageHeader(m, "From"),
		getMessageHeader(m, "Subject"))
}
