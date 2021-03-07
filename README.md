# WCrawler

WCrawler is a simple web crawler CLI tool.

**NOTE:** This tool was created mainly for practice purposes and therefore doesn't rely on any library that facilitates crawling.

![Usage example video](https://user-images.githubusercontent.com/17534422/109546768-85aec680-7ac2-11eb-8c72-2dbf7c7223a8.mp4)

# Usage

Exploring the Web:

```
❯ wcrawler explore --help
Explore the web by following links up to a pre-determined depth.
A depth of zero means no limit.

Usage:
  wcrawler explore URL [flags]

Flags:
  -d, --depth uint        depth of recursion (default 5)
  -h, --help              help for explore
  -s, --nostats           don't show live stats
  -o, --output string     file to save results (default "./web_graph.json")
  -r, --retry uint        retry requests when they timeout (default 2)
  -z, --stayinsubdomain   follow links only in the same subdomain
  -t, --timeout uint      HTTP requests timeout in seconds (default 10)
  -w, --workers uint      number of workers making concurrent requests (default 100)
```

Visualizing the graph in the browser:

```
❯ wcrawler view --help
View web links relationships in the browser

Usage:
  wcrawler view [flags]

Flags:
  -h, --help            help for view
  -i, --input string    file containing the data (default "./web_graph.json")
  -n, --noautoopen      don't open browser automatically
  -o, --output string   HTML output file (default "./web_graph.html")
```

This will generate a webpage and load it on your default browser.

# Example

This will crawl the web starting at the `example.com` website up to a max of 8 depth levels, using 5 workers with a 6 second timeout per request and saving the collected data to `/tmp/result.json`.

```
wcrawler explore https://example.com -d 8 -w 5 -t 6 -o /tmp/result.json
```

This command will then generate an HTML file with a graph view of the data collected and load it onto the default web browser.

```
wcrawler view -i /tmp/result.json
```

---

# Third party libraries being used (directly):

```
- github.com/gosuri/uilive     [updating terminal output in realtime]
- github.com/spf13/cobra       [CLI args and flags parsing]
- github.com/stretchr/testify  [writing unit tests]
- golang.org/x/net             [HTML parsing]
- github.com/oleiade/lane      [Provides a Queue data structure implementation]
```

---

# Staying up to date

To update wcrawler to the latest version, use `go get -u github.com/gustavooferreira/wcrawler`.

---

# Build

To build this project run:

```
make build
```

---

# Tests

To run tests:

```
make test
```

To get coverage:

```
make coverage
```

---

# Contributing

I'd normally be more than happy to accept pull requests, but given that I've created this project with the sole intent of practicing, it doesn't make sense for me to accept other people's work.

However, feel free to fork the project and add whatever new features you fill like.

I'd still be glad if you notice a bug and report it by opening an issue.

---

# License

This project is licensed under the terms of the MIT license.
