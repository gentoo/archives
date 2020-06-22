package main

import (
	"archives/pkg/app"
	"archives/pkg/database"
	"archives/pkg/importer"
	"flag"
	"fmt"
	"time"
)

var flagvar int

func main() {

	fmt.Println("Starting Gentoo Archives.")

	waitForPostgres()
	database.Connect()
	defer database.DBCon.Close()

	// main part

	fullImport := flag.Bool("full-import", false, "Start a full import, importing all mails")
	incrementalImport := flag.Bool("incremental-import", false, "Start a incremental import, importing only new mails")
	serve := flag.Bool("serve", false, "Start serving the web application")

	flag.Parse()

	if *fullImport {
		importer.FullImport()
	}

	if *incrementalImport {
		importer.IncrementalImport()
	}

	if *serve {
		app.Serve()
	}

	importer.WaitGroup.Wait()
}

// TODO this has to be solved differently
// wait for postgres to come up
func waitForPostgres() {
	time.Sleep(4 * time.Second)
}
