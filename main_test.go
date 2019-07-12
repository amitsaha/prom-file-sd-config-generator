package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func cleanup() {
	err := os.Remove("/tmp/test.json")
	if err != nil {
		panic(err)
	}

}
func TestMain(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `
<a href="http://127.0.0.1:9100/foo/bar">target1</a>
<a href="http://127.0.0.1:9300/foo/bar">target3</a>`)
	}))
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		log.Fatal("Error creating HTTP client", err)
	}
	client := &http.Client{}
	targetLinks := getTargetLinks(client, req)

	var expectedT []Target
	t1 := Target{
		Targets: []string{"127.0.0.1:9100"},
		Labels: map[string]string{
			"__metrics_path__": "/foo/bar",
		},
	}
	t2 := Target{
		Targets: []string{"127.0.0.1:9300"},
		Labels: map[string]string{
			"__metrics_path__": "/foo/bar",
		},
	}

	expectedT = append(expectedT, t1, t2)
	expectedTargets, err := json.MarshalIndent(expectedT, "", "  ")
	if err != nil {
		panic(err)
	}
	b := generateFileSdConfig(targetLinks)
	if string(expectedTargets) != string(b) {
		t.Errorf("Expected %s, Got %s\n", expectedTargets, string(b))
	}

	writeToFile(b, "/tmp/test.json")
	defer cleanup()

	data, err := ioutil.ReadFile("/tmp/test.json")
	if err != nil {
		panic(err)
	}
	if string(data) != string(b) {
		t.Errorf("Expected %s, Got %s\n", expectedTargets, string(b))
	}
}
