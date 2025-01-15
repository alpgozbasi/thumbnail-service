package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/alpgozbasi/thumbnail-service/internal/config"
	"github.com/alpgozbasi/thumbnail-service/internal/worker"
)

// run starts the http server and listens for results
func Run(cfg *config.Config, jobs chan<- worker.Job, results <-chan worker.Result, wg *sync.WaitGroup) error {
	// prepare upload and thumbnail dirs
	if err := os.MkdirAll("./uploads", 0755); err != nil {
		return err
	}
	if err := os.MkdirAll("./thumbnails", 0755); err != nil {
		return err
	}

	// handle file uploads
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// check request method
		if r.Method != http.MethodPost {
			http.Error(w, "use post", http.StatusMethodNotAllowed)
			return
		}

		// parse form file
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// create local file
		localPath := filepath.Join("./uploads", header.Filename)
		out, err := os.Create(localPath)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		// copy file
		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// send job
		jobID := "job-" + header.Filename
		j := worker.Job{
			ID:        jobID,
			InputPath: localPath,
			OutputDir: "./thumbnails",
		}
		jobs <- j

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("file uploaded and job queued"))
	})

	// listen for worker results
	go func() {
		for res := range results {
			if res.Err != nil {
				log.Printf("job %s error: %v", res.JobID, res.Err)
				continue
			}
			log.Printf("job %s completed: %s", res.JobID, res.OutputPath)
		}
	}()

	// configure and start http server
	srv := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	return nil
}
