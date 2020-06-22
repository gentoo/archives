package list

import (
	"archives/pkg/cache"
	"archives/pkg/config"
	"archives/pkg/database"
	"archives/pkg/models"
	"net/http"
)

func Browse(w http.ResponseWriter, r *http.Request) {
	templateData := cache.Get("/lists")
	if templateData == nil {
		http.NotFound(w,r)
		return
	}
	renderBrowseTemplate(w, templateData)
}

func ComputeBrowseTemplateData() interface{} {
	var res []struct {
		Name string
		MessageCount int
	}
	err := database.DBCon.Model((*models.Message)(nil)).
		ColumnExpr("list as name, count(*) as message_count").
		Group("list").
		Select(&res)

	if err != nil {
		return nil
	}

	var currentMailingLists []models.MailingList
	var frozenArchives []models.MailingList

	for _, list := range res {
		if contains(config.CurrentMailingLists(), list.Name) {
			currentMailingLists = append(currentMailingLists, models.MailingList{
				Name: list.Name,
				MessageCount: list.MessageCount,
			})
		}else if contains(config.FrozenArchives(), list.Name) {
			frozenArchives = append(frozenArchives, models.MailingList{
				Name: list.Name,
				MessageCount: list.MessageCount,
			})
		}
	}

	return struct {
		CurrentMailingLists []models.MailingList
		FrozenArchives      []models.MailingList
	}{
		CurrentMailingLists: currentMailingLists,
		FrozenArchives:      frozenArchives,
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
