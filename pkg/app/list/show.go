package list

import (
	"archives/pkg/cache"
	"archives/pkg/database"
	"archives/pkg/models"
	"github.com/go-pg/pg/v10/orm"
	"net/http"
	"strings"
)

func Show(w http.ResponseWriter, r *http.Request) {
	listName := strings.ReplaceAll(r.URL.Path, "/", "")
	templateData := cache.Get("/" + listName + "/")
	if templateData == nil {
		http.NotFound(w,r)
		return
	}
	renderShowTemplate(w, listName, templateData)
}


func ComputeShowTemplateData(listName string) interface{} {
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
		ColumnExpr("to_char(date, 'YYYY-MM') AS combined_date").
		ColumnExpr("count(*) AS message_count").
		Group("combined_date").
		Order("combined_date DESC").
		Select(&res)

	if err != nil {
		return nil
	}
	return res
}