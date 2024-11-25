package crud

type Song struct {
	P       int    `json:"p"`
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Author  string `json:"author"`
	Text    string `json:"text"`
	Release string `json:"release"`
	Link    string `json:"link"`
}

type QueryList struct {
	P      int    `json:"p"`
	Filter string `json:"filter"`
}