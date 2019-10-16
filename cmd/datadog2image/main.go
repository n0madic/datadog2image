// Take a screenshot from DataDog public dashboard
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/n0madic/datadog2image"
)

var (
	height      int64
	httpListen  string
	outputFile  string
	sourceURL   string
	waitLoading int64
	width       int64
)

const refreshTime = 60

func init() {
	flag.Int64Var(&height, "height", 1200, "Screenshot height")
	flag.StringVar(&httpListen, "http", "", "TCP address to HTTP listen on")
	flag.StringVar(&outputFile, "output", "", "Output filename (screenshot.png or index.html)")
	flag.StringVar(&sourceURL, "url", "", "Public dashboard url")
	flag.Int64Var(&waitLoading, "wait", 5, "Dashboard load waiting time in seconds")
	flag.Int64Var(&width, "width", 1920, "Screenshot width")
}

func main() {
	flag.Parse()

	if httpListen != "" {
		indexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			url := r.URL.Query().Get("url")
			if r.URL.Query().Get("width") != "" {
				if n, err := strconv.ParseInt(r.URL.Query().Get("width"), 10, 64); err == nil {
					width = n
				}
			}
			if r.URL.Query().Get("height") != "" {
				if n, err := strconv.ParseInt(r.URL.Query().Get("height"), 10, 64); err == nil {
					height = n
				}
			}
			if url != "" {
				w.Header().Set("Content-Type", "text/html")
				now := time.Now()
				dash := datadog2image.NewDashboard(url).GetScreenshot(width, height, waitLoading).AddTimestamp(&now)
				if dash.Error != nil {
					log.Println(dash.Error.Error())
					w.WriteHeader(http.StatusServiceUnavailable)
					w.Write([]byte(dash.Error.Error()))
				}
				w.Write(dash.HTML(refreshTime))
			} else {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("ERROR: dashboard url required!"))
			}
		})
		http.Handle("/", gziphandler.GzipHandler(indexHandler))
		log.Println("Starting HTTP server at", httpListen)
		log.Fatal(http.ListenAndServe(httpListen, nil))
	} else if sourceURL == "" || outputFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	var buf []byte
	now := time.Now()
	dash := datadog2image.NewDashboard(sourceURL).GetScreenshot(width, height, waitLoading).AddTimestamp(&now)
	if dash.Error != nil {
		log.Fatal(dash.Error)
	}
	if strings.HasSuffix(outputFile, ".png") {
		buf = dash.PNG()
	} else if strings.HasSuffix(outputFile, ".html") || strings.HasSuffix(outputFile, ".htm") {
		buf = dash.HTML(refreshTime)
	} else {
		log.Fatal("Unknown output file format: ", outputFile)
	}

	f, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err := f.Write(buf); err != nil {
		log.Fatal(err)
	}
}
