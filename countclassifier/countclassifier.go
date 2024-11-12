package countclassifier

import (
	"context"
	"image"

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
	ModelName = "countclassifier"
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
	CountThresholds map[uint]string    `json:"count_thresholds"`
}

// Validate validates the config and returns implicit dependencies,
// this Validate checks if the camera and detector exist for the module's vision model.
func (cfg *Config) Validate(path string) ([]string, error) {
	if cfg.DetectorName == "" {
		return nil, errors.New("attribute detector_name cannot be left blank")
	}
	return []string{cfg.DetectorName}, nil
}

type countcls struct {
	resource.Named
	logger     logging.Logger
	properties vision.Properties
	detName    string
	detector   vision.Service
	labels     map[string]float64
	thresholds map[uint]string
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
			DetectionSupported:      false,
			ObjectPCDsSupported:     false,
		},
	}

	if err := cc.Reconfigure(ctx, deps, conf); err != nil {
		return nil, err
	}
	return cc, nil
}

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
	return nil
}

func (cc *countcls) count(dets []objdet.Detection) string {
	return ""
}

func (cc *countcls) DetectionsFromCamera(
	ctx context.Context,
	cameraName string,
	extra map[string]interface{},
) ([]objdet.Detection, error) {
	return nil, errUnimplemented
}

func (cc *countcls) Detections(ctx context.Context, img image.Image, extra map[string]interface{}) ([]objdet.Detection, error) {
	return nil, errUnimplemented
}

func (cc *countcls) ClassificationsFromCamera(
	ctx context.Context,
	cameraName string,
	n int,
	extra map[string]interface{},
) (classification.Classifications, error) {
	return nil, nil
}

func (cc *countcls) Classifications(ctx context.Context, img image.Image,
	n int, extra map[string]interface{},
) (classification.Classifications, error) {
	return nil, nil
}

func (cc *countcls) GetObjectPointClouds(
	ctx context.Context,
	cameraName string,
	extra map[string]interface{},
) ([]*viz.Object, error) {
	return nil, errUnimplemented
}

func (cc *countcls) GetProperties(ctx context.Context, extra map[string]interface{}) (*vision.Properties, error) {
	return &cc.properties, nil
}

func (cc *countcls) CaptureAllFromCamera(
	ctx context.Context,
	cameraName string,
	opt viscapture.CaptureOptions,
	extra map[string]interface{},
) (viscapture.VisCapture, error) {
	return viscapture.VisCapture{}, nil
}

func (cc *countcls) Close(ctx context.Context) error {
	return nil
}

func (cc *countcls) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}
