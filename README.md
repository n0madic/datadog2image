# datadog2image
Take a screenshot from DataDog public dashboard. Good for monitoring on TV.

## Install

From source:

```
go get -u github.com/n0madic/datadog2image/cmd/datadog2image
```

## Help

```
Usage of datadog2image:
  -http string
	TCP address to HTTP listen on
  -output string
	Output filename (screenshot.png or index.html)
  -url string
	Public dashboard url
  -wait int
	Dashboard load waiting time in seconds (default 4)
```

## Usage

Save to file:

```
$ datadog2image -url https://p.datadoghq.com/sb/... -output screenshot.png
```
or
```
$ datadog2image -url https://p.datadoghq.com/sb/... -output index.html
```

Run as web service:

```
$ datadog2image -http :8000
$ curl http://localhost:8000/?url=https://p.datadoghq.com/sb/...
```

## Use as package

```go
package main

import (
  "os"
  "time"

  "github.com/n0madic/datadog2image"
)

func main() {
    now := time.Now()
    dash := datadog2image.NewDashboard("https://p.datadoghq.com/...").GetScreenshot(4).AddTimestamp(&now)
    if dash.Error != nil {
        panic(dash.Error)
    }
    buf := dash.PNG()
    f, err := os.Create("screenshot.png")
    if err != nil {
        panic(err)
    }
    defer f.Close()
    if _, err := f.Write(buf); err != nil {
        panic(err)
    }
}
```
