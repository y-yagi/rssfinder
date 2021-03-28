package rssfinder

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type RSS struct {
	Type  string
	Href  string
	Title string
}

var rssTypes = map[string]struct{}{
	"application/rss+xml":  {},
	"application/atom+xml": {},
	"application/rdf+xml":  {},
	"application/rss":      {},
	"application/atom":     {},
	"application/rdf":      {},
	"text/rss+xml":         {},
	"text/atom+xml":        {},
	"text/rdf+xml":         {},
	"text/rss":             {},
	"text/atom":            {},
	"text/rdf":             {},
}

func Run(url string) ([]*RSS, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		err := fmt.Errorf("status code isn't success: %d %s", res.StatusCode, res.Status)
		return nil, err
	}

	node, err := html.Parse(res.Body)
	if err != nil {
		err := fmt.Errorf("html parse error: %v", err)
		return nil, err
	}

	var result []*RSS
	findRSSs(node, &result)

	return result, nil
}

func buildRSS(node *html.Node) *RSS {
	rss := &RSS{}
	for _, v := range node.Attr {
		if v.Key == "type" {
			if _, found := rssTypes[v.Val]; !found {
				return nil
			}
			rss.Type = v.Val
		}

		if v.Key == "href" {
			rss.Href = v.Val
		}

		if v.Key == "title" {
			rss.Title = v.Val
		}
	}

	if len(rss.Type) == 0 {
		return nil
	}

	return rss
}

func findRSSs(node *html.Node, result *[]*RSS) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if c.DataAtom == atom.Link {
				rss := buildRSS(c)
				if rss != nil {
					*result = append(*result, rss)
				}
			}
			findRSSs(c, result)
		}
	}
}
