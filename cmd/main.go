package main

import (
	"log"
	"github.com/DenisSud/project/internal/display"
	"github.com/DenisSud/project/internal/detection"
	"github.com/DenisSud/project/internal/notification"
	"github.com/DenisSud/project/internal/capture"
)

func main() {
	// Initialize components
	cap, err := capture.NewCapture()
	if err != nil {
		log.Fatalf("Failed to initialize capture: %v", err)
	}
	defer cap.Close()

	detector, err := detection.NewDetector()
	if err != nil {
		log.Fatalf("Failed to initialize detector: %v", err)
	}

	notifier, err := notification.NewNotifier()
	if err != nil {
		log.Fatalf("Failed to initialize notifier: %v", err)
	}

	ui, err := display.NewUI()
	if err != nil {
		log.Fatalf("Failed to initialize UI: %v", err)
	}
	defer ui.Close()

	// Main processing loop
	for {
		frame, err := cap.GetFrame()
		if err != nil {
			log.Printf("Error capturing frame: %v", err)
			continue
		}

		objects, err := detector.Detect(frame)
		if err != nil {
			log.Printf("Error detecting objects: %v", err)
			continue
		}

		if len(objects) > 0 {
			notifier.Notify(objects)
		}

		ui.Display(frame, objects)
	}
}
