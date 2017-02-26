package main

type defType int

const (
	union defType = 1
	rec           = 2
	enum          = 3
)

type nameList []string

type def struct {
	name string
	pos  string
	t    defType
	u    nameList
	f    fieldList
	sql  string
}

type defList []def

type fieldList []field

type field struct {
	name    string
	t       string
	isArray bool
}
