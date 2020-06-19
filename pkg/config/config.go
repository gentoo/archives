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
	return []string{
		"gentoo-announce",
		"eudev",
		"gentoo-accessibility",
		"gentoo-admin",
		"gentoo-alpha",
		"gentoo-alt",
		"gentoo-amd64",
		"gentoo-announce",
		"gentoo-automated-testing",
		"gentoo-bsd",
		"gentoo-catalyst",
		"gentoo-cluster",
		"gentoo-commits",
		"gentoo-containers",
		"gentoo-desktop",
		"gentoo-dev",
		"gentoo-dev-announce",
		"gentoo-devhelp",
		"gentoo-doc",
		"gentoo-doc-cvs",
		"gentoo-embedded",
		"gentoo-foundation-announce",
		"gentoo-genkernel",
		"gentoo-github-gentoo",
		"gentoo-guis",
		"gentoo-hardened",
		"gentoo-hppa",
		"gentoo-java",
		"gentoo-kernel",
		"gentoo-keys",
		"gentoo-licenses",
		"gentoo-lisp",
		"gentoo-mips",
		"gentoo-mirrors",
		"gentoo-musl",
		"gentoo-nfp",
		"gentoo-openstack",
		"gentoo-perl",
		"gentoo-pms",
		"gentoo-portage-dev",
		"gentoo-powerpc",
		"gentoo-ppc-stable",
		"gentoo-pr",
		"gentoo-project",
		"gentoo-proxy-maint",
		"gentoo-python",
		"gentoo-qa",
		"gentoo-releng",
		"gentoo-releng-autobuilds",
		"gentoo-science",
		"gentoo-scm",
		"gentoo-security",
		"gentoo-server",
		"gentoo-soc",
		"gentoo-sparc",
		"gentoo-systemd",
		"gentoo-translators",
		"gentoo-user",
		"gentoo-user-br",
		"gentoo-user-cs",
		"gentoo-user-de",
		"gentoo-user-el",
		"gentoo-user-es",
		"gentoo-user-fr",
		"gentoo-user-hu",
		"gentoo-user-id",
		"gentoo-user-kr",
		"gentoo-user-pl",
		"gentoo-user-ru",
		"gentoo-user-tr"}
}

func FrozenArchives() []string {
	return []string{
		"gentoo-arm",
		"gentoo-arm",
		"gentoo-au",
		"gentoo-council",
		"gentoo-cygwin",
		"gentoo-desktop-research",
		"gentoo-dev-lang",
		"gentoo-devrel",
		"gentoo-doc-de",
		"gentoo-doc-el",
		"gentoo-doc-es",
		"gentoo-doc-fi",
		"gentoo-doc-fr",
		"gentoo-doc-hu",
		"gentoo-doc-id",
		"gentoo-doc-lt",
		"gentoo-doc-nl",
		"gentoo-doc-pl",
		"gentoo-doc-ru",
		"gentoo-docs-it",
		"gentoo-extreme-security",
		"gentoo-forum-translations",
		"gentoo-gnustep",
		"gentoo-gwn",
		"gentoo-gwn-de",
		"gentoo-gwn-es",
		"gentoo-gwn-fr",
		"gentoo-gwn-nl",
		"gentoo-gwn-pl",
		"gentoo-ia64",
		"gentoo-installer",
		"gentoo-kbase",
		"gentoo-laptop",
		"gentoo-media",
		"gentoo-nx",
		"gentoo-osx",
		"gentoo-performance",
		"gentoo-ppc-dev",
		"gentoo-ppc-user",
		"gentoo-proctors",
		"gentoo-scire",
		"gentoo-trustees",
		"gentoo-uk",
		"gentoo-vdr",
		"gentoo-web-user",
		"gentoo-xbox",
		"gnap-dev",
		"tenshi-announce",
		"tenshi-user",
		"www-redesign"}
}

func getEnv(key string, fallback string) string {
	if os.Getenv(key) != "" {
		return os.Getenv(key)
	} else {
		return fallback
	}
}
