package infrastructure

import (
	"encoding/json"
	"log"
	"os"
)

// Configuration is the model to store the config data
type Configuration struct {
	RabbitURL    string `json:"rabbitUrl"`
	ExchangeName string `json:"exchangeName"`
	QueueName    string `json:"queueName"`
	Durable      bool   `json:"durable"`
}

// ReadConfig reads the json file and populates the Configuration struct
func ReadConfig() Configuration {
	file, _ := os.Open("./config.json")
	defer file.Close()
	log.Println("Loading in config file")

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)

	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Loaded in config file")
	return configuration
}
