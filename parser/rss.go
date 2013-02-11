package parser

import (
	"encoding/xml"
	"net/http"
)

type enclosure struct {
	URL string `xml:"url,attr"`
}
type item struct {
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Link        string    `xml:"link"`
	Enclosure   enclosure `xml:"enclosure"`
}
type channel struct {
	Item []item `xml:"item"`
}
type rss struct {
	Channel channel `xml:"channel"`
}

func rssParser(uri string) (e [][]string, err error) {
	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	dec := xml.NewDecoder(res.Body)
	feed := new(rss)
	err = dec.Decode(feed)
	if err != nil {
		return nil, err
	}
	for _, i := range feed.Channel.Item {
		url := i.Link
		if url == "" {
			url = i.Enclosure.URL
		}
		id := i.Description
		if id == "" {
			id = i.Title
		}
		e = append(e, []string{id, url})
	}
	return
}

func init() {
	register("RSS", rssParser, nil)
}
