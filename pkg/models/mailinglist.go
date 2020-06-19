package models

type MailingList struct {
	Name         string
	Description  string
	Messages     []*Message
	MessageCount int
}
