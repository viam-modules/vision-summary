// Package main is a module which serves the countclassifier vision service
package main

import (
	"context"

	"go.viam.com/rdk/services/vision"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/utils"

	"github.com/viam-modules/vision-summary/countclassifier"
)

func main() {
	utils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs("count-classifier"))
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) (err error) {
	myMod, err := module.NewModuleFromArgs(ctx)
	if err != nil {
		return err
	}

	err = myMod.AddModelFromRegistry(ctx, vision.API, countclassifier.Model)
	if err != nil {
		return err
	}

	err = myMod.Start(ctx)
	defer myMod.Close(ctx)
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}
