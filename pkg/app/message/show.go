// Used to show the landing page of the application

package message

import (
	"archives/pkg/database"
	"archives/pkg/models"
	"net/http"
	"sort"
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

	//
	// Compute Replies
	//
	var replies []*models.Message
	var queryParts []string
	for _, reference := range message.References {
		queryParts = append(queryParts, "reference_id = '" + reference.Id + "'")
	}
	queryParts = append(queryParts, "reference_id = '" + message.Id + "'")
	query := strings.Join(queryParts, " OR ")

	var refs []*models.MessageToReferences
	err = database.DBCon.Model(&refs).
		Where(query).
		Select()

	// part 2
	if len(refs) > 0 {
		var nextQueryParts []string
		for _, reference := range refs {
			nextQueryParts = append(nextQueryParts, "id = '" + reference.MessageId + "'")
		}
		nextQuery := strings.Join(nextQueryParts, " OR ")


		err = database.DBCon.Model(&replies).
			Where(nextQuery).
			Where("date >= '" + message.Date.Format("2006-01-02 15:04:05") + "'").
			Where("NOT id = ?", message.Id).
			Select()

		//
		// If In-Reply is null, but there are references we will use the last message
		// in the thread as In-Reply-To message
		//
		var inReplyTo []*models.Message
		if message.InReplyToId == "" || message.InReplyTo == nil {
			err = database.DBCon.Model(&inReplyTo).
				Where(nextQuery).
				Where("date <= '" + message.Date.Format("2006-01-02 15:04:05") + "'").
				Where("NOT id = ?", message.Id).
				Order("date DESC").
				Limit(1).
				Select()

			if err == nil && len(inReplyTo) > 0 {
				message.InReplyTo = inReplyTo[0]
			}
		}
	}

	sort.Slice(replies, func(i,j int) bool {
		return replies[i].Date.Before(replies[j].Date)
	})

	renderMessageTemplate(w, listName, message, replies)
}
