package popular

import (
	"net/http"
)

func Threads(w http.ResponseWriter, r *http.Request) {
	threads, err := GetPopularThreads(25, "2020-06-01")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	renderPopularThreads(w, threads)
}
