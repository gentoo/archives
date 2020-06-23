package list

import (
	"archives/pkg/database"
	"archives/pkg/models"
	"github.com/go-pg/pg/v10/orm"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func Threads(w http.ResponseWriter, r *http.Request) {

	urlParts := strings.Split(r.URL.Path, "/threads/")
	if len(urlParts) != 2 {
		http.NotFound(w, r)
		return
	}

	listName := strings.ReplaceAll(urlParts[0], "/", "")
	trailingUrlParts := strings.Split(urlParts[1], "/")
	combinedDate := trailingUrlParts[0]
	currentPage := 1
	if len(trailingUrlParts) > 1 {
		parsedCurrentPage, err := strconv.Atoi(trailingUrlParts[1])
		if err == nil {
			currentPage = parsedCurrentPage
		}
	}
	offset := (currentPage - 1) * 50

	var messages []*models.Message
	query := database.DBCon.Model(&messages).
		Column("id", "subject", "from", "date").
		Where("to_char(date, 'YYYY-MM') = ?", combinedDate).
		Where(`starts_thread = TRUE`).
		WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q = q.WhereOr(`subject LIKE '[` + listName + `]%'`).
				WhereOr(`subject LIKE 'Re: [` + listName + `]%'`)
			return q, nil
		}).
		Order("date DESC")

	messagesCount, _ := query.Count()
	query.Limit(50).Offset(offset).Select()

	maxPages := int(math.Ceil(float64(messagesCount) / float64(50)))

	renderThreadsTemplate(w, listName, combinedDate, currentPage, maxPages, messages)

}
