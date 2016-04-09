package main 

import (
    "fmt"
    "net/http"
    "net/url"
)

func main() {
    proxy := func(_ *http.Request) (*url.URL, error) {
        return url.Parse("http://127.0.0.1:8087")
    }

    transport := &http.Transport{Proxy: proxy}

    client := &http.Client{Transport: transport}
    resp, err := client.Get("http://www.google.com")

    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(resp)
}

