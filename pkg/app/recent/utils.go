// miscellaneous utility functions used for the landing page of the application

package recent

import (
	"html/template"
	"net/http"
)

// renderIndexTemplate renders all templates used for the landing page
func renderRecentTemplate(w http.ResponseWriter, templateData interface{}) {
	templates := template.Must(
		template.Must(
			template.New("Show").
				ParseGlob("web/templates/layout/*.tmpl")).
			ParseGlob("web/templates/recent/recent.tmpl"))

	templates.ExecuteTemplate(w, "recent.tmpl", templateData)
}
