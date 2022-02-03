package Piper

import "time"

type EventTypeName string

// Event Marker interface for events
type Event interface {
	isEvent()
}

type EventListener interface {
	onEvent(e Event) error
}

type EventPublisher struct {
	listeners []EventListener
}

func (b *EventPublisher) AddListener(l EventListener) {

	for _, listener := range b.listeners {
		if listener == l {
			return
		}
	}

	b.listeners = append(b.listeners, l)
}

func (b *EventPublisher) Publish(e Event) error {

	for _, listener := range b.listeners {
		if err := listener.onEvent(e); err != nil {
			return err
		}
	}

	return nil
}

type MultiEventListener struct {
	Listeners []EventListener
}

func (m MultiEventListener) onEvent(e Event) error {
	for _, listener := range m.Listeners {
		if err := listener.onEvent(e); err != nil {
			return err
		}
	}

	return nil
}

// EVENTS

//goland:noinspection GoNameStartsWithPackageName
type PipelineStartedEvent struct {
	StartedAt     time.Time
	StageNames    []StageName
	Configuration Configuration
}

func (e PipelineStartedEvent) isEvent() {}

//goland:noinspection GoNameStartsWithPackageName
type PipelineEndedEvent struct {
	Configuration Configuration
	StartedAt     time.Time
	EndedAt       time.Time
	StageNames    []StageName
	Output        interface{}
	Error         error
	LastStageName StageName
}

func (e PipelineEndedEvent) isEvent() {}

type StageStartedEvent struct {
	Configuration Configuration
	StartedAt     time.Time
	StageName     StageName
	Input         interface{}
}

func (e StageStartedEvent) isEvent() {}

type StageEndedEvent struct {
	Configuration Configuration
	StartedAt     time.Time
	EndedAt       time.Time
	StageName     StageName
	Input         interface{}
}

func (e StageEndedEvent) isEvent() {}
