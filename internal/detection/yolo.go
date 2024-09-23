package detection

import (
	"fmt"
	"image"
	"math"

	"gocv.io/x/gocv"
	ort "github.com/yalue/onnxruntime_go"
)

type YOLOv8Detector struct {
	session    ort.AdvancedSession
	inputName  string
	outputName string
	classes    []string
}

func NewYOLOv8Detector(modelPath string) (*YOLOv8Detector, error) {
	ort.SetSharedLibraryPath(getSharedLibPath())
	err := ort.InitializeEnvironment()
	if err != nil {
		return nil, err
	}

	inputShape := ort.NewShape(1, 3, 640, 640)
	inputTensor, err := ort.NewEmptyTensor[float32](inputShape)
	if err != nil {
		return nil, err
	}

	outputShape := ort.NewShape(1, 84, 8400)
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		return nil, err
	}

	options, err := ort.NewSessionOptions()
	if err != nil {
		return nil, err
	}
	defer options.Destroy()

	session, err := ort.NewAdvancedSession(modelPath,
		[]string{"images"}, []string{"output0"},
		[]ort.ArbitraryTensor{inputTensor}, []ort.ArbitraryTensor{outputTensor},
		options)
	if err != nil {
		return nil, err
	}

	return &YOLOv8Detector{
		session:    session,
		inputName:  "images",
		outputName: "output0",
		classes:    yolo_classes,
	}, nil
}

func (y *YOLOv8Detector) Detect(img gocv.Mat) ([]Object, error) {
	resized := gocv.NewMat()
	defer resized.Close()
	gocv.Resize(img, &resized, image.Point{X: 640, Y: 640}, 0, 0, gocv.InterpolationLinear)

	blob := gocv.BlobFromImage(resized, 1.0/255.0, image.Pt(640, 640), gocv.NewScalar(0, 0, 0, 0), true, false)
	defer blob.Close()

	input := blob.DataPtrFloat32()
	inputTensor := y.session.Inputs()[0]
	copy(inputTensor.GetData(), input)

	err := y.session.Run()
	if err != nil {
		return nil, err
	}

	output := y.session.Outputs()[0].GetData()
	return y.postprocess(output, img.Cols(), img.Rows())
}

func (y *YOLOv8Detector) postprocess(output []float32, imgWidth, imgHeight int) ([]Object, error) {
	var objects []Object
	rows := 8400
	dimensions := 84

	for i := 0; i < rows; i++ {
		row := output[i*dimensions : (i+1)*dimensions]
		classId, confidence := y.getMaxConfidence(row[4:])
		if confidence < 0.25 {
			continue
		}

		x := row[0]
		y := row[1]
		w := row[2]
		h := row[3]

		left := (x - w/2) * float32(imgWidth)
		top := (y - h/2) * float32(imgHeight)
		right := (x + w/2) * float32(imgWidth)
		bottom := (y + h/2) * float32(imgHeight)

		objects = append(objects, Object{
			Class:      y.classes[classId],
			Confidence: float64(confidence),
			BoundingBox: image.Rect(
				int(math.Max(0, float64(left))),
				int(math.Max(0, float64(top))),
				int(math.Min(float64(imgWidth), float64(right))),
				int(math.Min(float64(imgHeight), float64(bottom))),
			),
		})
	}

	return objects, nil
}

func (y *YOLOv8Detector) getMaxConfidence(scores []float32) (int, float32) {
	maxIndex := 0
	maxScore := float32(0)
	for i, score := range scores {
		if score > maxScore {
			maxScore = score
			maxIndex = i
		}
	}
	return maxIndex, maxScore
}

// Implement getSharedLibPath() function here
// This function should return the path to the ONNX Runtime shared library
// based on the current operating system and architecture

var yolo_classes = []string{
	"person", "bicycle", "car", "motorcycle", "airplane", "bus", "train", "truck", "boat",
	"traffic light", "fire hydrant", "stop sign", "parking meter", "bench", "bird", "cat",
	"dog", "horse", "sheep", "cow", "elephant", "bear", "zebra", "giraffe", "backpack",
	"umbrella", "handbag", "tie", "suitcase", "frisbee", "skis", "snowboard", "sports ball",
	"kite", "baseball bat", "baseball glove", "skateboard", "surfboard", "tennis racket",
	"bottle", "wine glass", "cup", "fork", "knife", "spoon", "bowl", "banana", "apple",
	"sandwich", "orange", "broccoli", "carrot", "hot dog", "pizza", "donut", "cake", "chair",
	"couch", "potted plant", "bed", "dining table", "toilet", "tv", "laptop", "mouse",
	"remote", "keyboard", "cell phone", "microwave", "oven", "toaster", "sink", "refrigerator",
	"book", "clock", "vase", "scissors", "teddy bear", "hair drier", "toothbrush",
}
