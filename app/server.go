package app

import (
	"encoding/json"
	"fmt"
	"github.com/bal3000/BalStreamer.Caster/infrastructure"
	"github.com/bal3000/BalStreamer.Caster/models"
	"github.com/streadway/amqp"
	"log"
)

const (
	routingKey         string = "chromecast-key"
	streamToChromecast string = "StreamToChromecastEvent"
	stopChromecast     string = "StopPlayingStreamEvent"
)

type Server struct {
	RabbitMQ           infrastructure.RabbitMQ
	ChromecastStreamer infrastructure.Streamer
}

func NewServer(rabbit infrastructure.RabbitMQ, streamer infrastructure.Streamer) *Server {
	return &Server{RabbitMQ: rabbit, ChromecastStreamer: streamer}
}

func (s *Server) Run() error {
	// Start listening for cast events
	err := s.RabbitMQ.StartConsumer(routingKey, processMessages, 2)
	if err != nil {
		return err
	}

	// Find chromecasts
	itemAdded := make(chan string)
	itemRemoved := make(chan string)
	if err := s.ChromecastStreamer.DiscoverChromecasts(itemAdded, itemRemoved); err != nil {
		return err
	}

	// Send events out once one is found or removed
	for item := range itemAdded {
		message := &models.ChromecastFoundEvent{Chromecast: item}
		if err := s.RabbitMQ.SendMessage(routingKey, message); err != nil {
			return err
		}
	}

	for item := range itemRemoved {
		message := &models.ChromecastLostEvent{Chromecast: item}
		if err := s.RabbitMQ.SendMessage(routingKey, message); err != nil {
			return err
		}
	}

	return nil
}

func processMessages(d amqp.Delivery) bool {
	fmt.Printf("processing message: %s", string(d.Body))

	// find if event is start or stop
	event, err := getMessage(d)
	if err != nil {
		log.Fatalln(err)
		return false
	}

	// If its not a stream event ignore it and carry on
	if event == nil {
		return true
	}

	// find chromecast to send to - might have to turn into a channel

	// send to correct chromecast streamer method

	return true
}

func getMessage(d amqp.Delivery) (models.StreamEvent, error) {
	if d.Type == streamToChromecast {
		event := new(models.StreamToChromecastEvent)
		if err := json.Unmarshal(d.Body, event); err != nil {
			return nil, err
		}
		return event, nil
	} else if d.Type == stopChromecast {
		event := new(models.StopPlayingStreamEvent)
		if err := json.Unmarshal(d.Body, event); err != nil {
			return nil, err
		}
		return event, nil
	}

	return nil, nil
}
