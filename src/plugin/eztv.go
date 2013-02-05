package plugin

import (
	"code.google.com/p/go-html-transform/h5"
	"code.google.com/p/go-html-transform/html/transform"
	"strings"
)

func sel(n *h5.Node, selector string) []*h5.Node {
	return transform.NewSelectorQuery(selector).Apply(n)
}

func EZTV(b string) (e [][]string, err error) {
	h := strings.Replace(b, "<center>", "", 1) // repair (remove) first unclosed <center> tag
	p := h5.NewParserFromString(string(h))
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
	register("EZTV", EZTV)
}
