package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fvdveen/fuzzy"
)

func main() {
	token := flag.String("token", "", "The discord bot token")
	prefix := flag.String("prefix", "$", "The bots prefix")
	flag.Parse()

	bot, err := fuzzy.New(
		fuzzy.WithToken(*token),
		fuzzy.WithPrefix(*prefix),
	)
	if err != nil {
		log.Fatal(err)
	}

	_ = bot.RegisterCommand(
		fuzzy.HelpCommand("ping-pong bot", "An example ping-pong bot"),
		fuzzy.NewCommand("ping", "sends pong", func(ctx fuzzy.Context) {
			ctx.SendMessage("pong")
		}),
		fuzzy.NewCommand("pong", "sends ping", func(ctx fuzzy.Context) {
			ctx.SendMessage("ping")
		}),
	)

	err = bot.Open()
	if err != nil {
		log.Fatal(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	err = bot.Close()
	if err != nil {
		log.Fatal(err)
	}
}
