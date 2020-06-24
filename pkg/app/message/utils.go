// miscellaneous utility functions used for the landing page of the application

package message

import (
	"archives/pkg/models"
	"html/template"
	"net/http"
	"strings"
)

// renderIndexTemplate renders all templates used for the landing page
func renderMessageTemplate(w http.ResponseWriter, listName string, message *models.Message, replies []*models.Message) {
	templates := template.Must(
		template.Must(
			template.New("Show").
				Funcs(getFuncMap()).
				ParseGlob("web/templates/layout/*.tmpl")).
			ParseGlob("web/templates/message/show.tmpl"))

	templateData := struct {
		ListName  string
		Message   *models.Message
		Replies   []*models.Message
	}{
		ListName:  listName,
		Message:   message,
		Replies:   replies,
	}

	templates.ExecuteTemplate(w, "show.tmpl", templateData)
}

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"formatAddr": formatAddr,
		"formatAddrList": formatAddrList,
	}
}

func formatAddr(addr string) string {
	if strings.Contains(addr, "@lists.gentoo.org") || strings.Contains(addr, "@gentoo.org") {
		addr = strings.ReplaceAll(addr, "@lists.gentoo.org", "@l.g.o")
		addr = strings.ReplaceAll(addr, "@gentoo.org", "@g.o")
	} else {
		start := false
		for i := len(addr) - 1; i > 0; i-- {
			if addr[i] == '@' {
				break
			}
			if start {
				out := []rune(addr)
				out[i] = '×'
				addr = string(out)
			}
			if addr[i] == '.' {
				start = true
			}
		}
	}
	return addr
}

func formatAddrList(addrList []string) string {
	var formatedAddrList []string
	for _, addr := range addrList {
		formatedAddrList = append(formatedAddrList, formatAddr(addr))
	}
	return strings.Join(formatedAddrList, ", ")
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}
