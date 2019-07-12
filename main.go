package main

import (
	"encoding/json"
	"flag"
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

// Target represents a scraping target
type Target struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func generateFileSdConfig(targetLinks []*url.URL) []byte {

	var targets []Target

	for _, link := range targetLinks {
		t := Target{}
		t.Targets = []string{link.Host}
		t.Labels = map[string]string{
			"__metrics_path__": link.Path,
		}
		targets = append(targets, t)
	}

	d, _ := json.MarshalIndent(targets, "", "  ")
	return d
}

func writeToFile(d []byte, fileSdConfigPath string) {
	file, err := os.Create(fileSdConfigPath)
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

func getTargetLinks(client *http.Client, req *http.Request) []*url.URL {
	var targetLinks []*url.URL
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return targetLinks
	} else {
		defer resp.Body.Close()
	}

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Print("Error loading HTTP response body", err)
	}

	if resp.StatusCode != 200 {
		log.Printf("Got back non-200 response code: %d", resp.StatusCode)
	}

	document.Find("a").Each(func(index int, element *goquery.Selection) {
		link, exists := element.Attr("href")
		if exists {
			parsedURL, err := url.Parse(link)
			if err != nil {
				log.Printf("Error parsing %s: %s", link, err)
			}
			targetLinks = append(targetLinks, parsedURL)
		}
	})

	return targetLinks
}

func main() {

	targetScrapeURL := flag.String("target-source", "", "HTTP URL of the target source")
	targetScrapeInterval := flag.Int64("scrape-interval", 5, "Interval in seconds between consecutive scrapes")
	fileSdConfigPath := flag.String("config-path", "./file_sd_config.json", "Path of the SD config JSON file")

	flag.Parse()

	if len(*targetScrapeURL) == 0 {
		flag.Usage()
		log.Fatal("Please specify target-source")
	}

	req, err := http.NewRequest("GET", *targetScrapeURL, nil)
	if err != nil {
		log.Fatal("Error creating HTTP client", err)
	}
	client := &http.Client{}
	ticker := time.NewTicker(time.Duration(*targetScrapeInterval) * time.Second)
	go func() {
		for range ticker.C {
			targetLinks := getTargetLinks(client, req)
			b := generateFileSdConfig(targetLinks)
			writeToFile(b, *fileSdConfigPath)
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
