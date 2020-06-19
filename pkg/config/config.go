package config

import (
	"os"
	"strings"
)

func MailDirPath() string {
	mailDir := getEnv("ARCHIVES_MAILDIR_PATH", "/var/archives/.maildir/")
	if !strings.HasSuffix(mailDir, "/") {
		mailDir = mailDir + "/"
	}
	return mailDir
}

func Port() string {
	return getEnv("ARCHIVES_PORT", "5000")
}

func PostgresUser() string {
	return getEnv("ARCHIVES_POSTGRES_USER", "admin")
}

func PostgresPass() string {
	return getEnv("ARCHIVES_POSTGRES_PASS", "admin")
}

func PostgresDb() string {
	return getEnv("ARCHIVES_POSTGRES_DB", "garchives")
}

func PostgresHost() string {
	return getEnv("ARCHIVES_POSTGRES_HOST", "localhost")
}

func PostgresPort() string {
	return getEnv("ARCHIVES_POSTGRES_PORT", "5432")
}

func CacheControl() string {
	return getEnv("ARCHIVES_CACHE_CONTROL", "max-age=300")
}

func IndexMailingLists() [][]string {
	return [][]string{
		{"gentoo-dev", "is the main technical development mailing list of Gentoo"},
		{"gentoo-project", "contains non-technical discussion and propositions for the Gentoo Council"},
		{"gentoo-announce", "contains important news for all Gentoo stakeholders"},
		{"gentoo-user", "is our main support and Gentoo-related talk mailing list"},
		{"gentoo-commits", " - Lots of commits"},
		{"gentoo-dev-announce", "conveys important changes to all developers and interested users"}}
}

func AllPublicMailingLists() []string {
	var allMailingLists []string
	allMailingLists = append(allMailingLists, CurrentMailingLists()...)
	allMailingLists = append(allMailingLists, FrozenArchives()...)
	return allMailingLists
}

func CurrentMailingLists() []string {
	return []string{"gentoo-announce", "gentoo-commits", "gentoo-dev", "gentoo-dev-announce", "gentoo-nfp", "gentoo-project", "gentoo-user"}
}

func FrozenArchives() []string {
	return []string{"gentoo-arm", "gentoo-au", "gentoo-council", "gentoo-cygwin", "gentoo-desktop-research"}
}

func getEnv(key string, fallback string) string {
	if os.Getenv(key) != "" {
		return os.Getenv(key)
	} else {
		return fallback
	}
}
