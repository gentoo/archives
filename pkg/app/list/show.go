package list

import (
	"archives/pkg/database"
	"archives/pkg/models"
	"github.com/go-pg/pg/v10/orm"
	"net/http"
	"strings"
)

func Show(w http.ResponseWriter, r *http.Request) {

	listName := strings.ReplaceAll(r.URL.Path, "/", "")

	var res []struct {
		CombinedDate string
		MessageCount int
	}
	err := database.DBCon.Model((*models.Message)(nil)).
		WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q = q.WhereOr(`subject LIKE '[` + listName + `]%'`).
				WhereOr(`subject LIKE 'Re: [` + listName + `]%'`)
			return q, nil
		}).
		//WhereGroup(func(q *orm.Query) (*orm.Query, error) {
		//	q = q.WhereOr(`to LIKE '%` + listName + `@lists.gentoo.org%'`).
		//		WhereOr(`cc LIKE '%` + listName + `@lists.gentoo.org%'`).
		//		WhereOr(`to LIKE '%` + listName + `@gentoo.org%'`).
		//		WhereOr(`cc LIKE '%` + listName + `@gentoo.org%'`)
		//	return q, nil
		//}).
		ColumnExpr("to_char(date, 'YYYY-MM') AS combined_date").
		ColumnExpr("count(*) AS message_count").
		Group("combined_date").
		Order("combined_date DESC").
		Select(&res)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	renderShowTemplate(w, listName, res)
}
