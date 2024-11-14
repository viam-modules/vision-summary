// Package main is a module which serves the countclassifier vision service
package main

import (
	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/services/vision"

	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"

	"github.com/viam-modules/vision-summary/countclassifier"
	"github.com/viam-modules/vision-summary/countsensor"
)

func main() {
	module.ModularMain(resource.APIModel{vision.API, countclassifier.Model},
	                   resource.APIModel{sensor.API, countsensor.Model})
}
