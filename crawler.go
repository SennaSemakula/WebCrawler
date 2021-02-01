package main

import (
	"os"
	"flag"
	"fmt"
	"strings"
	"net/http"
	"io"
	"golang.org/x/net/html"
	"log"
)

type Crawler struct {
	url string
	results []string
	depth int
	cache map[string][]string
}

var (
	client *http.Client
	infoLog  = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime|log.Lshortfile)
	errLog  = log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
)

func init() {
	client = &http.Client{}
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

					if uri := strings.HasPrefix(u.Val, "/"); !uri {
						continue
					}
					if _, ok := c.cache[c.url + u.Val]; ok {
						continue
					}
					uris = append(uris, u.Val)
				}
			}
		}
	}

	return uris, nil
}

func(c *Crawler) Crawl(url string, depth int) error {
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

	c.cache[url] = urls
	c.depth += 1

	if c.depth >= len(urls) {
		return nil
	}

	for range urls {
		// Have a concurrent function here
		if c.depth >= len(urls) {
			return nil
		}
		c.Crawl(c.url + urls[c.depth], c.depth)

	}

	return nil
}


func main() {
	url := flag.String("url", "https://monzo.com", "URL to start crawling")
	flag.Parse()

	results := make([]string, 0)
	cache := make(map[string][]string, 0)

	c := &Crawler{*url, results, 0, cache}

	if c.depth != 0 {
		log.Fatal("crawler must have initialised depth of 0")
	}

	c.Crawl(*url, 0)

}