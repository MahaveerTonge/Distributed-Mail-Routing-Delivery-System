package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	mode         = flag.String("mode", "async", "sync or async")
	downstreamMs = flag.Int("downstream-ms", 200, "simulated SMTP delivery delay in ms")
	workers      = flag.Int("workers", 8, "async worker pool size")
	accepted     atomic.Int64
	delivered    atomic.Int64
	queue        = make(chan int64, 50000)
)

func main() {
	flag.Parse()
	if *mode == "async" {
		for i := 0; i < *workers; i++ {
			go func() {
				for range queue {
					time.Sleep(time.Duration(*downstreamMs) * time.Millisecond)
					delivered.Add(1)
				}
			}()
		}
	}

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		id := accepted.Add(1)
		if *mode == "sync" {
			time.Sleep(time.Duration(*downstreamMs) * time.Millisecond)
			delivered.Add(1)
			w.WriteHeader(http.StatusOK)
			return
		}
		select {
		case queue <- id:
			w.WriteHeader(http.StatusAccepted)
		default:
			http.Error(w, "queue full", http.StatusServiceUnavailable)
		}
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "accepted=%d delivered=%d queued=%d\n",
			accepted.Load(), delivered.Load(), len(queue))
	})

	fmt.Printf("Running in %s mode (downstream=%dms, workers=%d) on :8080\n", *mode, *downstreamMs, *workers)
	http.ListenAndServe(":8080", nil)
}