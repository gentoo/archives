package search

import (
	"archives/pkg/models"
	"html/template"
	"net/http"
)

type SearchData struct {
	SearchQuery        string
	ShowThreads        bool
	SearchResultsCount int
	CurrentPage        int
	MaxPages           int
	Messages           []*models.Message
}

// renderIndexTemplate renders all templates used for the landing page
func renderSearchTemplate(w http.ResponseWriter, showThreads bool, searchQuery string, messagesCount int, currentPage int, maxPages int, messages []*models.Message) {
	templates := template.Must(
		template.Must(
			template.Must(
				template.New("Show").
					Funcs(getFuncMap()).
					ParseGlob("web/templates/layout/*.tmpl")).
				ParseGlob("web/templates/search/components/pagination.tmpl")).
			ParseGlob("web/templates/search/*.tmpl"))

	templates.ExecuteTemplate(w, "searchresults.tmpl", buildSearchData(showThreads, searchQuery, messagesCount, currentPage, maxPages, messages))
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
		"mul": func(a, b int) int {
			return a * b
		},
		"makeRange": makeRange,
	}
}

func buildSearchData(showThreads bool, searchQuery string, messagesCount int, currentPage int, maxPages int, messages []*models.Message) SearchData {
	return SearchData{
		SearchQuery:        searchQuery,
		ShowThreads:        showThreads,
		SearchResultsCount: messagesCount,
		CurrentPage:        currentPage,
		MaxPages:           maxPages,
		Messages:           messages,
	}
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

// getParameterValue returns the value of a given parameter
func getParameterValue(parameterName string, r *http.Request) string {
	results, ok := r.URL.Query()[parameterName]
	if !ok {
		return ""
	}
	if len(results) == 0 {
		return ""
	}
	return results[0]
}
