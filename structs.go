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

type searchResultJSON struct {
	Result []struct {
		ID      int    `json:"id"`
		MediaID string `json:"media_id"`
		Title   struct {
			Pretty string `json:"pretty"`
		} `json:"title"`
		Images struct {
			Thumbnail struct {
				T string `json:"t"`
			} `json:"thumbnail"`
		} `json:"images"`
	} `json:"result"`
	NumPages int `json:"num_pages"`
	PerPage  int `json:"per_page"`
}

type doujinshiJSON struct {
	ID      int    `json:"id"`
	MediaID string `json:"media_id"`
	Title   struct {
		Pretty string `json:"pretty"`
	} `json:"title"`
	Images struct {
		Pages []struct {
			T string `json:"t"`
		} `json:"pages"`
	} `json:"images"`
	Tags []struct {
		Name string `json:"name"`
	} `json:"tags"`
	NumPages int `json:"num_pages"`
}

type tagJSON struct {
	Tags []struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	} `json:"tags"`
}
