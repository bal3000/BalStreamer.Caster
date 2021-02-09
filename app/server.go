package app

import (
	"fmt"
	"github.com/bal3000/BalStreamer.Caster/infrastructure"
	"github.com/bal3000/BalStreamer.Caster/models"
	"github.com/streadway/amqp"
)

const routingKey string = "chromecast-key"

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

	// find chromecast to send to - might have to turn into a channel

	// send to correct chromecast streamer method

	return true
}
