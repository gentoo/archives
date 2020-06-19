package models

type Threads []struct {
	Id         string
	Headers    map[string][]string
	Subject    string
	Count int
}
