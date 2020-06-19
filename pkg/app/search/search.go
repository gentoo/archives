package search

import (
	"archives/pkg/config"
	"archives/pkg/database"
	"archives/pkg/models"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func Search(w http.ResponseWriter, r *http.Request) {

	//
	// Parse search params
	//
	searchTerm := getParameterValue("q", r)
	showThreads := getParameterValue("threads", r) != ""
	page, err := strconv.Atoi(getParameterValue("page", r))
	var currentPage int
	var offset int

	if err != nil {
		currentPage = 1
		offset = 0
	} else {
		currentPage = page
		offset = 50 * (page - 1)
	}

	//
	// Step 1: Search for List with the same name and redirect
	//
	for _, list := range config.AllPublicMailingLists() {
		if strings.TrimSpace(searchTerm) == list {
			http.Redirect(w, r, "/"+list+"/", http.StatusMovedPermanently)
			return
		}
	}

	//
	// Step 2: Search by Author
	//
	var searchResults []*models.Message
	query := database.DBCon.Model(&searchResults).
		WhereOr(`headers::jsonb->>'From' LIKE ?`, "%"+searchTerm+"%").
		Order("date DESC")
	if showThreads {
		query = query.Where(`NOT headers::jsonb ? 'References'`).Where(`NOT headers::jsonb ? 'In-Reply-To'`)
	}

	messagesCount, _ := query.Count()
	err = query.Limit(50).Offset(offset).Select()

	if err == nil && messagesCount > 0 && strings.TrimSpace(searchTerm) != "gentoo" {
		maxPages := int(math.Ceil(float64(messagesCount) / float64(50)))
		renderSearchTemplate(w, showThreads, searchTerm, messagesCount, currentPage, maxPages, searchResults)
		return
	}

	//
	// Step 3: Search by Subject
	//
	query = database.DBCon.Model(&searchResults).
		Where(`tsv_subject @@ to_tsquery(''?'')`, searchTerm)
	if showThreads {
		query = query.Where(`NOT headers::jsonb ? 'References'`).Where(`NOT headers::jsonb ? 'In-Reply-To'`)
	}

	messagesCount, _ = query.Count()
	err = query.Limit(50).Offset(offset).Select()

	if err == nil && messagesCount > 0 {
		maxPages := int(math.Ceil(float64(messagesCount) / float64(50)))
		renderSearchTemplate(w, showThreads, searchTerm, messagesCount, currentPage, maxPages, searchResults)
		return
	}

	//
	// Step 4: Search by Message Body
	//
	query = database.DBCon.Model(&searchResults).
		Where(`tsv_body @@ to_tsquery(''?'')`, searchTerm)
	if showThreads {
		query = query.Where(`NOT headers::jsonb ? 'References'`).Where(`NOT headers::jsonb ? 'In-Reply-To'`)
	}

	messagesCount, _ = query.Count()
	err = query.Limit(50).Offset(offset).Select()

	if err != nil {
		http.NotFound(w, r)
		return
	}
	maxPages := int(math.Ceil(float64(messagesCount) / float64(50)))
	renderSearchTemplate(w, showThreads, searchTerm, messagesCount, currentPage, maxPages, searchResults)
}
