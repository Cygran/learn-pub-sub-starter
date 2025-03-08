package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitSrv string = "amqp://guest:guest@localhost:5672/"
	fmt.Println("Starting Peril server...")

	//connect to rabbit amqp server
	connection, err := amqp.Dial(rabbitSrv)
	if err != nil {
		log.Fatalf("Unable to establish connection: %s", err)
		return
	}
	defer connection.Close()
	fmt.Println("Connection to Server successful")

	//create channel
	ch1, err := connection.Channel()
	if err != nil {
		log.Fatalf("Unable to create Channel on connection: %s", err)
		return
	}

	//publish message to the exchange
	err = pubsub.PublishJSON(ch1, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})
	if err != nil {
		log.Fatalf("Failed to publish message: %s", err)
		return
	}

	//handle closure from ctrl + c gracefully
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("Recieved interrupt. Closing connection")
	connection.Close()
}
