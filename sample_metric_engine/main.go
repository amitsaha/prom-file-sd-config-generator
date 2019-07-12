package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	SCHEME    string = "http"
	HOSTNAME  string = "127.0.0.1"
	PORT      int    = 2112
	indexPage string = ""
)

// Record metrics example from prometheus goclient docs
func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, indexPage)
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func main() {
	recordMetrics()

	r := mux.NewRouter()

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/", indexHandler)

	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		if t != "/" {
			indexPage += fmt.Sprintf(`<p><a href="%s://%s:%d%v">%v</a></p>`, SCHEME, HOSTNAME, PORT, t, t) + "\n"
		}
		return nil
	})

	http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
}
