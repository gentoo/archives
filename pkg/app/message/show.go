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
	err := database.DBCon.Model(message).
		Column("message.*").
		Relation("InReplyTo").
		Relation("References").
		WherePK().
		Select()

	if err != nil {
		http.NotFound(w, r)
		return
	}

	var queryParts []string
	for _, reference := range message.References {
		queryParts = append(queryParts, "reference_id = '" + reference.Id + "'")
	}
	query := strings.Join(queryParts, " OR ")

	var refs []*models.MessageToReferences
	err = database.DBCon.Model(&refs).
		Where(query).
		Select()

	// part 2
	// TODO only if len(refs) >= 1
	var nextQueryParts []string
	for _, reference := range refs {
		nextQueryParts = append(nextQueryParts, "id = '" + reference.MessageId + "'")
	}
	nextQuery := strings.Join(nextQueryParts, " OR ")

	var replies []*models.Message
	err = database.DBCon.Model(&replies).
		Where(nextQuery).
		// Where date is newer than message
		Select()

	renderMessageTemplate(w, listName, message, replies)
}
