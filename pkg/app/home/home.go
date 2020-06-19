// Used to show the landing page of the application

package home

import (
	"archives/pkg/app/popular"
	"archives/pkg/config"
	"archives/pkg/database"
	"archives/pkg/models"
	"github.com/go-pg/pg/v10/orm"
	"net/http"
	"time"
)

// Show renders a template to show the landing page of the application
func Show(w http.ResponseWriter, r *http.Request) {

	var mailingLists []models.MailingList

	for _, mailingList := range config.IndexMailingLists() {
		var messages []*models.Message
		database.DBCon.Model(&messages).
			WhereGroup(func(q *orm.Query) (*orm.Query, error) {
				q = q.WhereOr(`(headers::jsonb->>'Subject')::jsonb->>0 LIKE '[` + mailingList[0] + `]%'`).
					WhereOr(`(headers::jsonb->>'Subject')::jsonb->>0 LIKE 'Re: [` + mailingList[0] + `]%'`)
				return q, nil
			}).
			WhereGroup(func(q *orm.Query) (*orm.Query, error) {
				q = q.WhereOr(`headers::jsonb->>'To' LIKE '%` + mailingList[0] + `@lists.gentoo.org%'`).
					WhereOr(`headers::jsonb->>'Cc' LIKE '%` + mailingList[0] + `@lists.gentoo.org%'`).
					WhereOr(`headers::jsonb->>'To' LIKE '%` + mailingList[0] + `@gentoo.org%'`).
					WhereOr(`headers::jsonb->>'Cc' LIKE '%` + mailingList[0] + `@gentoo.org%'`)
				return q, nil
			}).
			Order("date DESC").
			Limit(5).
			Select()

		mailingLists = append(mailingLists, models.MailingList{
			Name:        mailingList[0],
			Description: mailingList[1],
			Messages:    messages,
		})
	}

	//
	// Get popular threads
	//
	popularThreads, err := popular.GetPopularThreads(10, "2020-06-01")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if len(popularThreads) > 5 {
		popularThreads = popularThreads[:5]
	}

	templateData := struct {
		MailingLists   []models.MailingList
		PopularThreads models.Threads
		MessageCount   string
		CurrentMonth   string
	}{
		MailingLists:   mailingLists,
		PopularThreads: popularThreads,
		MessageCount:   formatMessageCount(getAllMessagesCount()),
		CurrentMonth:   time.Now().Format("2006-01"),
	}

	renderIndexTemplate(w, templateData)
}
