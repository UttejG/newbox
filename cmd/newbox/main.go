package main

import (
	"fmt"
	"os"

	"github.com/uttejg/newbox/internal/adapter/output/detector"
)

func main() {
	d := &detector.SystemDetector{}
	platform, err := d.Detect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting platform: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(detector.FormatDetectionInfo(platform))
}
