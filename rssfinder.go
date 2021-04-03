package rssfinder

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Feed is a type of RSS feed.
type Feed struct {
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

// Find finds feeds from URL.
func Find(url string) ([]*Feed, error) {
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

	var feeds []*Feed
	findFeeds(node, &feeds, url)

	return feeds, nil
}

func buildFeed(node *html.Node, url string) *Feed {
	feed := &Feed{}
	for _, v := range node.Attr {
		if v.Key == "type" {
			if _, found := rssTypes[v.Val]; !found {
				return nil
			}
			feed.Type = v.Val
		}

		if v.Key == "href" {
			if strings.HasPrefix(v.Val, "http") {
				feed.Href = v.Val
			} else {
				feed.Href = path.Join(url, v.Val)
			}
		}

		if v.Key == "title" {
			feed.Title = v.Val
		}
	}

	if len(feed.Type) == 0 {
		return nil
	}

	return feed
}

func findFeeds(node *html.Node, feeds *[]*Feed, url string) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if c.DataAtom == atom.Link {
				rss := buildFeed(c, url)
				if rss != nil {
					*feeds = append(*feeds, rss)
				}
			}
			findFeeds(c, feeds, url)
		}
	}
}
