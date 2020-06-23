package popular

import (
	"archives/pkg/database"
	"archives/pkg/models"
	"html/template"
	"net/http"
)

// renderIndexTemplate renders all templates used for the landing page
func renderPopularThreads(w http.ResponseWriter, templateData interface{}) {
	templates := template.Must(
		template.Must(
			template.New("Popular").
				ParseGlob("web/templates/layout/*.tmpl")).
			ParseGlob("web/templates/popular/threads.tmpl"))

	templates.ExecuteTemplate(w, "threads.tmpl", templateData)
}

// utility methods

func GetPopularThreads(n int, date string) ([]*models.Message, error) {

	var recentMessages []*models.Message

	err := database.DBCon.Model(&recentMessages).
		Where(`starts_thread = TRUE`).
		Where("not date is null").
		OrderExpr("date DESC").
		Limit(n).
		Select()

	return recentMessages, err
}
