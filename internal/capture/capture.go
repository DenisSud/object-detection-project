package capture

import (
	"fmt"

	"gocv.io/x/gocv"
)

type Capture struct {
	webcam *gocv.VideoCapture
}

func NewCapture() (*Capture, error) {
	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		return nil, err
	}
	return &Capture{webcam: webcam}, nil
}

func (c *Capture) GetFrame() (gocv.Mat, error) {
	frame := gocv.NewMat()
	if ok := c.webcam.Read(&frame); !ok {
		return frame, fmt.Errorf("unable to read frame")
	}
	return frame, nil
}

func (c *Capture) Close() {
	c.webcam.Close()
}
