package countclassifier

import (
	"context"
	"image"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/vision"
	viz "go.viam.com/rdk/vision"
	"go.viam.com/rdk/vision/classification"
	objdet "go.viam.com/rdk/vision/objectdetection"
	"go.viam.com/rdk/vision/viscapture"
)

const (
	// ModelName is the name of the model
	ModelName = "count-classifier"
	// OverflowLabel is the label if the counts exceed what was specified by the user
	OverflowLabel = "Overflow"
)

var (
	// Model is the resource
	Model            = resource.NewModel("viam", "vision-summary", ModelName)
	errUnimplemented = errors.New("unimplemented")
)

func init() {
	resource.RegisterService(vision.API, Model, resource.Registration[vision.Service, *Config]{
		Constructor: newCountClassifier,
	})
}

// Config contains names for necessary resources (camera and vision service)
type Config struct {
	DetectorName    string             `json:"detector_name"`
	ChosenLabels    map[string]float64 `json:"chosen_labels"`
	CountThresholds map[string]int     `json:"count_thresholds"`
}

// Validate validates the config and returns implicit dependencies,
// this Validate checks if the camera and detector exist for the module's vision model.
func (cfg *Config) Validate(path string) ([]string, error) {
	if cfg.DetectorName == "" {
		return nil, errors.New("attribute detector_name cannot be left blank")
	}
	if len(cfg.CountThresholds) == 0 {
		return nil, errors.New("attribute count_thresholds is required")
	}
	testMap := map[int]string{}
	for label, v := range cfg.CountThresholds {
		if _, ok := testMap[v]; ok {
			return nil, errors.Errorf("cannot have two labels for the same threshold in count_thresholds. Threshold value %v appears more than once", v)
		}
		testMap[v] = label
	}
	return []string{cfg.DetectorName}, nil
}

// Bin stores the thresholds that turns counts into labels
type Bin struct {
	UpperBound int
	Label      string
}

// NewThresholds creates a list of thresholds for labeling counts
func NewThresholds(t map[string]int) []Bin {
	// first invert the map, Validate ensures a 1-1 mapping
	thresholds := map[int]string{}
	for label, val := range t {
		thresholds[val] = label
	}
	out := []Bin{}
	keys := []int{}
	for k := range thresholds {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, key := range keys {
		b := Bin{key, thresholds[key]}
		out = append(out, b)
	}
	return out
}

type countcls struct {
	resource.Named
	logger     logging.Logger
	properties vision.Properties
	detName    string
	detector   vision.Service
	labels     map[string]float64
	thresholds []Bin
}

func newCountClassifier(
	ctx context.Context,
	deps resource.Dependencies,
	conf resource.Config,
	logger logging.Logger) (vision.Service, error) {
	cc := &countcls{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
		properties: vision.Properties{
			ClassificationSupported: true,
			DetectionSupported:      true,
			ObjectPCDsSupported:     false,
		},
	}

	if err := cc.Reconfigure(ctx, deps, conf); err != nil {
		return nil, err
	}
	return cc, nil
}

// Reconfigure resets the underlying detector as well as the thresholds and labels for the count
func (cc *countcls) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	countConf, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return errors.Errorf("Could not assert proper config for %s", ModelName)
	}
	cc.detName = countConf.DetectorName
	cc.detector, err = vision.FromDependencies(deps, countConf.DetectorName)
	if err != nil {
		return errors.Wrapf(err, "unable to get vision service %v for count classifier", countConf.DetectorName)
	}
	// put everything in lower case
	labels := map[string]float64{}
	for l, c := range countConf.ChosenLabels {
		labels[strings.ToLower(l)] = c
	}
	cc.labels = labels
	cc.thresholds = NewThresholds(countConf.CountThresholds)
	return nil
}

func (cc *countcls) count(dets []objdet.Detection) (string, []objdet.Detection) {
	// get the number of boxes with the right label and confidences
	count := 0
	outDets := []objdet.Detection{}
	for _, d := range dets {
		label := strings.ToLower(d.Label())
		if conf, ok := cc.labels[label]; ok {
			if d.Score() >= conf {
				count++
				outDets = append(outDets, d)
			}
		}
	}
	// associated the number with the right label
	for _, thresh := range cc.thresholds {
		if count <= thresh.UpperBound {
			return thresh.Label, outDets
		}
	}
	return OverflowLabel, outDets
}

// Detections just calls the underlying detector
func (cc *countcls) DetectionsFromCamera(
	ctx context.Context,
	cameraName string,
	extra map[string]interface{},
) ([]objdet.Detection, error) {
	return cc.detector.DetectionsFromCamera(ctx, cameraName, extra)
}

// Detections just calls the underlying detector
func (cc *countcls) Detections(ctx context.Context, img image.Image, extra map[string]interface{}) ([]objdet.Detection, error) {
	return cc.detector.Detections(ctx, img, extra)
}

// ClassificationsFromCamera calls DetectionsFromCamera on the underlying service and counts valid boxes.
func (cc *countcls) ClassificationsFromCamera(
	ctx context.Context,
	cameraName string,
	n int,
	extra map[string]interface{},
) (classification.Classifications, error) {
	cls := []classification.Classification{}
	dets, err := cc.detector.DetectionsFromCamera(ctx, cameraName, extra)
	if err != nil {
		return nil, errors.Wrapf(err, "error from underlying detector %s", cc.detName)
	}
	label, _ := cc.count(dets)
	c := classification.NewClassification(1.0, label)
	cls = append(cls, c)
	return classification.Classifications(cls), nil
}

// Classifications calls Detections on the underlying service and counts valid boxes.
func (cc *countcls) Classifications(ctx context.Context, img image.Image,
	n int, extra map[string]interface{},
) (classification.Classifications, error) {
	cls := []classification.Classification{}
	dets, err := cc.detector.Detections(ctx, img, extra)
	if err != nil {
		return nil, errors.Wrapf(err, "error from underlying vision model %s", cc.detName)
	}
	label, _ := cc.count(dets)
	c := classification.NewClassification(1.0, label)
	cls = append(cls, c)
	return classification.Classifications(cls), nil
}

func (cc *countcls) GetObjectPointClouds(
	ctx context.Context,
	cameraName string,
	extra map[string]interface{},
) ([]*viz.Object, error) {
	return nil, errUnimplemented
}

// GetProperties returns the properties
func (cc *countcls) GetProperties(ctx context.Context, extra map[string]interface{}) (*vision.Properties, error) {
	return &cc.properties, nil
}

// CaptureAllFromCamera calls the underlying detector's method and adds a classification
func (cc *countcls) CaptureAllFromCamera(
	ctx context.Context,
	cameraName string,
	opt viscapture.CaptureOptions,
	extra map[string]interface{},
) (viscapture.VisCapture, error) {
	opt.ReturnDetections = true
	visCapture, err := cc.detector.CaptureAllFromCamera(ctx, cameraName, opt, extra)
	if err != nil {
		return visCapture, errors.Wrapf(err, "error from underlying detector %s", cc.detName)
	}
	label, dets := cc.count(visCapture.Detections)
	cls := []classification.Classification{}
	c := classification.NewClassification(1.0, label)
	cls = append(cls, c)
	visCapture.Classifications = classification.Classifications(cls)
	visCapture.Detections = dets
	return visCapture, nil
}

// Close does nothing
func (cc *countcls) Close(ctx context.Context) error {
	return nil
}

// DoCommand implements nothing
func (cc *countcls) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}
