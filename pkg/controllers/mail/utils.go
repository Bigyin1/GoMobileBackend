package mail

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Bigyin1/GoMobileBackend/pkg/crypter"
	"github.com/palantir/stacktrace"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"io/ioutil"
	"os"
	"strconv"
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

func outMappingToStr(mapping crypter.Mapping) string {
	res := ""
	for _, f := range mapping {
		res += fmt.Sprintf("%s ---> %s\n", f.Name, f.GetUrlOrErr())
	}
	return res
}

func getOutputMessage(mapping crypter.Mapping, to, from, subject string) *gmail.Message {
	m := gmail.Message{}
	messageStr := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from, to, subject, outMappingToStr(mapping)))
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

func getInitialHistoryId(path string) (uint64, error) {
	idstr, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, stacktrace.Propagate(err, "Failed to read history id file")
	}
	id, err := strconv.ParseUint(string(idstr), 10, 64)
	if err != nil {
		return 0, stacktrace.Propagate(err, "Wrong data in history id file")
	}
	return id, nil
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


