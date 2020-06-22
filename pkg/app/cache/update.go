package cache

import (
	"archives/pkg/app/home"
	"archives/pkg/app/list"
	"archives/pkg/app/popular"
	"archives/pkg/cache"
	"archives/pkg/config"
	"net/http"
)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	Update()
	w.Write([]byte("Updated."))
}

func Init(){
	cache.Init()
}

func Update(){
	cache.Put("/", home.ComputeTemplateData())
	cache.Put("/lists", list.ComputeBrowseTemplateData())
	cache.Put("/popular", popular.ComputeThreadsTemplateData())
	for _, listName := range config.AllPublicMailingLists() {
		cache.Put("/"+listName+"/", list.ComputeShowTemplateData(listName))
		cache.Put("/"+listName+"/", list.ComputeShowTemplateData(listName))
	}
}
