package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

const (
	URL string = "https://xkcd.com"
)

func fetch(n int) (*Result, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	url := strings.Join([]string{URL, strconv.Itoa(n), "info.0.json"}, "/")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("http request error: %e", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http error doing request: %e", err)
	}

	var data *Result
	if resp.StatusCode != http.StatusOK {
		data = new(Result)
	} else {
		if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
			return nil, fmt.Errorf("json decode error: %e", err)
		}
	}

	defer resp.Body.Close()
	return data, nil
}

func main() {
	n := 200
	result, err := fetch(n)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(result.Title)
}
