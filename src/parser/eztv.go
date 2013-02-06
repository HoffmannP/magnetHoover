package parser

import (
	"code.google.com/p/go-html-transform/h5"
	"code.google.com/p/go-html-transform/html/transform"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

func sel(n *h5.Node, selector string) []*h5.Node {
	return transform.NewSelectorQuery(selector).Apply(n)
}

func eztvParser(uri string) (e [][]string, err error) {
	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// p := h5.NewParser(body)
	h := strings.Replace(string(body), "<center>", "", 1) // repair (remove) first unclosed <center> tag
	p := h5.NewParserFromString(h)

	if err = p.Parse(); err != nil {
		return nil, err
	}
	rows := sel(
		sel(p.Tree(), "table.forum_header_noborder")[0],
		"tr.forum_header_border")
	for _, row := range rows {
		e = append(e, eztvSnipp(row))
	}
	return
}

func eztvUrl(id string) (string, error) {
	res, err := http.Get("http://eztv.it/showlist/")
	if err != nil {
		return id, err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return id, err
	}
	p := h5.NewParserFromString(strings.Replace(string(b), "Pending</b>", "Pending", -1))
	if err = p.Parse(); err != nil {
		return id, err
	}
	rows := sel(p.Tree(), "tr[name=hover]")
	for _, row := range rows {
		a := sel(sel(row, "td.forum_thread_post")[0], "a")[0]
		if a.Children[0].Data() == id {
			for _, attr := range a.Attr {
				if attr.Name == "href" {
					return "http://eztv.it" + attr.Value, nil
				}
			}
			return id, errors.New("URI not found")
		}
	}
	return id, errors.New("Show not found")
}

func eztvSnipp(r *h5.Node) []string {
	c := transform.NewSelectorQuery("td.forum_thread_post").Apply(r)
	id := transform.NewSelectorQuery("a").Apply(c[1])[0].Children[0].Data()
	ls := transform.NewSelectorQuery("a").Apply(c[2])
	var u string
	for _, l := range ls {
		u = eztvLinkHref(l)
		if strings.Split(u, ":")[0] == "magnet" {
			break
		}
	}
	return []string{id, u}
}

func eztvLinkHref(a *h5.Node) string {
	for _, a := range a.Attr {
		if strings.ToLower(a.Name) == "href" {
			return a.Value
		}
	}
	return ""
}

func init() {
	register("EZTV", eztvParser, eztvUrl)
}
