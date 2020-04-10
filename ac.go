package main

import (
	"github.com/GreatGodApollo/acgo/cmds"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var logger = &logrus.Logger{}

func main() {
	log.Println("Starting testBot")

	client, err := discordgo.New("Bot " + "NjUwMTQ1OTY0MDQwMTI2NDcx.XpD3cA.fsuD95W6txaipp3IcMS_qzJqLSQ")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	manager := cmds.NewManager(logger, true, CommandErrorFunc)
	manager.AddPrefix("t.")

	manager.AddCommand(cmds.DefaultHelp(0x532c60))

	client.AddHandler(manager.CommandHandler)

	err = client.Open()
	if err != nil {
		log.Fatal(err)
		return
	}

	manager.AddPrefix("<@!" + client.State.User.ID + "> ")
	manager.AddPrefix("<@!" + client.State.User.ID + ">")

	// Wait until a term signal is received
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close after term signal
	_ = client.Close()
}

func CommandErrorFunc(cmdm *cmds.Manager, ctx cmds.Context, err error) {
	cmdm.Logger.Println("Err: ", err.Error())
	_, err2 := ctx.Reply(":x: Error: `" + err.Error() + "` :x:")
	if err2 != nil {
		cmdm.Logger.Println("Err: ", err2.Error())
	}
}
