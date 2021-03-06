package main

import (
	"github.com/jwillp/Piper"
)

type LintingOutput struct {
	// Represents the linting errors encountered during the processing.
	Errors []error
}

type BuildStage struct {
}

func (s BuildStage) Name() Piper.StageName {
	return "Build"
}

func (s BuildStage) Run(input interface{}) (interface{}, error) {
	return nil, nil
}

func main() {
	p := Piper.NewPipeline().WithStages(
		Piper.NewStage("Build", func(input interface{}) (interface{}, error) {
			return nil, nil
		}),
		Piper.NewStage("UnitTest", func(input interface{}) (interface{}, error) {
			return nil, nil
		}),
		Piper.NewStage("IntegrationTest", func(input interface{}) (interface{}, error) {
			return nil, nil
		}),
		Piper.NewStage("DeployStaging", func(input interface{}) (interface{}, error) {
			return nil, nil
		}),
		Piper.NewStage("TestStaging", func(input interface{}) (interface{}, error) {
			return nil, nil
		}),
		Piper.NewStage("DeployProduction", func(input interface{}) (interface{}, error) {
			return nil, nil
		}),
		Piper.NewPipelineStage("Do something else", Piper.NewPipeline().WithStages(
			Piper.NewStage("Build", func(input interface{}) (interface{}, error) {
				return nil, nil
			}),
		)),
	).WithEventListeners(
		Piper.LoggerEventListener{},
	)

	if _, err := p.Run(nil); err != nil {
		panic(err)
	}
}
