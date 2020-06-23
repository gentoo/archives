// miscellaneous utility functions used for the landing page of the application

package list

import (
	"archives/pkg/models"
	"archives/pkg/utils"
	"html/template"
	"net/http"
	"strconv"
)

type ListData struct {
	ListName    string
	Date        string
	CurrentPage int
	MaxPages    int
	Messages    []*models.Message
}

// renderIndexTemplate renders all templates used for the landing page
func renderShowTemplate(w http.ResponseWriter, listName string, messageData interface{}) {
	templates := template.Must(
		template.Must(
			template.New("Show").
				ParseGlob("web/templates/layout/*.tmpl")).
			ParseGlob("web/templates/list/*.tmpl"))

	templteData := struct {
		ListName    string
		MessageData interface{}
	}{
		ListName:    listName,
		MessageData: messageData,
	}

	templates.ExecuteTemplate(w, "show.tmpl", templteData)
}

// renderIndexTemplate renders all templates used for the landing page
func renderMessagesTemplate(w http.ResponseWriter, listName string, date string, currentPage int, maxPages int, messages []*models.Message) {
	templates := template.Must(
		template.Must(
			template.Must(
				template.New("Show").
					Funcs(getFuncMap()).
					ParseGlob("web/templates/layout/*.tmpl")).
				ParseGlob("web/templates/list/components/*.tmpl")).
			ParseGlob("web/templates/list/*.tmpl"))

	templates.ExecuteTemplate(w, "messages.tmpl", buildListData(listName, date, currentPage, maxPages, messages))
}

// renderIndexTemplate renders all templates used for the landing page
func renderThreadsTemplate(w http.ResponseWriter, listName string, date string, currentPage int, maxPages int, messages []*models.Message) {
	templates := template.Must(
		template.Must(
			template.Must(
				template.New("Show").
					Funcs(getFuncMap()).
					ParseGlob("web/templates/layout/*.tmpl")).
				ParseGlob("web/templates/list/components/*.tmpl")).
			ParseGlob("web/templates/list/*.tmpl"))

	templates.ExecuteTemplate(w, "threads.tmpl", buildListData(listName, date, currentPage, maxPages, messages))
}

// renderIndexTemplate renders all templates used for the landing page
func renderBrowseTemplate(w http.ResponseWriter, lists interface{}) {
	templates := template.Must(
		template.Must(
			template.New("Show").
				Funcs(template.FuncMap{
					"toString" : strconv.Itoa,
					"formatCount" : utils.FormatMessageCount,
			}).
				ParseGlob("web/templates/layout/*.tmpl")).
			ParseGlob("web/templates/list/*.tmpl"))

	templates.ExecuteTemplate(w, "browse.tmpl", lists)
}

// utility methods

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"min": func(a, b int) int {
			if a < b {
				return a
			}
			return b
		},
		"max": func(a, b int) int {
			if a < b {
				return b
			}
			return a
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"makeRange": makeRange,
	}
}

func buildListData(listName string, date string, currentPage int, maxPages int, messages []*models.Message) ListData {
	return ListData{
		ListName:    listName,
		Date:        date,
		CurrentPage: currentPage,
		MaxPages:    maxPages,
		Messages:    messages,
	}
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
