package config

import (
	"os"
	"strconv"
)

type Config struct {
	ThumbnailWidth  int
	ThumbnailHeight int
	WorkerCount     int
}

func LoadConfig() (*Config, error) {
	// set default values
	c := &Config{
		ThumbnailWidth:  200,
		ThumbnailHeight: 200,
		WorkerCount:     5,
	}

	// override from environment if set
	if widthStr := os.Getenv("THUMBNAIL_WIDTH"); widthStr != "" {
		if w, err := strconv.Atoi(widthStr); err == nil {
			c.ThumbnailWidth = w
		}
	}

	if heightStr := os.Getenv("THUMBNAIL_HEIGHT"); heightStr != "" {
		if h, err := strconv.Atoi(heightStr); err == nil {
			c.ThumbnailHeight = h
		}
	}

	if workerStr := os.Getenv("WORKER_COUNT"); workerStr != "" {
		if wc, err := strconv.Atoi(workerStr); err == nil {
			c.WorkerCount = wc
		}
	}

	return c, nil
}
