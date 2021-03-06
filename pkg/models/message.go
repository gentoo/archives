package models

import (
	"mime"
	"net/mail"
	"strings"
	"time"
)

type Message struct {
	tableName struct{} `pg:",discard_unknown_columns"`

	Id       string `pg:",pk"`
	MessageIdField string
	Filename string

	List string

	From string
	To []string
	Cc []string

	Subject string
	Body string

	Date time.Time

	// fk
	InReplyTo *Message `pg:"fk:in_reply_to_id"` // fk specifies foreign key
	InReplyToId string

	// many to many
	//References []string
	RawReferences []string `pg:"-"`
	References []Message `pg:"many2many:message_to_references,joinFK:reference_id"`

	Attachments []Attachment

	StartsThread bool

	Comment string
	Hidden  bool

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

type MessageToReferences struct {
	MessageId string
	ReferenceId  string
}

func (m Message) MessageId() string {
	return m.MessageIdField
}

func (m Message) GetListNameFromSubject() string {
	subject := m.Subject
	listName := strings.Split(subject, "]")[0]
	listName = strings.ReplaceAll(listName, "[", "")
	listName = strings.ReplaceAll(listName, "Re:", "")
	listName = strings.TrimSpace(listName)
	return listName
}

func (m Message) GetAuthorName() string {
	addr, err := mail.ParseAddress(m.From)
	if err != nil {
		return ""
	}
	if addr.Name == "" {
		return addr.Address
	}
	return addr.Name
}

func (m Message) GetSubject() string {
	header, err := new(mime.WordDecoder).DecodeHeader(m.Subject)
	if err == nil {
		return header
	}
	return m.Subject
}

func (m Message) GetMessageId() string {
	messageId := m.MessageIdField
	messageId = strings.ReplaceAll(messageId, "<", "")
	messageId = strings.ReplaceAll(messageId, ">", "")
	messageId = strings.ReplaceAll(messageId, "\"", "")
	return messageId
}

func (m Message) GetInReplyTo() string {
	inReplyTo := m.InReplyTo.MessageIdField
	inReplyTo = strings.ReplaceAll(inReplyTo, "<", "")
	inReplyTo = strings.ReplaceAll(inReplyTo, ">", "")
	inReplyTo = strings.ReplaceAll(inReplyTo, " ", "")
	return inReplyTo
}


