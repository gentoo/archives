package models

import (
	"mime"
	"net/mail"
	"strings"
	"time"
)

type Message struct {
	Id       string `pg:",pk"`
	Filename string

	Headers     map[string][]string
	Body        map[string]string
	Attachments []Attachment

	Lists []string
	Date  time.Time

	//Search           types.ValueAppender // tsvector

	Comment string
	Hidden  bool

	//ParentId         string
	//Parent           Message -> pg fk?
}

type Header struct {
	Name    string
	Content string
}

type Body struct {
	ContentType string
	Content     string
}

type Attachment struct {
	Filename string
	Mime     string
	Content  string
}

func (m Message) GetSubject() string {
	return m.GetHeaderField("Subject")
}

func (m Message) GetListNameFromSubject() string {
	subject := m.GetSubject()
	listName := strings.Split(subject, "]")[0]
	listName = strings.ReplaceAll(listName, "[", "")
	listName = strings.ReplaceAll(listName, "Re:", "")
	listName = strings.TrimSpace(listName)
	return listName
}

func (m Message) GetAuthorName() string {
	addr, err := mail.ParseAddress(m.GetHeaderField("From"))
	if err != nil {
		return ""
	}
	return addr.Name
}

func (m Message) GetMessageId() string {
	messageId := m.GetHeaderField("Message-Id")
	messageId = strings.ReplaceAll(messageId, "<", "")
	messageId = strings.ReplaceAll(messageId, ">", "")
	messageId = strings.ReplaceAll(messageId, "\"", "")
	return messageId
}

func (m Message) GetInReplyTo() string {
	inReplyTo := m.GetHeaderField("In-Reply-To")
	inReplyTo = strings.ReplaceAll(inReplyTo, "<", "")
	inReplyTo = strings.ReplaceAll(inReplyTo, ">", "")
	inReplyTo = strings.ReplaceAll(inReplyTo, " ", "")
	return inReplyTo
}

func (m Message) GetHeaderField(key string) string {
	subject, found := m.Headers[key]
	if !found {
		return ""
	}
	header := strings.Join(subject, " ")
	if strings.Contains(header, "=?") {
		dec := new(mime.WordDecoder)
		decodedHeader, err := dec.DecodeHeader(header)
		if err != nil {
			return ""
		}
		return decodedHeader
	}
	return header
}

func (m Message) HasHeaderField(key string) bool {
	_, found := m.Headers[key]
	return found
}

func (m Message) GetBody() string {
	// Get text/plain body
	for contentType, content := range m.Body {
		if strings.Contains(contentType, "text/plain") {
			return content
		}
	}

	// If text/plain is not present, fall back to html
	for contentType, content := range m.Body {
		if strings.Contains(contentType, "text/html") {
			return content
		}
	}

	// If neither text/plain nor text/html is available return nothing
	return ""
}

func (m Message) HasAttachments() bool {
	for key, _ := range m.Body {
		if !(strings.Contains(key, "text/plain") || strings.Contains(key, "text/plain")) {
			return true
		}
	}
	return false
}

func (m Message) GetAttachments() []Attachment {
	var attachments []Attachment
	for key, content := range m.Body {
		if !(strings.Contains(key, "text/plain") || strings.Contains(key, "text/plain")) {
			attachments = append(attachments, Attachment{
				Filename: getAttachmentFileName(key),
				Mime:     strings.Split(key, ";")[0],
				Content:  content,
			})
		}
	}
	return attachments
}

// utility methods

func getAttachmentFileName(contentTypeHeader string) string {
	parts := strings.Split(contentTypeHeader, "name=")
	if len(parts) < 2 {
		return "unknown"
	}
	return strings.ReplaceAll(parts[1], "\"", "")
}
