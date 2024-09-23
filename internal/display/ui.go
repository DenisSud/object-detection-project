package display

import (
	"gocv.io/x/gocv"
)

type UI struct {
	window *gocv.Window
}

func NewUI() (*UI, error) {
	window := gocv.NewWindow("Object Detection")
	return &UI{window: window}, nil
}

func (u *UI) Display(frame gocv.Mat, objects []Object) {
	// Draw bounding boxes and labels on the frame
	for _, obj := range objects {
		gocv.Rectangle(&frame, obj.BoundingBox, color.RGBA{255, 0, 0, 255}, 2)
		gocv.PutText(&frame, obj.Class, image.Pt(obj.BoundingBox.Min.X, obj.BoundingBox.Min.Y-10),
			gocv.FontHersheyPlain, 1.2, color.RGBA{255, 0, 0, 255}, 2)
	}
	u.window.IMShow(frame)
	u.window.WaitKey(1)
}

func (u *UI) Close() {
	u.window.Close()
}
