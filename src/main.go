package main

import (
	"peonbot/params"
	"peonbot/peonbot"
	"peonbot/verbose"

	"log"
)

func main() {
	log.Printf("Starting up...\n")

	/* Scoop up command lines args */
	p, err := params.New()
	if err != nil {
		panic(err)
	}
	/* Set verbose printer */
	verbose.SetPrinter(p.Args.Verbose())

	/* Connect bot to battle.net */
	bot := peonbot.New(p.Token(), p.Config.Blist(), p.Config.Greetings(),
		p.Config.Pusers())

	if err := bot.Start(); err != nil {
		panic(err)
	}
	defer bot.Conn.Close()

	/* Listen for user input from stdin */
	go bot.ListenStdin()

	/* Listen for responses from websocket */
	go bot.ListenWebsocket()

	/* Event loop */
event_loop:
	for {
		select {
		case event := <-bot.Chbnt():
			if err := bot.HandleEvent(event); err != nil {
				bot.Vprintf("Got error from handling event: %v\n", err)
			}
		case err := <-bot.Cherr():
			bot.Vprintf("Got error reading from websocket: %v\n", err)
			break event_loop
		case msg := <-bot.Chsin():
			bot.SendMessage(msg)
		}

	}

	log.Printf("Event loop broken. Shutting down...\n")
}
