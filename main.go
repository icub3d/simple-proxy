package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

var (
	addr   = flag.String("addr", ":8080", "the addr:port on which to listen.")
	scheme = flag.String("scheme", "https", "the scheme to rewrite")
	host   = flag.String("host", "", "the host to rewrite")
	path   = flag.String("path", "/", "the path to prefix the original request with")
	query  = flag.String("query", "", "the query parameters to rewrite")
)

func handler(w http.ResponseWriter, r *http.Request) {
	u := url.URL{}
	u.Scheme = *scheme
	u.Host = *host

	if r.URL.Path == "/" || r.URL.Path == "" {
		r.URL.Path = "/index.html"
	}

	u.Path = filepath.Join(*path, r.URL.Path)

	if r.URL.RawQuery == "" || *query == "" {
		u.RawQuery = *query + r.URL.RawQuery
	} else {
		u.RawQuery = *query + "&" + r.URL.RawQuery
	}
	fmt.Println(u.String())
	req, err := http.NewRequest(r.Method, u.String(), nil)
	if err != nil {
		fmt.Println("making url:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("handling %v: %v", r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if strings.HasSuffix(r.URL.Path, ".html") {
		w.Header().Set("Content-Type", "text/html")
	}
	io.Copy(w, resp.Body)
}

func main() {
	flag.Parse()
	http.Handle("/", http.HandlerFunc(handler))
	http.ListenAndServe(*addr, nil)
}
