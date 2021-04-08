package rssfinder

import (
	"fmt"
	"net/http"
	"net/url"
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

type finder struct {
	url string
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
	f := finder{url: url}
	return f.find()
}

func (f *finder) find() ([]*Feed, error) {
	res, err := http.Get(f.url)
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
		err := fmt.Errorf("html parse error: %w", err)
		return nil, err
	}

	var feeds []*Feed
	f.findFeeds(node, &feeds)

	return feeds, nil
}

func (f *finder) buildFeed(node *html.Node) *Feed {
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
				u, _ := url.Parse(f.url)
				u.Path = path.Join(u.Path, v.Val)
				feed.Href = u.String()
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

func (f *finder) findFeeds(node *html.Node, feeds *[]*Feed) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if c.DataAtom == atom.Link {
				rss := f.buildFeed(c)
				if rss != nil {
					*feeds = append(*feeds, rss)
				}
			}
			f.findFeeds(c, feeds)
		}
	}
}
