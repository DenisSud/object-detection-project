package detection

import (
 "image"
)

// Detector interface defines methods for object detection
type Detector interface {
 Detect(img []byte) ([]Object, error)
}

// Object represents a detected object
type Object struct {
 Class      string
 Confidence float64
 BoundingBox image.Rectangle
}
