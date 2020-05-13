//Package nhentai - A library for pulling doujinshi from nhentai using the JSON api.
//Provides the ability to search based on keyword, tag or id, with id providing full infomation about the doujinshi.
package nhentai

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var exts map[string]string

func init() {
	exts = make(map[string]string)
	exts["j"] = ".jpg"
	exts["p"] = ".png"
	exts["g"] = ".gif"
}

func httpGet(url string) (r []byte, code int, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("User-Agent", "Xnopyts_nHentai_Scraper/0.1")

	resp, err := client.Do(req)

	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}

//Search - Search for a term on nHentai and return results
func Search(query string, page int) (*SearchResponse, error) {
	requrl := strings.Replace(query, " ", "+", -1)
	requrl = url.QueryEscape(requrl)
	requrl = "https://nhentai.net/api/galleries/search?q=" + requrl + "&page=" + strconv.Itoa(page)

	body, _, err := httpGet(requrl)

	if err != nil {
		return nil, err
	}

	var result searchResultJSON

	json.Unmarshal(body, &result)

	if result.NumPages == 0 {
		return nil, errors.New("the search query returned no pages")
	}

	var parsedResults SearchResponse
	parsedResults.MaxPage = result.NumPages

	for _, v := range result.Result {
		var r SearchResult
		r.ID = strconv.Itoa(v.ID)
		r.ThumbURL = "https://t.nhentai.net/galleries/" + v.MediaID + "/cover" + exts[v.Images.Thumbnail.T]
		r.Title = v.Title.Pretty
		r.URL = "https://nhentai.net/g/" + r.ID
	}

	return &parsedResults, nil
}

//Get - get a doujinshi by nHentai id
func Get(id int) (*Doujinshi, error) {
	requrl := "https://nhentai.net/api/gallery/" + strconv.Itoa(id)

	body, status, err := httpGet(requrl)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, errors.New("could not find a result for the given id")
	}

	var result doujinshiJSON
	var doujinshi Doujinshi

	json.Unmarshal(body, &result)

	doujinshi.ID = strconv.Itoa(result.ID)
	doujinshi.MediaID = result.MediaID
	doujinshi.Title = result.Title.Pretty
	doujinshi.URL = "https://nhentai.net/g/" + doujinshi.ID
	doujinshi.TotalPages = len(result.Images.Pages)

	for _, v := range result.Tags {
		doujinshi.Tags = append(doujinshi.Tags, v.Name)
	}

	for i, v := range result.Images.Pages {
		var p Page
		p.Ext = exts[v.T]
		p.Num = i + 1
		doujinshi.Pages = append(doujinshi.Pages, p)
	}

	return &doujinshi, nil
}

//Tag - Search for a tag on nHentai and return results
func Tag(query string, page int) (*SearchResponse, error) {
	query = strings.ToLower(query)
	requrl := strings.Replace(query, " ", "-", -1)
	requrl = url.QueryEscape(requrl)
	requrl = "https://nhentai.net/tag/" + requrl

	body, _, err := httpGet(requrl)

	if err != nil {
		return nil, err
	}

	furl := regexp.MustCompile("href=\"\\/g\\/.*\\/\"")
	tagURL := string(furl.Find(body))

	tagURL = "https://nhentai.net/api/gallery/" + tagURL[9:len(tagURL)-2]
	body, _, err = httpGet(tagURL)

	if err != nil {
		return nil, err
	}

	var tagIDs tagJSON
	var tagID string

	json.Unmarshal(body, &tagIDs)

	for _, v := range tagIDs.Tags {
		v.Name = strings.ToLower(v.Name)
		v.Name = strings.Replace(v.Name, " ", "-", -1)
		if query == v.Name {
			tagID = strconv.Itoa(v.ID)
			break
		}
	}

	requrl = "https://nhentai.net/api/galleries/tagged?tag_id=" + tagID + "&page=" + strconv.Itoa(page)

	body, _, err = httpGet(requrl)

	if err != nil {
		return nil, err
	}

	var result searchResultJSON

	json.Unmarshal(body, &result)

	if result.NumPages == 0 {
		return nil, errors.New("the search query returned no pages")
	}

	var parsedResults SearchResponse
	parsedResults.MaxPage = result.NumPages

	for _, v := range result.Result {
		var r SearchResult
		r.ID = strconv.Itoa(v.ID)
		r.ThumbURL = "https://t.nhentai.net/galleries/" + v.MediaID + "/cover" + exts[v.Images.Thumbnail.T]
		r.Title = v.Title.Pretty
		r.URL = "https://nhentai.net/g/" + r.ID
	}

	return &parsedResults, nil
}
