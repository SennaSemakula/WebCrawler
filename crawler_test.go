package main

import (
	"testing"
	"strings"
)


func TestFetch(t *testing.T) {
	t.Run("Fetch", func(t *testing.T) {
		url := "https://monzo.com"
		_, err := Fetch(url)
	
		if err != nil {
			t.Errorf("unable to fetch %q", url)
		}
	})
	
	t.Run("EmptyFetch", func(t *testing.T) {
		_, err := Fetch("")

		expected := "url cannot be empty."
	
		if err.Error() != expected {
			t.Errorf("expected %q but got %q", err.Error(), expected)
		}
	})
	
	t.Run("FetchHealth", func(t *testing.T) {
		url := "https://monzo.com/lolol"
		_, err := Fetch(url)
	
		expected := "endpoint is unhealthy!"
	
		if err.Error() != expected {
			t.Errorf("expected %q but got %q", err.Error(), expected)
		}
    })
}


func TestCrawl(t *testing.T) {
	cache := make(map[string][]string)
	c := Crawler{"https://monzo.com", []string{}, 0, cache}

	t.Run("CrawlEmpty", func(t *testing.T){
		err := c.Crawl("", 0)
		expected := "url cannot be empty."
	
		if err.Error() != expected {
			t.Errorf("expected %q but got %q", expected, err.Error())
		}
	})

	t.Run("CrawlDepth", func(t *testing.T){
		err := c.Crawl("https://monzo.com", -2)
		expected := "depth cannot be less than 0"
	
		if err.Error() != expected {
			t.Errorf("expected %q but got %q", expected, err.Error())
		}
	})

}

func TestHTMLParse(t *testing.T) {

	cache := make(map[string][]string)
	c := Crawler{"https://monzo.com", []string{}, 0, cache}

	t.Run("ValidHTML", func(t *testing.T) {
		html := `
		<html>
		<head>
		
		</head>
	
		<body>
		<div>
		<a href="/i/business" class="main-navigation__links__link">Business</a>
		<a href="/i/travel" class="main-navigation__links__link">Business</a>
		<a href="/i/expenses" class="main-navigation__links__link">Business</a>
		</div>
		</body>	
		</html>
		`

		r := strings.NewReader(html)
		expected := 3
		got, _ := c.ParseHTML(r)
	
		if expected != len(got) {
			t.Errorf("expected %v but got %v", expected, got)
		}
	})

	t.Run("EmptyLinks", func(t *testing.T) {
		html := `
		<html>
		<head>
		
		</head>
	
		<body>
		<div>
		</div>
		</body>	
		</html>
		`

		r := strings.NewReader(html)
		expected := 0
		got, _ := c.ParseHTML(r)
	
		if expected != len(got) {
			t.Errorf("expected %v but got %v", expected, got)
		}
	})

}
