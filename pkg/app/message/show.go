// Used to show the landing page of the application

package message

import (
	"archives/pkg/database"
	"archives/pkg/models"
	"net/http"
	"strings"
)

// Show renders a template to show the landing page of the application
func Show(w http.ResponseWriter, r *http.Request) {

	urlParts := strings.Split(r.URL.Path, "/")
	listName := urlParts[1]
	messageHash := urlParts[len(urlParts)-1]

	message := &models.Message{Id: messageHash}
	err := database.DBCon.Select(message)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	var inReplyTos []*models.Message
	var inReplyTo *models.Message
	if message.InReplyTo != nil {
		err = database.DBCon.Model(&inReplyTos).
			Where(`(headers::jsonb->>'Message-Id')::jsonb ? '` + message.InReplyTo.Id + `'`).
			Select()
		if err != nil || len(inReplyTos) < 1 {
			inReplyTo = nil
		} else {
			inReplyTo = inReplyTos[0]
		}
	} else {
		inReplyTo = nil
	}

	var replies []*models.Message
	database.DBCon.Model(&replies).
		Where(`(headers::jsonb->>'References')::jsonb ? '` + message.Id + `'`).
		WhereOr(`(headers::jsonb->>'In-Reply-To')::jsonb ? '` + message.Id + `'`).
		Order("date ASC").Select()

	renderMessageTemplate(w, listName, message, inReplyTo, replies)
}
