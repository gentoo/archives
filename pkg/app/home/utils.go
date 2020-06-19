// miscellaneous utility functions used for the landing page of the application

package home

import (
	"archives/pkg/database"
	"archives/pkg/models"
	"github.com/go-pg/pg/v10"
	"html/template"
	"net/http"
	"strconv"
)

// renderIndexTemplate renders all templates used for the landing page
func renderIndexTemplate(w http.ResponseWriter, templateData interface{}) {
	templates := template.Must(
		template.Must(
			template.New("Show").
				Funcs(template.FuncMap{
					"makeMessage" : func(headers map[string][]string) models.Message {
						return models.Message{
							Headers:     headers,
						}
					},
				}).
				ParseGlob("web/templates/layout/*.tmpl")).
			ParseGlob("web/templates/home/*.tmpl"))

	templates.ExecuteTemplate(w, "home.tmpl", templateData)
}

// utility methods

func getAllMessagesCount() int {
	var messsageCount int
	database.DBCon.Model((*models.Message)(nil)).QueryOne(pg.Scan(&messsageCount), `
		SELECT
			count(DISTINCT messages.headers->>'Message-Id')
		FROM
			messages;
	`)
	return messsageCount
}

// formatMessageCount returns the formatted number of
// messages containing a thousands comma
func formatMessageCount(messageCount int) string {
	packages := strconv.Itoa(messageCount)
	if len(string(messageCount)) == 9 {
		return packages[:3] + "," + packages[3:6] + "," + packages[6:]
	} else if len(packages) == 8 {
		return packages[:2] + "," + packages[2:5] + "," + packages[5:]
	} else if len(packages) == 7 {
		return packages[:1] + "," + packages[1:4] + "," + packages[4:]
	} else if len(packages) == 6 {
		return packages[:3] + "," + packages[3:]
	} else if len(packages) == 5 {
		return packages[:2] + "," + packages[2:]
	} else if len(packages) == 4 {
		return packages[:1] + "," + packages[1:]
	} else {
		return packages
	}
}
