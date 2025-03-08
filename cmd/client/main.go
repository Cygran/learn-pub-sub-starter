package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitSrv string = "amqp://guest:guest@localhost:5672/"
	fmt.Println("Starting Peril client...")
	//connect to rabbit amqp server
	connection, err := amqp.Dial(rabbitSrv)
	if err != nil {
		log.Fatalf("Unable to establish connection: %s", err)
		return
	}
	defer connection.Close()
	fmt.Println("Connection to Server successful")

	//User Login
	userName, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Unable to fecth username: %s", err)
	}

	//Create and bind to pause queue
	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, userName)
	_, queue, err := pubsub.DeclareAndBind(connection, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.QueueTypeTransient)
	if err != nil {
		log.Fatalf("Unable to establish queue on exchange: %s", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	//Game loop
	gs := gamelogic.NewGameState(userName)
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "spawn":
			err := gs.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "move":
			_, err = gs.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("unknown command")
		}
	}
}
