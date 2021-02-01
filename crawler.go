package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Crawler struct {
	url     string
	results []string
	depth   int
	cache   map[string][]string
}

var (
	client  *http.Client
	infoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime|log.Lshortfile)
	errLog  = log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
)

func init() {
	client = &http.Client{}
}

func (c *Crawler) New(url string) (*Crawler, error) {

	results := make([]string, 0)
	cache := make(map[string][]string, 0)

	c = &Crawler{url, results, 0, cache}

	r, err := Fetch(c.url)

	if err != nil {
		return nil, err
	}

	urls, err := c.ParseHTML(r)

	if err != nil {
		return nil, err
	}

	c.results = append(c.results, urls...)

	return c, nil

}

func Fetch(url string) (io.Reader, error) {
	if len(url) <= 0 {
		return nil, fmt.Errorf("url cannot be empty.")
	}

	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)

	infoLog.Printf("Visiting %q\n", url)

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("endpoint is unhealthy!")
	}

	if err != nil {
		resp.Body.Close()
		return nil, err
	}

	return resp.Body, nil
}

func (c *Crawler) ParseHTML(r io.Reader) ([]string, error) {
	uris := make([]string, 0)
	z := html.NewTokenizer(r)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the doc
			return uris, nil
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if isAnchor {
				for _, u := range t.Attr {
					// TODO:
					// Maybe there may be an absolute link with monzo so create a function to validate the domain that takes in a url as a tring
					if uri := strings.HasPrefix(u.Val, "/"); !uri {
						continue
					}
					if _, ok := c.cache[c.url+u.Val]; ok {
						continue
					}
					uris = append(uris, u.Val)
				}
			}
		}
	}

	return uris, nil
}

func (c *Crawler) Crawl(url string, depth int) error {
	urls := make([]string, 0)

	if len(url) <= 0 {
		return fmt.Errorf("url cannot be empty.")
	}

	if depth < 0 {
		return fmt.Errorf("depth cannot be less than 0")
	}

	if tags, ok := c.cache[url]; ok {
		infoLog.Printf("%q is already in cache\n", url)
		urls = append(urls, tags...)
	} else {
		body, err := Fetch(url)

		if err != nil {
			return err
		}

		fetchedUrls, err := c.ParseHTML(body)

		urls = append(urls, fetchedUrls...)

		if err != nil {
			return err
		}
	}
	infoLog.Printf("URLs found:\n")
	for _, u := range urls {
		infoLog.Println(c.url + u)
	}

	if c.depth >= len(c.results) {
		return nil
	}
	c.cache[url] = urls
	c.depth += 1

	for range urls {
		if c.depth >= len(urls)-1 {
			return nil
		}
		c.Crawl(c.url+c.results[c.depth], c.depth)
	}

	return nil
}

func main() {
	url := flag.String("url", "https://monzo.com", "URL to start crawling")
	flag.Parse()

	var crawl *Crawler
	c, err := crawl.New(*url)

	if err != nil {
		errLog.Println(err)
	}

	if c.depth != 0 {
		log.Fatal("crawler must have initialised depth of 0")
	}

	c.Crawl(*url, 0)

}
