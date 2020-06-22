package popular

import (
	"archives/pkg/cache"
	"net/http"
)

func Threads(w http.ResponseWriter, r *http.Request) {
	templateData := cache.Get("/popular")
	if templateData == nil {
		http.NotFound(w,r)
		return
	}
	renderPopularThreads(w, templateData)
}

func ComputeThreadsTemplateData() interface{} {
	threads, err := GetPopularThreads(25, "2020-06-01")
	if err != nil {
		return nil
	}
	return threads
}