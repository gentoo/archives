package main

import (
	"archives/pkg/app"
	"archives/pkg/database"
	"archives/pkg/importer"
	"flag"
	"fmt"
)

var flagvar int

func main() {

	fmt.Println("Starting Gentoo Archives.")

	waitForPostgres()
	database.Connect()
	defer database.DBCon.Close()

	// main part

	fullImport := flag.Bool("fullimport", false, "Start a full import, importing all mails")
	serve := flag.Bool("serve", false, "Start serving the web application")
	flag.Parse()

	if *fullImport {
		importer.FullImport()
	}

	if *serve {
		app.Serve()
	}

}

// TODO this has to be solved differently
// wait for postgres to come up
func waitForPostgres() {
	//time.Sleep(2 * time.Second)
}
