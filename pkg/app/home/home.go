// Used to show the landing page of the application

package home

import (
	"archives/pkg/app/popular"
	"archives/pkg/cache"
	"archives/pkg/config"
	"archives/pkg/database"
	"archives/pkg/models"
	"github.com/go-pg/pg/v10/orm"
	"net/http"
	"time"
)

// Show renders a template to show the landing page of the application
func Show(w http.ResponseWriter, r *http.Request) {
	templateData := cache.Get("/")
	if templateData == nil {
		http.NotFound(w,r)
		return
	}
	renderIndexTemplate(w, templateData)
}

func ComputeTemplateData() interface{} {
	var mailingLists []models.MailingList

	for _, mailingList := range config.IndexMailingLists() {
		var messages []*models.Message
		database.DBCon.Model(&messages).
			WhereGroup(func(q *orm.Query) (*orm.Query, error) {
				q = q.WhereOr(`subject LIKE '[` + mailingList[0] + `]%'`).
					WhereOr(`subject LIKE 'Re: [` + mailingList[0] + `]%'`)
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
		return nil
	}
	if len(popularThreads) > 5 {
		popularThreads = popularThreads[:5]
	}

	return struct {
		MailingLists   []models.MailingList
		PopularThreads []*models.Message
		MessageCount   string
		CurrentMonth   string
	}{
		MailingLists:   mailingLists,
		PopularThreads: popularThreads,
		MessageCount:   formatMessageCount(getAllMessagesCount()),
		CurrentMonth:   time.Now().Format("2006-01"),
	}
}