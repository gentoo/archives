package importer

import (
	"archives/pkg/config"
	"archives/pkg/database"
	"archives/pkg/models"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var overAllcounter int
var importedCounter int
var startTime time.Time


func FullImport() {

	fmt.Println("Init import...")
	filepath.Walk(config.MailDirPath(), initImport)

	overAllcounter = 0
	importedCounter = 0
	startTime = time.Now()
	filepath.Walk(config.MailDirPath(), func(path string, info os.FileInfo, err error) error {
		if overAllcounter % 1000 == 0 {
			fmt.Println(strconv.Itoa(overAllcounter) + ": " + time.Now().Sub(startTime).String())
		}
		overAllcounter++
		if err != nil {
			return err
		}
		if !info.IsDir() && getDepth(path, config.MailDirPath()) >= 1 && isPublicList(path) {
			importedCounter++
			importMail(path, info.Name())
		}
		return nil
	})

	fmt.Println("Finished full import. Imported " + strconv.Itoa(importedCounter) + " messages.")
}

func IncrementalImport() {
	var messages []*models.Message
	err := database.DBCon.Model(&messages).
		Column("filename").
		Select()

	if err != nil {
		fmt.Println("Problem during import, aborting:")
		fmt.Println(err)
		return
	}

	fmt.Println("Init import...")
	filepath.Walk(config.MailDirPath(), initImport)

	overAllcounter = 0
	importedCounter = 0
	startTime = time.Now()
	filepath.Walk(config.MailDirPath(), func(path string, info os.FileInfo, err error) error {
		if overAllcounter % 1000 == 0 {
			fmt.Println(strconv.Itoa(overAllcounter) + ": " + time.Now().Sub(startTime).String())
		}
		overAllcounter++
		if err != nil {
			return err
		}
		if !info.IsDir() && getDepth(path, config.MailDirPath()) >= 1 && isPublicList(path) && !fileIsAlreadyPresent(path, messages) {
			importedCounter++
			importMail(path, info.Name())
		}
		return nil
	})

	fmt.Println("Finished incremental import. Imported " + strconv.Itoa(importedCounter) + " new messages.")
}

func RecomputeThreads() {

	fmt.Println("Init thread computation...")
	filepath.Walk(config.MailDirPath(), initImport)

	fmt.Println("Import thread references into database...")
	overallLength := len(mails)
	counter := 0
	for _, mail := range mails {
		if counter % 1000 == 0 {
			fmt.Println("Imported thread references for " + strconv.Itoa(counter) +  " of " + strconv.Itoa(overallLength) + " mails.")
		}
		counter++
		insertReferencesToMail(mail.RawReferences, mail.Id, mail.From)
	}

	fmt.Println("Finished thread computation.")
}

// utility methods

func fileIsAlreadyPresent(path string, messages []*models.Message) bool {
	for _, message := range messages {
		if strings.Contains(strings.TrimRight(path, ",S"), strings.TrimRight(message.Filename, ",S")){
			return true
		}
	}
	return false
}
