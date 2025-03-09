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

	//create message publishing channel
	publishCh, err := connection.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	//User Login
	userName, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Unable to fecth username: %s", err)
	}

	//create Game State
	gs := gamelogic.NewGameState(userName)

	//create, bind and subscibe to pause queue
	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilDirect, routing.PauseKey+"."+userName, routing.PauseKey, pubsub.SimpleQueueTransient, handlerPause(gs))
	if err != nil {
		log.Fatalf("Unable to subscribe to pause queue: %s", err)
	}

	//create, bind and subscribe to move queue
	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+userName, routing.ArmyMovesPrefix+".*", pubsub.SimpleQueueTransient, handlerMove(gs))
	if err != nil {
		log.Fatalf("Unable to subscribe to moves queue: %s", err)
	}

	//Game Loop REPL
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
			mv, err := gs.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilTopic,
				routing.ArmyMovesPrefix+"."+mv.Player.Username,
				mv,
			)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				continue
			}
			fmt.Printf("Moved %v units to %s\n", len(mv.Units), mv.ToLocation)

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
