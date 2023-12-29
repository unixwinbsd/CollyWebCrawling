## Colly Web Crawler
Based on https://github.com/gocolly/colly


Continue with the following command.
# go mod colly
# go get go.uber.org/automaxprocs
# go get github.com/gocolly/colly

To run the web crawler we have to add several "flags".
-daemon: Run crawler on daemon mode
-delay int: Duration to wait before creating a new request to the matching domains (default 1)
-domain string: Set url for crawling. Example: https://example.com
-header string: Set header for crawler request. Example: header_name:header_value
-parallelism int: Parallelism is the number of the maximum allowed concurrent requests of the matching domains (default 2).
-sleep int: Time in seconds to wait before run crawler again (default 60)

Example of use:
# go run crawlingblogsite.go -domain https://www.unixwinbsd.site -header header_name:header_value -daemon



