package importer

import (
	"archives/pkg/config"
	"archives/pkg/database"
	"archives/pkg/models"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type MailIdentifier struct {
	ArchivesHash string
	MessageId string
	To string
}

// TODO
var mails []*models.Message

// TODO
var WaitGroup sync.WaitGroup

// TODO
func initImport(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.IsDir() && getDepth(path, config.MailDirPath()) >= 1 && isPublicList(path) {

		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			return err
		}

		m, err := mail.ReadMessage(file)
		if err != nil {
			return err
		}

		references := strings.FieldsFunc(m.Header.Get("References"), func(r rune) bool {
			return r == ',' || r == ' '
		})
		mails = append(mails, &models.Message{
			Id:           m.Header.Get("X-Archives-Hash"),
			Filename:     info.Name(),
			From:         m.Header.Get("From"),
			To:           strings.Split(m.Header.Get("To"), ","),
			Subject:      m.Header.Get("Subject"),
			MessageIdField:    m.Header.Get("Message-Id"),
			RawReferences: references,
		})
	}
	return nil
}

// TODO
func importMail(path, filename string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	r := bytes.NewReader(content)
	m, err := mail.ReadMessage(r)
	if err != nil {
		return err
	}

	WaitGroup.Add(1)
	go importIntoDatabase(path, filename, m)

	return nil
}

func importIntoDatabase(path, filename string, m *mail.Message) {
	bodyContent, attachments := parseBody(m)
	msg := models.Message{
		Id:          m.Header.Get("X-Archives-Hash"),
		MessageIdField:   m.Header.Get("Message-Id"),
		Filename:    filename,
		From:        m.Header.Get("From"),
		To:          parseAddressList(m.Header.Get("To")),
		Cc:          parseAddressList(m.Header.Get("Cc")),
		Subject:     m.Header.Get("Subject"),

		List:        getListName(path),

		// TODO
		Date:        getDate(m.Header),
		InReplyToId:   getInReplyToMail(m.Header.Get("In-Reply-To"), m.Header.Get("From")),
		//References:  getReferencesToMail(strings.Split(m.Header.Get("References"), ","), m.Header.Get("From")),
		Body:        bodyContent,
		Attachments: attachments,

		StartsThread: m.Header.Get("In-Reply-To") == "" && m.Header.Get("References") == "",

		Comment:     "",
		Hidden:      false,
	}

	err := insertMessage(msg)
	if err != nil {
		fmt.Println("Error during importing Mail")
		fmt.Println(err)
	}

	references := strings.FieldsFunc(m.Header.Get("References"), func(r rune) bool {
		return r == ',' || r == ' '
	})
	insertReferencesToMail(references, m.Header.Get("X-Archives-Hash"), m.Header.Get("From"))

	WaitGroup.Done()
}

func parseAddressList(addressList string) []string {
	result := strings.Split(addressList, ",")
	if len(result) == 1 && strings.TrimSpace(result[0]) == "" {
		return nil
	}
	return result
}

func getInReplyToMail(messageId, from string) string {
	if messageId == "" {
		return ""
	}
	// step 1 TODO add description
	for _, mail := range mails {
		if mail.MessageId() == messageId && strings.Contains(strings.Join(mail.To, ", "), from) {
			return mail.Id
		}
	}
	// step 2 TODO add description
	for _, mail := range mails {
		if mail.MessageId() == messageId {
			return mail.Id
		}
	}
	return ""
}


func insertReferencesToMail(references []string, messageId, from string) []*models.Message {
	var referencesToMail []*models.Message
	for _, reference := range references {
		if strings.TrimSpace(reference) == "" {
			continue
		}
		// step 1 TODO add description
		for _, mail := range mails {
			if mail.MessageId() == reference  && strings.Contains(strings.Join(mail.To, ", "), from) {
				referencesToMail = append(referencesToMail, mail)
			}
		}
		// step 2 TODO add description
		for _, mail := range mails {
			if mail.MessageId() == reference {
				referencesToMail = append(referencesToMail, mail)
			}
		}
	}

	for _, reference := range referencesToMail {
		_, err := database.DBCon.Model(&models.MessageToReferences{
			MessageId: messageId,
			ReferenceId: reference.Id,
		}).Insert()

		if err != nil {
			fmt.Println("Err inserting Message to references")
			fmt.Println(err)
		}
	}

	return referencesToMail
}

func getDepth(path, maildirPath string) int {
	return strings.Count(strings.ReplaceAll(path, maildirPath, ""), "/")
}

func parseBody(m *mail.Message) (string, []models.Attachment) {
	header := m.Header
	body := m.Body
	foundPlainText := false
	if isMultipartMail(header) {
		var attachments []models.Attachment
		boundary := regexp.MustCompile(`boundary="(.*?)"`).
			FindStringSubmatch(
				header.Get("Content-Type"))
		if len(boundary) != 2 {
			//err
			return "", attachments
		}
		parsedBody := ""
		mr := multipart.NewReader(body, boundary[1])
		for {
			p, err := mr.NextPart()
			if err != nil {
				return parsedBody, attachments
			}
			bodyContent, err := ioutil.ReadAll(p)
			if err != nil {
				fmt.Println("Error while reading the body:")
				fmt.Println(err)
				continue
			}
			if strings.Contains(p.Header.Get("Content-Type"), "text/plain") {
				parsedBody = string(bodyContent)
				foundPlainText = true
			} else if strings.Contains(p.Header.Get("Content-Type"), "text/html") {
				if !foundPlainText {
					parsedBody = string(bodyContent)
				}
			} else if strings.Contains(p.Header.Get("Content-Type"), "multipart") {
				////////////////////////////////////////////////////////////////////////////////////////////////////////
				////////////////////////////////////////////////////////////////////////////////////////////////////////
				boundary := regexp.MustCompile(`boundary="(.*?)"`).
					FindStringSubmatch(p.Header.Get("Content-Type"))
				if len(boundary) != 2 {
					//err
					continue
				}
				mr := multipart.NewReader(bytes.NewReader(bodyContent), boundary[1])
				for {
					p, err := mr.NextPart()
					if err != nil {
						break
					}
					bodyContent, err := ioutil.ReadAll(p)
					if err != nil {
						fmt.Println("Error while reading the body:")
						fmt.Println(err)
						continue
					}
					if strings.Contains(p.Header.Get("Content-Type"), "text/plain") {
						parsedBody = string(bodyContent)
					} else if strings.Contains(p.Header.Get("Content-Type"), "text/html") {
						if !foundPlainText {
							parsedBody = string(bodyContent)
						}
					} else {
						attachments = append(attachments, models.Attachment{
							Filename: getAttachmentFileName(p.Header.Get("Content-Type")),
							Mime:     p.Header.Get("Content-Type"),
							Content:  string(bodyContent),
						})
					}
				}
				////////////////////////////////////////////////////////////////////////////////////////////////////////
				////////////////////////////////////////////////////////////////////////////////////////////////////////
			} else {
				attachments = append(attachments, models.Attachment{
					Filename: getAttachmentFileName(p.Header.Get("Content-Type")),
					Mime:     p.Header.Get("Content-Type"),
					Content:  string(bodyContent),
				})
			}
		}
		return parsedBody, attachments
	} else {
		content, _ := ioutil.ReadAll(body)
		return string(content), nil
	}
}


func getAttachments(header mail.Header, body io.Reader) []models.Attachment {

	if !isMultipartMail(header) {
		return nil
	}

	boundary := regexp.MustCompile(`boundary="(.*?)"`).
		FindStringSubmatch(
			header.Get("Content-Type"))
	if len(boundary) != 2 {
		return nil
	}
	var attachments []models.Attachment
	mr := multipart.NewReader(body, boundary[1])
	for {
		p, err := mr.NextPart()
		if err != nil {
			return attachments
		}
		content, err := ioutil.ReadAll(p)
		if err != nil {
			fmt.Println("Error while reading the body:")
			fmt.Println(err)
			continue
		}

		attachments = append(attachments, models.Attachment{
			Filename: getAttachmentFileName(p.Header.Get("Content-Type")),
			Mime:     p.Header.Get("Content-Type"),
			Content:  string(content),
		})

	}
	return attachments
}

func getAttachmentFileName(contentTypeHeader string) string {
	parts := strings.Split(contentTypeHeader, "name=")
	if len(parts) < 2 {
		return "unknown"
	}
	return strings.ReplaceAll(parts[1], "\"", "")
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


func getListName(path string) string {
	listName := strings.ReplaceAll(path, config.MailDirPath() + ".", "")
	listName = strings.Split(listName, "/")[0]
	return listName
}

func insertMessage(message models.Message) error {
	_, err := database.DBCon.Model(&message).
		Value("tsv_subject", "to_tsvector(?)", message.Subject).
		Value("tsv_body", "to_tsvector(?)", message.Body).
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
