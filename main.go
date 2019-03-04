package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/ppai-plivo/mnlookup/radix"
	"github.com/ppai-plivo/mnlookup/server"
)

const (
	csvFile = "processed_prefix_data.csv"
)

func main() {

	log.Printf("Opening file: %s\n", csvFile)
	f, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Parsing file: %s\n", csvFile)
	tree, err := radix.New(f)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	runtime.GC()
	log.Printf("Loaded prefixes: %+v\n", tree.CountByNumType())

	svc := server.NewService(tree)

	srv := server.New(svc)
	go func(s *http.Server) {
		log.Println("Starting HTTP service")
		if err := s.ListenAndServe(); err != nil {
			log.Fatalf("http.Server.ListenAndServe() failed: %s\n", err)
		}
	}(srv)

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh
	log.Println("Received interrupt. Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("http.Server.Shutdown() failed: %s\n", err)
	}
}
