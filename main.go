package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
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

type Job struct {
	number int
}

const (
	URL string = "https://xkcd.com"
)

var jobs = make(chan Job, 100)
var results = make(chan Result, 100)
var resultCollection []Result

func getResults(done chan bool) {
	for result := range results {
		if result.Num != 0 {
			fmt.Println("Retrieving issue", result.Num)
			resultCollection = append(resultCollection, result)
		}
	}
	done <- true
}

func worker(wg *sync.WaitGroup) {
	for job := range jobs {
		result, err := fetch(job.number)
		if err != nil {
			log.Println("error in fetching: ", err)
		}
		results <- *result
	}
	wg.Done()
}

func createWorkerPool(numWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(&wg)
	}

	wg.Wait()
	close(results)
}

func allocateJobs(numJobs int) {
	for i := 0; i < numJobs; i++ {
		jobs <- Job{i + 1}
	}
	close(jobs)
}

func fetch(n int) (*Result, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	url := strings.Join([]string{URL, fmt.Sprintf("%d", n), "info.0.json"}, "/")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("http request error: %e", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http error doing request: %e", err)
	}

	var data Result

	// error from web service, empty struct to avoid disruption of process
	if resp.StatusCode != http.StatusOK {
		data = Result{}
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("json err: %v", err)
		}
	}

	resp.Body.Close()

	return &data, nil
}

func writeToFile(data []byte) error {
	f, err := os.Create("xkcd.json")
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	numJobs := 3000
	go allocateJobs(numJobs)

	done := make(chan bool)
	go getResults(done)

	numWorkers := 100
	createWorkerPool(numWorkers)

	<-done

	data, err := json.MarshalIndent(resultCollection, "", "    ")
	if err != nil {
		log.Fatalf("error in marshalling result collection: %e", err)
	}

	err = writeToFile(data)
	if err != nil {
		log.Fatal(err)
	}
}
