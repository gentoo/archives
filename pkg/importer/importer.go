package importer

import (
	"archives/pkg/config"
	"log"
	"os"
	"path/filepath"
)

func FullImport() {
	err := filepath.Walk(config.MailDirPath(),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && getDepth(path, config.MailDirPath()) >= 1 {
				if isPublicList(path) {
					importMail(info.Name(), path, config.MailDirPath())
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}
