package main

import (
	"net/http"
	"strings"
)

type Getter struct {
	Url    string
	Query  map[string]string
	Header map[string]string
}

func (g Getter) Get() (*http.Response, error) {
	var urlSlice []string
	for k, v := range g.Query {
		urlSlice = append(urlSlice, k+"="+v)
	}
	var url = g.Url
	if len(urlSlice) != 0 {
		url = strings.Join(urlSlice, "&")
		url = g.Url + "?" + url
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range g.Header {
		req.Header.Add(k, v)
	}
	var client = &http.Client{}
	return client.Do(req)

}
