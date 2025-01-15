package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alpgozbasi/thumbnail-service/internal/config"
	"github.com/alpgozbasi/thumbnail-service/internal/server"
	"github.com/alpgozbasi/thumbnail-service/internal/worker"
)

func main() {
	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// create channels
	jobs := make(chan worker.Job, 100)
	results := make(chan worker.Result, 100)

	// create waitgroup
	var wg sync.WaitGroup

	// start worker pool
	worker.Start(cfg, jobs, results, &wg)

	// start http server
	if err := server.Run(cfg, jobs, results, &wg); err != nil {
		log.Fatalf("server run error: %v", err)
	}

	// handle interrupt signals (ctrl+c, kill)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// wait for a signal
	<-sig
	log.Println("shutdown signal received")

	// close job channel so workers can finish
	close(jobs)
	// wait for all workers
	wg.Wait()
	// close results channel
	close(results)

	log.Println("graceful shutdown complete")
}
