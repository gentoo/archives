package popular

import (
	"archives/pkg/database"
	"archives/pkg/models"
	"html/template"
	"net/http"
	"strconv"
)

// renderIndexTemplate renders all templates used for the landing page
func renderPopularThreads(w http.ResponseWriter, templateData interface{}) {
	templates := template.Must(
		template.Must(
			template.New("Popular").
				Funcs(template.FuncMap{
					"makeMessage" : func(headers map[string][]string) models.Message {
						return models.Message{
							Headers:     headers,
						}
					},
			}).
				ParseGlob("web/templates/layout/*.tmpl")).
			ParseGlob("web/templates/popular/*.tmpl"))

	templates.ExecuteTemplate(w, "threads.tmpl", templateData)
}

// utility methods

func GetPopularThreads(n int, date string) (models.Threads, error) {
	var popularThreads models.Threads
	err := database.DBCon.Model(&popularThreads).
		TableExpr(`(SELECT id, headers, regexp_replace(regexp_replace(regexp_replace(regexp_replace(headers::jsonb->>'Subject','^\["',''),'"\]$',''),'^Re:\s',''), '^\[.*\]', '') AS c FROM messages WHERE date >= '2020-06-12'::date) t`).
		ColumnExpr(`c as Subject, jsonb_agg(id)->>0 as Id, jsonb_agg(headers)->>0 as Headers, Count(*) as Count`).
		GroupExpr(`c`).
		OrderExpr(`count DESC`).
		Limit(n).
		Select()

	return popularThreads, err
}

func GetMessagesFromPopularThreads(threads models.Threads) []*models.Message {
	var popularThreads []*models.Message
	for _, thread := range threads {
		var messages []*models.Message
		err := database.DBCon.Model(&messages).
			Where(`headers::jsonb->>'Subject' LIKE '%` + thread.Id + `%'`).
			Select()
		if err == nil && len(messages) > 0 {
			messages[0].Comment = strconv.Itoa(thread.Count)
			popularThreads = append(popularThreads, messages[0])
		}
	}
	return popularThreads
}
