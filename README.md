# Web Crawler
Simple web crawler that starts off with an initial url and prints all the links visited.

# Version
- Currently using go version 1.15.5

# Installation
1. Run ```go get``` to install the dependencies
2. Generate binary in local directory using ```go build```


# Example
- Start crawling with ```./crawler -url https://monzo.com``` or ```go run crawler.go``` (defaults to monzo domain)
- Visited URLs and their corresponding links should be printed to standard output
![Alt text](example.png?raw=true "Crawler results")


# Testing
1. Run ```go test```
2. Use ```go test -v``` for verbosity

# Contributors 
- Senna Semakula-Buuza