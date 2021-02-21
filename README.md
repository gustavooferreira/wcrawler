# WCrawler

WCrawler is a simple web crawler CLI tool.

![gif usage image](docs/images/usage.gif "Usage example")

# Usage

Exploring the Web:

```
❯ wcrawler explore -h
Explore the web by following links up to a pre-determined depth

Usage:
  wcrawler explore URL [flags]

Flags:
  -f, --file string    file to save results (default "./web_graph.json")
  -d, --depth uint     depth of recursion (default 10)
  -s, --stats          show live stats (default true)
  -t, --timeout uint   HTTP requests timeout (default 10)
  -w, --workers uint   number of workers making concurrent requests (default 10)
```

Visualizing the graph in the browser:

```
❯ wcrawler view -h
View web links relationships in the browser

Usage:
  wcrawler view [flags]

Flags:
  -f, --file string   file containing the data (default "./web_graph.json")

```

This will generate a webpage and load it on your default browser.

# Example

This will crawl the web starting at the `example.com` website up to a max of 8 depth levels, using 5 workers with a 6 second timeout per request and saving the collected data to `/tmp/result.json`.

```
wcrawler explore https://example.com -d 8 -w 5 -t 6 -f /tmp/result.json
```

This command will then generate an HTML file with a graph view of the data collected and load it onto the default web browser.

```
wcrawler view -f /tmp/result.json
```

---

# Staying up to date

To update wcrawler to the latest version, use `go get -u github.com/gustavooferreira/wcrawler`.

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

Please feel free to submit issues, fork the repository and send pull requests!

When submitting an issue, we ask that you please include a complete test function that demonstrates the issue.

---

# License

This project is licensed under the terms of the MIT license.
