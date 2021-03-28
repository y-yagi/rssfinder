# rssfinder

[![Build Status](https://github.com/y-yagi/rssfinder/actions/workflows/ci.yml/badge.svg)](https://github.com/y-yagi/rssfinder/actions/workflows/ci.yml)

`rssfinder` is a library that finds RSS feeds from an URL.

Example:

```go
package main

import (
	"fmt"

	"github.com/y-yagi/rssfinder"
)

func main() {
	feeds, err := rssfinder.Find("https://github.com/y-yagi/rssfinder/tags")
	if err != nil {
		fmt.Printf("Find error: %v\n", err)
		return
	}

	for _, feed := range feeds {
		fmt.Printf("Type: '%v', Title: '%v', Href: '%v'\n", feed.Type, feed.Title, feed.Href)
		// Type: 'application/atom+xml', Title: 'Recent Commits to rssfinder:master', Href: 'https://github.com/y-yagi/rssfinder/commits/master.atom'
	}
}
```
