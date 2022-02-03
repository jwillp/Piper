package Piper

import (
	"fmt"
	"github.com/benbjohnson/clock"
)

// Configuration Configuration of an Pipeline.
type Configuration interface {
	isConfiguration()
}

type EmptyConfiguration struct {
}

func (c EmptyConfiguration) isConfiguration() {}

// StageName Name of a Stage
type StageName string

// Stage Represents a stage in a pipeline.
type Stage interface {
	Name() StageName

	Run(i interface{}) (interface{}, error)
}

// Pipeline Service responsible for the execution of a Pipeline.
type Pipeline struct {
	Configuration  Configuration
	Stages         []Stage
	EventPublisher *EventPublisher
	Clock          clock.Clock
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		EmptyConfiguration{},
		[]Stage{},
		&EventPublisher{
			listeners: []EventListener{},
		},
		clock.New(),
	}
}

func (p *Pipeline) WithConfiguration(c Configuration) *Pipeline {
	p.Configuration = c

	return p
}

func (p *Pipeline) UsingEventPublisher(ep *EventPublisher) *Pipeline {
	p.EventPublisher = ep

	return p
}

func (p *Pipeline) WithEventListeners(l ...EventListener) *Pipeline {
	for _, listener := range l {
		p.EventPublisher.AddListener(listener)
	}

	return p
}

func (p *Pipeline) WithStages(s ...Stage) *Pipeline {
	p.Stages = append(p.Stages, s...)

	return p
}

func (p *Pipeline) Run(input interface{}) (o interface{}, err error) {

	var output interface{}

	pipelineStartedAt := p.Clock.Now()

	if err := p.notifyListeners(PipelineStartedEvent{
		StartedAt:     pipelineStartedAt,
		StageNames:    p.stageNames(),
		Configuration: p.Configuration,
	}); err != nil {
		return nil, err
	}

	defer func(p *Pipeline) {
		if dispatchError := p.notifyListeners(PipelineEndedEvent{
			StartedAt:     pipelineStartedAt,
			EndedAt:       p.Clock.Now(),
			StageNames:    p.stageNames(),
			Configuration: p.Configuration,
		}); dispatchError != nil {
			err = dispatchError
		}
	}(p)

	input = nil
	for _, stage := range p.Stages {
		stageStartedAt := p.Clock.Now()

		if err := p.notifyListeners(StageStartedEvent{
			Configuration: p.Configuration,
			StartedAt:     pipelineStartedAt,
			StageName:     stage.Name(),
			Input:         input,
		}); err != nil {
			return nil, err
		}

		//goland:noinspection GoDeferInLoop
		defer func(p *Pipeline) {
			if dispatchError := p.notifyListeners(StageEndedEvent{
				Configuration: p.Configuration,
				StartedAt:     stageStartedAt,
				EndedAt:       p.Clock.Now(),
				StageName:     stage.Name(),
				Input:         input,
				Output:        output,
				Error:         err,
			}); dispatchError != nil {
				err = dispatchError
			}
		}(p)

		if out, err := stage.Run(input); err != nil {
			return nil, err
		} else {
			input = out
			output = out
		}
	}

	return output, nil
}

func (p *Pipeline) notifyListeners(e Event) error {
	return p.EventPublisher.Publish(e)
}

func (p *Pipeline) stageNames() []StageName {
	var names []StageName
	for _, s := range p.Stages {
		names = append(names, s.Name())
	}

	return names
}

type LoggerEventListener struct {
}

func (l LoggerEventListener) OnEvent(e Event) error {
	fmt.Printf("%T", e)
	fmt.Println()

	return nil
}
