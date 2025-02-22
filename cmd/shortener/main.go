package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	base "github.com/jcoene/go-base62"
)

var urlMap = map[string]string{}
var counter int

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", PostURL)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

func PostURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		GetURL(res, req)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Cannot read body!", http.StatusBadRequest)
		return
	}

	counter++

	urlShort := base.Encode(int64(counter*1000000))
	urlMap[urlShort] = string(body)

	res.Header().Set("Content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", urlShort)))
}

func GetURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET method allowed!", http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(req.URL.Path, "/")
	if id != "" {
		url, ok := urlMap[id]
		if !ok {
			http.Error(res, "URL not found!", http.StatusBadRequest)
			return
		}

	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(res, "URL is empty!", http.StatusBadRequest)
		return
	}
}