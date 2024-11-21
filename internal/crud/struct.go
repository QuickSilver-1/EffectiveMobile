package crud

type Song struct {
	P       int
	Id      int
	Name    string
	Author  string
	Text    string
	Release string
	Link    string
}

type QueryList struct {
	P      int
	Filter string
}