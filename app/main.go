package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	started atomic.Bool
	ready   atomic.Bool
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("PID: %d", os.Getpid())

	sigCh := make(chan os.Signal, 32)
	signal.Notify(sigCh)

	go func() {
		for sig := range sigCh {
			log.Printf("received signal: %v", sig)

			switch sig {
			case syscall.SIGTERM:
				log.Println("SIGTERM: marking pod as NOT READY")
				ready.Store(false)

				log.Println("sleeping for 20s before exit...")
				time.Sleep(20 * time.Second)

				log.Println("exiting")
				os.Exit(0)

			case syscall.SIGINT:
				log.Println("SIGINT: exiting")
				os.Exit(0)

			default:
				log.Printf("no special handling for %v", sig)
			}
		}
	}()

	log.Println("doing initialization")
	time.Sleep(15 * time.Second)
	log.Println("initialization done")

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			n := rand.IntN(10) + 1
			log.Printf("timer: generated number %d", n)

			switch n {
			case 10:
				log.Println("random 10 -> exiting with code 1")
				os.Exit(1)
			case 1:
				log.Println("random 1 -> exiting with code 0")
				os.Exit(0)
			}
		}
	}()

	started.Store(true)
	ready.Store(true)

	http.HandleFunc("/startup", func(w http.ResponseWriter, r *http.Request) {
		log.Println("startup probe")
		if !started.Load() {
			http.Error(w, "starting", http.StatusServiceUnavailable)
			return
		}
		fmt.Fprint(w, "ok")
	})

	http.HandleFunc("/liveness", func(w http.ResponseWriter, r *http.Request) {
		n := rand.IntN(10) + 1

		if n == 10 {
			log.Println("simulatring abnormal application exit")
			os.Exit(1)
		}
		log.Println("liveness probe, random number:", n)
		fmt.Fprint(w, "ok")
	})

	http.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		n := rand.IntN(10) + 1

		if n == 10 {
			log.Println("simulatring normal application exit")
			os.Exit(0)
		}
		log.Println("readiness probe, random number:", n)
		fmt.Fprint(w, "ok")
	})

	log.Println("HTTP server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
