package list

import (
	"archives/pkg/config"
	"archives/pkg/models"
	"net/http"
)

func Browse(w http.ResponseWriter, r *http.Request) {

	// Count number of messages in the current mailing lists
	var currentMailingLists []models.MailingList
	for _, listName := range config.CurrentMailingLists() {
		messageCount, _ := countMessages(listName)
		currentMailingLists = append(currentMailingLists, models.MailingList{
			Name:         listName,
			MessageCount: messageCount,
		})
	}

	// Count number of messages in the frozen archives
	var frozenArchives []models.MailingList
	for _, listName := range config.FrozenArchives() {
		messageCount, _ := countMessages(listName)
		frozenArchives = append(frozenArchives, models.MailingList{
			Name:         listName,
			MessageCount: messageCount,
		})
	}

	browseData := struct {
		CurrentMailingLists []models.MailingList
		FrozenArchives      []models.MailingList
	}{
		CurrentMailingLists: currentMailingLists,
		FrozenArchives:      frozenArchives,
	}

	renderBrowseTemplate(w, browseData)
}
