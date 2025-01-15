package transport

type ArticleResponse struct {
	Status   string   `json:"status"`
	Data     Article  `json:"data"`
	Metadata Metadata `json:"metadata"`
}

type ArticlesResponse struct {
	Status   string    `json:"status"`
	Data     []Article `json:"data"`
	Metadata Metadata  `json:"metadata"`
}

type Article struct {
	Id              string `json:"id"`
	ClientId        string `json:"clientId"`
	ClientArticleId string `json:"clientArticleId"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	PublishDate     string `json:"publishDate"`
}

type Metadata struct {
	CreatedAt string `json:"createdAt"`
}
