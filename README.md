# WCrawler

[![Build Status](https://travis-ci.com/gustavooferreira/wcrawler.svg?branch=master)](https://travis-ci.com/gustavooferreira/wcrawler)
[![codecov](https://codecov.io/gh/gustavooferreira/wcrawler/branch/master/graph/badge.svg)](https://codecov.io/gh/gustavooferreira/wcrawler)
[![Go Report Card](https://goreportcard.com/badge/github.com/gustavooferreira/wcrawler)](https://goreportcard.com/report/github.com/gustavooferreira/wcrawler)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gustavooferreira/wcrawler)](https://pkg.go.dev/github.com/gustavooferreira/wcrawler)

WCrawler is a simple web crawler CLI tool.

**NOTE:** This tool was created mainly for practice purposes and therefore doesn't rely on any library that facilitates crawling.

![Usage example video](https://user-images.githubusercontent.com/17534422/109546768-85aec680-7ac2-11eb-8c72-2dbf7c7223a8.mp4)

**Watch this &#9757;**

\<according to [this](https://twitter.com/natfriedman/status/1365393828622921728), github is supposed to be able to display mp4 videos on markdown, but doesn't seem to work. Let's hope it's a Blue/Green deployment thing and wait>

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
  -m, --treemode          doesn't add links which would point back to known nodes
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

Spheres are coloured based on the URL subdomain, you can pan, tilt and rotate the scene, drag the spheres and move them around, hover to check the URL they represent and click on them to go straight to that URL.

**NOTE:** If you want to see a nice graph, make sure to run `wcrawler explore` with the `-m` flag.
Tree mode doesn't create links back to the original URLs making for much nicer visualizations.
Its utility? None, but the graphs are undeniably more beautiful.

Naturally, if you want a proper graph of the links visited and where they point to, just disregard the `-m` option. Don't try to visualize that, however, cos it's going to look ugly, if not freeze your browser entirely. Consider yourself warned :)

# Example

This following command will crawl the web starting at the `example.com` website up to a max of 8 depth levels, using 5 workers with a 6 second timeout per request and saving the collected data to `/tmp/result.json`.

```
wcrawler explore https://example.com -d 8 -w 5 -t 6 -o /tmp/result.json
```

This following command will then generate an HTML file with a graph view of the data collected and load it onto the default web browser. Only try to visualize the graph if you have specified the `-m` option! It's going to be the wrong graph, but it's going to look nice!

```
wcrawler view -i /tmp/result.json
```

---

# Considerations

Here I'm going to discuss the design decisions and a few caveats, but only when I'm actually done with the project.

Still have a few more things to do like:

- Add logic to fetch website's robots.txt file and adhere to whatever it's in there. At the moment we are just crawling everything (feeling like an outlaw here at the minute)
- Show last 10 errors in the CLI while crawling
- Make output more colorful
- Docs, docs and more docs
- Increase coverage and run some benchmarks (I'm pretty sure I can speed up some parts and reduce allocations, even though this program is I/O bound more than anything else so won't benefit much from these optimizations, but practice is practice)
- Add golangci-lint to travis-ci (cos it's quite nice)
- Organize code in a way that makes it for a useful library (mostly done)

---

# Third party libraries being used (directly):

Could have written the whole thing without using any library, but reusability is not a bad idea at all!

The only rule I had was to not use any library that facilitates crawling.

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

The `wcrawler` binary will be placed inside the `bin/` folder.

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

## Free tip

> If you run `make` without any targets, it will display all options available on the makefile followed by a short description.

---

# Contributing

I'd normally be more than happy to accept pull requests, but given that I've created this project with the sole intent of practicing, it doesn't make sense for me to accept other people's work.

However, feel free to fork the project and add whatever new features you feel like.

I'd still be glad if you notice a bug and report it by opening an issue.

---

# License

This project is licensed under the terms of the MIT license.
