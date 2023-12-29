package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/gocolly/colly"
	_ "go.uber.org/automaxprocs"
)

var (
	domain, header            string
	parallelism, delay, sleep int
	daemon                    bool
)

func init() {
	flag.StringVar(&domain, "domain", "", "Set url for crawling. Example: https://example.com")
	flag.IntVar(&parallelism, "parallelism", 2, "Parallelism is the number of the maximum allowed concurrent requests of the matching domains")
	flag.IntVar(&delay, "delay", 1, "Delay is the duration to wait before creating a new request to the matching domains")
	flag.BoolVar(&daemon, "daemon", false, "Run crawler on daemon mode")
	flag.IntVar(&sleep, "sleep", 60, "Time in seconds to wait before run crawler again")
	flag.StringVar(&header, "header", "", "Set header for crawler request. Example: header_name:header_value")
}

func crawler() {
	u, err := url.Parse(domain)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// Instantiate default collector
	c := colly.NewCollector(
		// Turn on asynchronous requests
		colly.Async(true),
		// Visit only domain
		colly.AllowedDomains(u.Host),
	)

	// Limit the number of threads
	c.Limit(&colly.LimitRule{
		DomainGlob:  u.Host,
		Parallelism: parallelism,
		Delay:       time.Duration(delay) * time.Second,
	})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	if len(header) > 0 {
		c.OnRequest(func(r *colly.Request) {
			reg := regexp.MustCompile(`(.*):(.*)`)
			headerName := reg.ReplaceAllString(header, "${1}")
			headerValue := reg.ReplaceAllString(header, "${2}")
			r.Headers.Set(headerName, headerValue)
		})
	}

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.Request.URL, "\t", r.StatusCode)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(r.Request.URL, "\t", r.StatusCode, "\nError:", err)
	})

	fmt.Print("Started crawler\n")
	// Start scraping
	c.Visit(domain)
	// Wait until threads are finished
	c.Wait()
}

func main() {
	flag.Parse()

	if len(domain) == 0 {
		fmt.Fprintf(os.Stderr, "Flag -domain required\n")
		os.Exit(1)
	}

	if daemon {
		for {
			crawler()
			fmt.Printf("Sleep %v seconds before run crawler again\n", sleep)
			time.Sleep(time.Duration(sleep) * time.Second)
		}
	} else {
		crawler()
	}
}