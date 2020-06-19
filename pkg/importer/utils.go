package importer

import (
	"archives/pkg/config"
	"archives/pkg/database"
	"archives/pkg/models"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"time"
)

func importMail(name, path, maildirPath string) {
	file, _ := os.Open(path)
	m, _ := mail.ReadMessage(file)

	msg := models.Message{
		Id:          m.Header.Get("X-Archives-Hash"),
		Filename:    name,
		Headers:     m.Header,
		Attachments: nil,
		Body:        getBody(m.Header, m.Body),
		Date:        getDate(m.Header),
		Lists:       getLists(m.Header),
		List:        getListName(path),
		Comment:     "",
		Hidden:      false,
	}

	err := insertMessage(msg)

	if err != nil {
		fmt.Println("Error during importing Mail")
		fmt.Println(err)
	}
}

func getDepth(path, maildirPath string) int {
	return strings.Count(strings.ReplaceAll(path, maildirPath, ""), "/")
}

func getBody(header mail.Header, body io.Reader) map[string]string {
	if isMultipartMail(header) {
		boundary := regexp.MustCompile(`boundary="(.*?)"`).
			FindStringSubmatch(
				header.Get("Content-Type"))
		if len(boundary) != 2 {
			//err
			return map[string]string{
				"text/plain": "",
			}
		}
		return getBodyParts(body, boundary[1])
	} else {
		content, _ := ioutil.ReadAll(body)
		return map[string]string{
			getContentType(header): string(content),
		}
	}
}

func getBodyParts(body io.Reader, boundary string) map[string]string {
	bodyParts := make(map[string]string)
	mr := multipart.NewReader(body, boundary)
	for {
		p, err := mr.NextPart()
		if err != nil {
			return bodyParts
		}
		slurp, err := ioutil.ReadAll(p)
		if err != nil {
			log.Fatal(err)
		}
		bodyParts[p.Header.Get("Content-Type")] = string(slurp)
	}
	return bodyParts
}

func getContentType(header mail.Header) string {
	contentTypes := regexp.MustCompile(`(.*?);`).
		FindStringSubmatch(
			header.Get("Content-Type"))
	if len(contentTypes) < 2 {
		// assume text/plain if we don't find a Content-Type header e.g. for git patches
		return "text/plain"
	}
	return contentTypes[1]
}

func getDate(header mail.Header) time.Time {
	date, _ := header.Date()
	return date
}

func isMultipartMail(header mail.Header) bool {
	return strings.Contains(getContentType(header), "multipart")
}

func getLists(header mail.Header) []string {
	var lists []string
	// To
	adr, _ := mail.ParseAddressList(header.Get("To"))
	for _, v := range adr {
		lists = append(lists, v.Address)
	}
	// Cc
	adr, _ = mail.ParseAddressList(header.Get("Cc"))
	for _, v := range adr {
		lists = append(lists, v.Address)
	}
	return lists
}

func getListName(path string) string {
	listName := strings.ReplaceAll(path, config.MailDirPath() + ".", "")
	listName = strings.Split(listName, "/")[0]
	return listName
}

func insertMessage(message models.Message) error {
	_, err := database.DBCon.Model(&message).
		Value("tsv_subject", "to_tsvector(?)", message.GetSubject()).
		Value("tsv_body", "to_tsvector(?)", message.GetBody()).
		OnConflict("(id) DO NOTHING").
		Insert()
	return err
}

func isPublicList(path string) bool {
	for _, publicList := range config.AllPublicMailingLists(){
		if publicList == getListName(path) {
			return true
		}
	}
	return false
}
