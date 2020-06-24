// Used to show the landing page of the application

package recent

import (
	"archives/pkg/cache"
	"archives/pkg/config"
	"archives/pkg/database"
	"archives/pkg/models"
	"net/http"
	"time"
)

// Show renders a template to show the landing page of the application
func Show(w http.ResponseWriter, r *http.Request) {
	templateData := cache.Get("/recent")
	if templateData == nil {
		http.NotFound(w,r)
		return
	}
	renderRecentTemplate(w, templateData)
}

func ComputeTemplateData() interface{} {
	var mailingLists []models.MailingList

	for _, mailingList := range config.IndexMailingLists() {
		var messages []*models.Message
		database.DBCon.Model(&messages).
			Where("list = ?", mailingList[0]).
			Where("not date is null").
			Order("date DESC").
			Limit(5).
			Select()

		mailingLists = append(mailingLists, models.MailingList{
			Name:        mailingList[0],
			Description: mailingList[1],
			Messages:    messages,
		})
	}

	return struct {
		MailingLists   []models.MailingList
		PopularThreads []*models.Message
		MessageCount   string
		CurrentMonth   string
	}{
		MailingLists:   mailingLists,
		CurrentMonth:   time.Now().Format("2006-01"),
	}
}
