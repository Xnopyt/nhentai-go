package nhentai

//SearchResponse - The reponse from a nhentai search query
type SearchResponse struct {
	MaxPage int
	Results []SearchResult
}

//SearchResult - A single result from a search
type SearchResult struct {
	ID       string
	Title    string
	ThumbURL string
	URL      string
}

//Doujinshi - Information and content from a given id
type Doujinshi struct {
	ID         string
	Title      string
	TotalPages int
	Pages      []Page
	Tags       []string
	URL        string
	MediaID    string
}

//Page - A single page of a doujinshi
type Page struct {
	Ext string
	Num int
}
