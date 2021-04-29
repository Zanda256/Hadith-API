package main

import (
	"HadithAPI/cmd/handlers"
	"HadithAPI/cmd/scheduler"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(2)
	l := log.New(os.Stdout, "Hadith-API", log.LstdFlags)
	hh := handlers.NewHadiths(l)
	hh.Initiate()

	scheduler.Jt = scheduler.NewJobTicker()

	sm := http.NewServeMux()
	sm.Handle("/", hh)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      sm,
		IdleTimeout:  12 * time.Minute,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	s.ListenAndServe()
}
