Gokogiri
========
LibXML bindings for the Go programming language.
------------------------------------------------
By Zhigang Chen and Hampton Catlin


This is a major rewrite from v0 in the following places:

- Separation of XML and HTML
- Put more burden of memory allocation/deallocation on Go
- Fragment parsing -- no more deep-copy
- Serialization
- Some API adjustment

## Installation

```bash
# Linux
sudo apt-get install libxml2-dev
# Mac
brew install libxml2

go get github.com/Jberlinsky/gokogiri
```

## Running tests

```bash
go test github.com/Jberlinsky/gokogiri/...
```

## Basic example

```go
package main

import (
  "net/http"
  "io/ioutil"
  "github.com/Jberlinsky/gokogiri"
)

func main() {
  // fetch and read a web page
  resp, _ := http.Get("http://www.google.com")
  page, _ := ioutil.ReadAll(resp.Body)

  // parse the web page
  doc, _ := gokogiri.ParseHtml(page)

  // perform operations on the parsed page -- consult the tests for examples

  // important -- don't forget to free the resources when you're done!
  doc.Free()
}
```
