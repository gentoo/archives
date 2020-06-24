// Used to show the landing page of the application

package home

import (
	"archives/pkg/app/popular"
	"archives/pkg/cache"
	"archives/pkg/models"
	"archives/pkg/utils"
	"net/http"
	"strconv"
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
		PopularThreads []*models.Message
		MessageCount   string
		CurrentMonth   string
	}{
		PopularThreads: popularThreads,
		MessageCount:   utils.FormatMessageCount(strconv.Itoa(getAllMessagesCount())),
		CurrentMonth:   time.Now().Format("2006-01"),
	}
}
