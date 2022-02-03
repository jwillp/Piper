package Piper

type ClosureStageClosure func(input StageInput) (StageOutput, error)

// ClosureStage allows defining a stage using a closure.
type ClosureStage struct {
	name    StageName
	closure ClosureStageClosure
}

func NewStage(n StageName, c ClosureStageClosure) ClosureStage {
	return ClosureStage{n, c}
}

func (s ClosureStage) Name() StageName {
	return s.name
}

func (s ClosureStage) Run(input StageInput) (StageOutput, error) {
	return s.closure(input)
}

// PipelineStage is a stage allowing to forward the input to another pipeline.
type PipelineStage struct {
	name     StageName
	pipeline *Pipeline
}

func NewPipelineStage(n StageName, p *Pipeline) PipelineStage {
	return PipelineStage{name: n, pipeline: p}
}

func (p PipelineStage) Name() StageName {
	return p.name
}

func (p PipelineStage) Run(input StageInput) (StageOutput, error) {
	return p.pipeline.Run(input)
}
