package countclassifier

import (
	"context"

	"github.com/pkg/errors"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/vision"
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
	var err error
	cc.detName = conf.DetectorName
	cc.detector, err = vision.FromDependencies(deps, conf.DetectorName)
	if err != nil {
		return errors.Wrapf(err, "unable to get vision service %v for count classifier", conf.DetectorName)
	}
	return nil
}
