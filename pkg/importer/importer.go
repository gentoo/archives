package importer

import (
	"archives/pkg/config"
	"fmt"
	"path/filepath"
)

func FullImport() {
	fmt.Println("Init import...")
	filepath.Walk(config.MailDirPath(), initImport)
	fmt.Println("Start import...")
	filepath.Walk(config.MailDirPath(), importMail)
	fmt.Println("Finished import.")
}
