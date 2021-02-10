package models

import (
	"encoding/json"
	"time"
)

const (
	PlayStream string = "PlayStreamEvent"
	StopStream string = "StopStreamEvent"
)

// ChromecastEventMessage interface for transforming messages
type ChromecastEventMessage interface {
	TransformMessage() ([]byte, error)
}

type StreamEvent interface {
	GetType() string
	GetStream() string
	GetChromecast() string
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

func (event *StreamToChromecastEvent) GetType() string {
	return PlayStream
}

func (event *StreamToChromecastEvent) GetStream() string {
	return event.Stream
}

func (event *StreamToChromecastEvent) GetChromecast() string {
	return event.ChromeCastToStream
}

func (event *StopPlayingStreamEvent) GetType() string {
	return StopStream
}

func (event *StopPlayingStreamEvent) GetStream() string {
	return ""
}

func (event *StopPlayingStreamEvent) GetChromecast() string {
	return event.ChromeCastToStop
}

// ChromecastEvent event when a chromecast is found
type ChromecastFoundEvent struct {
	Chromecast string `json:"chromecast"`
}

// ChromecastEvent event when a chromecast is found
type ChromecastLostEvent struct {
	Chromecast string `json:"chromecast"`
}

// TransformMessage transforms the message to a masstransit one and then turns into JSON
func (message *ChromecastFoundEvent) TransformMessage() ([]byte, error) {
	return json.Marshal(message)
}

// TransformMessage transforms the message to a masstransit one and then turns into JSON
func (message *ChromecastLostEvent) TransformMessage() ([]byte, error) {
	return json.Marshal(message)
}
