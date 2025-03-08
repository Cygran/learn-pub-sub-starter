package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitSrv string = "amqp://guest:guest@localhost:5672/"
	fmt.Println("Starting Peril server...")
	connection, err := amqp.Dial(rabbitSrv)
	if err != nil {
		log.Fatalf("Unable to establish connection: %s", err)
		return
	}
	defer connection.Close()
	fmt.Println("Connection to Server successful")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("Recieved interrupt. Closing connection")
	connection.Close()
}
