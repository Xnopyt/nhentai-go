package nhentai

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
)

//Search - Search for a term on nHentai and return results
func Search(query string, page int) (*SearchResponse, error) {
	requrl := strings.Replace(query, " ", "+", -1)
	requrl = url.QueryEscape(requrl)
	requrl = "https://nhentai.net/search?q=" + requrl + "&page=" + strconv.Itoa(page)

	client := &http.Client{}
	req, err := http.NewRequest("GET", requrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Xnopyts_nHentai_Scraper/0.1")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(string(body))

	if doc.Error != nil {
		return nil, doc.Error
	}

	var parsedResults SearchResponse

	findmaxpage := doc.Find("a", "class", "last")

	if findmaxpage.Error != nil {
		return nil, errors.New("the search query returned no pages")
	}

	maxpage, _ := strconv.Atoi(strings.Split(findmaxpage.Attrs()["href"], "&page=")[1])

	parsedResults.MaxPage = maxpage

	results := doc.FindAll("div", "class", "gallery")

	for _, v := range results {
		var result SearchResult
		result.Title = v.Find("div", "class", "caption").Text()
		id := v.Find("a").Attrs()["href"]
		result.ID = strings.Split(id, "/")[2]
		result.URL = "https://nhentai.net/g/" + result.ID
		result.ThumbURL = v.Find("img").Attrs()["data-src"]
		parsedResults.Results = append(parsedResults.Results, result)
	}

	return &parsedResults, nil
}

//Get - get a doujinshi by nHentai id
func Get(id int) (*Doujinshi, error) {
	requrl := "https://nhentai.net/g/" + strconv.Itoa(id)

	client := &http.Client{}
	req, err := http.NewRequest("GET", requrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Xnopyts_nHentai_Scraper/0.1")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("could not find a result for the given id")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(string(body))

	if doc.Error != nil {
		return nil, doc.Error
	}

	info := doc.Find("div", "id", "info")

	if info.Error != nil {
		return nil, info.Error
	}

	var doujinshi Doujinshi

	doujinshi.Title = info.Find("h1").Text()
	doujinshi.ID = strconv.Itoa(id)
	doujinshi.URL = requrl

	tagContainer := info.FindAll("div", "class", "tag-container")
	var tags soup.Root

	for _, v := range tagContainer {
		if strings.Contains(v.Text(), "Tags:") {
			tags = v
		}
	}

	tagA := tags.FindAll("a")

	for _, v := range tagA {
		href := v.Attrs()["href"]
		tag := strings.Replace(strings.Split(href, "/")[2], "-", " ", -1)
		doujinshi.Tags = append(doujinshi.Tags, tag)
	}

	imgContainer := doc.Find("div", "id", "thumbnail-container")
	imgs := imgContainer.FindAll("div", "class", "thumb-container")

	for _, v := range imgs {
		var page Page
		urlsplit := strings.Split(v.Find("img").Attrs()["data-src"], "/")
		doujinshi.MediaID = urlsplit[4]
		x := strings.Split(urlsplit[5], "t.")
		doujinshi.TotalPages, _ = strconv.Atoi(x[0])
		page.Num, _ = strconv.Atoi(x[0])
		page.Ext = x[1]
		doujinshi.Pages = append(doujinshi.Pages, page)
	}

	return &doujinshi, nil
}

//Tag - Search for a tag on nHentai and return results
func Tag(query string, page int) (*SearchResponse, error) {
	query = strings.ToLower(query)
	requrl := strings.Replace(query, " ", "-", -1)
	requrl = url.QueryEscape(requrl)
	requrl = "https://nhentai.net/tag/" + requrl + "?page=" + strconv.Itoa(page)

	client := &http.Client{}
	req, err := http.NewRequest("GET", requrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Xnopyts_nHentai_Scraper/0.1")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(string(body))

	if doc.Error != nil {
		return nil, doc.Error
	}

	var parsedResults SearchResponse

	findmaxpage := doc.Find("a", "class", "last")

	if findmaxpage.Error != nil {
		return nil, errors.New("the search query returned no pages")
	}

	maxpage, _ := strconv.Atoi(strings.Split(findmaxpage.Attrs()["href"], "?page=")[1])

	parsedResults.MaxPage = maxpage

	results := doc.FindAll("div", "class", "gallery")

	for _, v := range results {
		var result SearchResult
		result.Title = v.Find("div", "class", "caption").Text()
		id := v.Find("a").Attrs()["href"]
		result.ID = strings.Split(id, "/")[2]
		result.URL = "https://nhentai.net/g/" + result.ID
		result.ThumbURL = v.Find("img").Attrs()["data-src"]
		parsedResults.Results = append(parsedResults.Results, result)
	}

	return &parsedResults, nil
}
