package main

import (
	"fmt"
	"github.com/bal3000/BalStreamer.Caster/app"
	"github.com/bal3000/BalStreamer.Caster/infrastructure"
	"os"
)

var config infrastructure.Configuration

func init() {
	config = infrastructure.ReadConfig()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "startup error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	streamer, err := infrastructure.NewChromecastStreamer()
	if err != nil {
		return err
	}
	defer streamer.CloseAndCleanUp()

	rabbitMq, err := infrastructure.NewRabbitMQConnection(&config)
	if err != nil {
		return err
	}
	defer rabbitMq.CloseChannel()

	server := app.NewServer(rabbitMq, streamer)
	err = server.Run()
	if err != nil {
		return err
	}

	return nil
}
