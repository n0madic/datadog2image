// Take a screenshot from DataDog public dashboard
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/n0madic/datadog2image"
)

var (
	httpListen  string
	outputFile  string
	sourceURL   string
	waitLoading int64
)

const refreshTime = 60

func init() {
	flag.StringVar(&httpListen, "http", "", "TCP address to HTTP listen on")
	flag.StringVar(&outputFile, "output", "", "Output filename (screenshot.png or index.html)")
	flag.StringVar(&sourceURL, "url", "", "Public dashboard url")
	flag.Int64Var(&waitLoading, "wait", 4, "Dashboard load waiting time in seconds")
}

func main() {
	flag.Parse()

	if httpListen != "" {
		http.HandleFunc("/", indexRequest)
		log.Println("Starting HTTP server at", httpListen)
		log.Fatal(http.ListenAndServe(httpListen, nil))
	} else if sourceURL == "" || outputFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	var buf []byte
	now := time.Now()
	dash := datadog2image.NewDashboard(sourceURL).GetScreenshot(waitLoading).AddTimestamp(&now)
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

func indexRequest(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url != "" {
		w.Header().Set("Content-Type", "text/html")
		now := time.Now()
		dash := datadog2image.NewDashboard(url).GetScreenshot(waitLoading).AddTimestamp(&now)
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
}
