package cache

import (
	"archives/pkg/app/home"
	"archives/pkg/app/list"
	"archives/pkg/app/popular"
	"archives/pkg/app/recent"
	"archives/pkg/cache"
	"archives/pkg/config"
	"fmt"
	"net/http"
	"time"
)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	Update()
	w.Write([]byte("Updated."))
}

func Init(){
	cache.Init()
}

func Update(){
	fmt.Println("Updating caches...")

	startTime := time.Now()
	cache.Put("/", home.ComputeTemplateData())
	fmt.Println("> Updated '/' in " + time.Now().Sub(startTime).String())

	startTime = time.Now()
	cache.Put("/recent", recent.ComputeTemplateData())
	fmt.Println("> Updated '/recent' in " + time.Now().Sub(startTime).String())

	startTime = time.Now()
	cache.Put("/lists", list.ComputeBrowseTemplateData())
	fmt.Println("> Updated '/lists' in " + time.Now().Sub(startTime).String())

	startTime = time.Now()
	cache.Put("/popular", popular.ComputeThreadsTemplateData())
	fmt.Println("> Updated '/popular' in " + time.Now().Sub(startTime).String())

	startTime = time.Now()
	for _, listName := range config.AllPublicMailingLists() {
		tmpStartTime := time.Now()
		cache.Put("/"+listName+"/", list.ComputeShowTemplateData(listName))
		fmt.Println(">> Updated '/" + listName + "/' in " + time.Now().Sub(tmpStartTime).String())
	}
	fmt.Println("> Updated '/{{list}}/' in " + time.Now().Sub(startTime).String())
}
