package worker

import (
	"log"
	"sync"

	"github.com/alpgozbasi/thumbnail-service/internal/config"
	"github.com/alpgozbasi/thumbnail-service/internal/image"
)

// job carries the input path of the image
type Job struct {
	ID        string
	InputPath string
	OutputDir string
}

// result carries the output path of the generated thumbnail or any error
type Result struct {
	JobID      string
	OutputPath string
	Err        error
}

// start spawns multiple worker goroutines based on cfg.WorkerCount
func Start(cfg *config.Config, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	for w := 1; w <= cfg.WorkerCount; w++ {
		wg.Add(1)
		go worker(cfg, w, jobs, results, wg)
	}
}

// worker listens for jobs, generates thumbnails, and sends results
func worker(cfg *config.Config, workerID int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		log.Printf("worker %d received job: %s", workerID, job.ID)

		out, err := image.GenerateThumbnail(
			job.InputPath,
			job.OutputDir,
			cfg.ThumbnailWidth,
			cfg.ThumbnailHeight,
		)

		results <- Result{
			JobID:      job.ID,
			OutputPath: out,
			Err:        err,
		}
	}
}
