package app

import "github.com/bal3000/BalStreamer.Caster/infrastructure"

type Server struct {
	RabbitMQ           infrastructure.RabbitMQ
	ChromecastStreamer infrastructure.Streamer
}

func NewServer() *Server {

	return &Server{}
}
