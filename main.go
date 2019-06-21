package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Target struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

var targetLinks []*url.URL

func processElement(index int, element *goquery.Selection) {

	link, exists := element.Attr("href")
	if exists {
		parsedUrl, err := url.Parse(link)
		if err != nil {
			log.Printf("Error parsing %s", link, err)
		}
		targetLinks = append(targetLinks, parsedUrl)
	}

}

func generateFileSdConfig() {

	var targets []Target

	for _, link := range targetLinks {
		t := Target{}
		t.Targets = []string{link.Host}
		t.Labels = map[string]string{
			"__metrics_path__": link.Path,
		}
		targets = append(targets, t)
	}

	d, _ := json.Marshal(targets)
	file, err := os.Create("file_sd_targets.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(d)
	if err != nil {
		log.Fatal(err)
	}
	file.Sync()

}

func main() {

	req, err := http.NewRequest("GET", os.Args[1], nil)
	if err != nil {
		log.Fatal("Error creating HTTP client", err)
	}
	client := &http.Client{}

	ticker := time.NewTicker(5000 * time.Millisecond)
	go func() {
		for range ticker.C {
			resp, err := client.Do(req)
			if err != nil {
				log.Print(err)
				continue
			} else {
				defer resp.Body.Close()
			}
			// Create a goquery document from the HTTP response
			document, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Print("Error loading HTTP response body", err)
			}

			if resp.StatusCode != 200 {
				log.Printf("Got back non-200 response code: %s", resp.StatusCode)
			}

			document.Find("a").Each(processElement)
			generateFileSdConfig()

			targetLinks = targetLinks[:0]
		}
	}()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		done <- true
	}()

	<-done
	fmt.Println("exiting")
}
