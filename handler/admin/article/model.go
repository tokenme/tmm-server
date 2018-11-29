package article



type Article struct {
	Fileid      int    `json:"fileid" form:"fileid" `
	Author      string `json:"author" form:"author" `
	Title       string `json:"title" form:"title" `
	Link        string `json:"link" form:"link" `
	SourceUrl   string `json:"source_url" form:"source_url"`
	Cover       string `json:"cover" form:"cover" `
	PublishedOn string `json:"published_on" form:"published_on" `
	Digest      string `json:"digest" form:"digest"`
	Content     string `json:"content" form:"content"`
	Sortid      int    `json:"sortid" form:"sortid" `
	Published   int    `json:"published" form:"published"`
	Online      int    `json:"online" form:"online"`
	Id          int    `json:"id"  form:"id"`
}
