package models

type CrawlData struct {
	URL             string   `json:"url"`
	StatusCode      int      `json:"status_code"`
	ResponseTime    int      `json:"response_time"`
	ContentHash     string   `json:"content_hash"`
	Title           string   `json:"title"`
	H1              string   `json:"h1"`
	MetaDescription string   `json:"meta_description"`
	Links           []string `json:"links"`
}
