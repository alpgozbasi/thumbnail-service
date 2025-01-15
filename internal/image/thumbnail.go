package image

import (
	"fmt"
	"path/filepath"

	"github.com/disintegration/imaging"
)

func GenerateThumbnail(inputPath string, outputDir string, width, height int) (string, error) {
	// open original image
	src, err := imaging.Open(inputPath)
	if err != nil {
		return "", err
	}

	// resize
	dst := imaging.Resize(src, width, height, imaging.Lanczos)

	// prepare output path
	filename := filepath.Base(inputPath)
	thumbName := fmt.Sprintf("%s_%dx%d%s", filenameWithoutExt(filename), width, height, filepath.Ext(filename))
	outputPath := filepath.Join(outputDir, thumbName)

	// save thumbnail
	err = imaging.Save(dst, outputPath)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

// filenamewithoutext returns name part without extension
func filenameWithoutExt(fname string) string {
	ext := filepath.Ext(fname)
	return fname[:len(fname)-len(ext)]
}
