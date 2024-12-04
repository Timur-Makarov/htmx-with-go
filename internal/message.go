package internal

type Message struct {
	Images []string `json:"images"`
	UID    string   `json:"uid"`
	Text   string   `json:"text"`
}
