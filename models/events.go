package models

import (
	"encoding/json"
	"time"
)

// EventMessage interface for transforming messages to masstransit ones
type EventMessage interface {
	TransformMessage() ([]byte, error)
}

// StreamToChromecastEvent the send to chromecast event
type StreamToChromecastEvent struct {
	ChromeCastToStream string    `json:"chromeCastToStream"`
	Stream             string    `json:"stream"`
	StreamDate         time.Time `json:"streamDate"`
}

// StopPlayingStreamEvent the stop cast event
type StopPlayingStreamEvent struct {
	ChromeCastToStop string    `json:"chromeCastToStop"`
	StopDateTime     time.Time `json:"stopDateTime"`
}

// ChromecastEvent event when a chromecast is found
type ChromecastFoundEvent struct {
	Chromecast interface{} `json:"chromecast"`
}

// ChromecastEvent event when a chromecast is found
type ChromecastLostEvent struct {
	Chromecast interface{} `json:"chromecast"`
}

// TransformMessage transforms the message to a masstransit one and then turns into JSON
func (message *StreamToChromecastEvent) TransformMessage() ([]byte, error) {
	return json.Marshal(message)
}

// TransformMessage transforms the message to a masstransit one and then turns into JSON
func (message *StopPlayingStreamEvent) TransformMessage() ([]byte, error) {
	return json.Marshal(message)
}
