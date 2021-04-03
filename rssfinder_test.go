package rssfinder_test

import (
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/y-yagi/rssfinder"
)

func TestFind(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body := `
<html>
  <head>
  	<meta charset="UTF-8">
  	<link rel="profile" href="https://gmpg.org/xfn/11" />
  	<link rel="alternate" type="application/rss+xml" title="feed" href="http://localhost/feed/" />
  	<link rel="stylesheet" type="text/css" media="all" />
  </head>
	<body />
</html>
	`
		_, err := w.Write([]byte(body))
		if err != nil {
			panic(err)
		}
	}

	testserver := httptest.NewServer(http.HandlerFunc(handler))
	defer testserver.Close()

	result, err := rssfinder.Find(testserver.URL)
	if err != nil {
		t.Error(err)
	}

	if len(result) != 1 {
		t.Errorf("Expect rss count is 1, but got %v.", len(result))
	}

	if result[0].Type != "application/rss+xml" {
		t.Errorf("Expect type is 'application/rss+xml', but got %v.", result[0].Type)
	}

	if result[0].Title != "feed" {
		t.Errorf("Expect type is 'feed, but got %v.", result[0].Title)
	}

	if result[0].Href != "http://localhost/feed/" {
		t.Errorf("Expect type is 'http://localhost/feed/', but got %v.", result[0].Href)
	}
}

func TestFind_RelatvePath(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body := `
<html>
  <head>
		<meta charset="UTF-8">
		<link rel="alternate" type="application/rss+xml" title="RSS" href="rss">
  </head>
	<body />
</html>
	`
		_, err := w.Write([]byte(body))
		if err != nil {
			panic(err)
		}
	}

	testserver := httptest.NewServer(http.HandlerFunc(handler))
	defer testserver.Close()

	result, err := rssfinder.Find(testserver.URL)
	if err != nil {
		t.Error(err)
	}

	if len(result) != 1 {
		t.Errorf("Expect rss count is 1, but got %v.", len(result))
	}

	if result[0].Type != "application/rss+xml" {
		t.Errorf("Expect type is 'application/rss+xml', but got %v.", result[0].Type)
	}

	if result[0].Title != "RSS" {
		t.Errorf("Expect type is 'feed, but got %v.", result[0].Title)
	}

	expected := path.Join(testserver.URL, "rss")
	if result[0].Href != expected {
		t.Errorf("Expect href is '%v', but got %v.", expected, result[0].Href)
	}
}
