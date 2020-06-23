// miscellaneous utility functions used for the landing page of the application

package home

import (
	"archives/pkg/database"
	"archives/pkg/models"
	"github.com/go-pg/pg/v10"
	"html/template"
	"net/http"
)

// renderIndexTemplate renders all templates used for the landing page
func renderIndexTemplate(w http.ResponseWriter, templateData interface{}) {
	templates := template.Must(
		template.Must(
			template.New("Show").
				ParseGlob("web/templates/layout/*.tmpl")).
			ParseGlob("web/templates/home/home.tmpl"))

	templates.ExecuteTemplate(w, "home.tmpl", templateData)
}

// utility methods

func getAllMessagesCount() int {
	var messsageCount int
	database.DBCon.Model((*models.Message)(nil)).QueryOne(pg.Scan(&messsageCount), `
		SELECT
			count(DISTINCT messages.message_id)
		FROM
			messages;
	`)
	return messsageCount
}
