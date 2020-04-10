package main

import (
	"github.com/GreatGodApollo/acgo/cmds"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var logger = &log.Logger{}

func main() {
	log.Println("Starting testBot")

	client, err := discordgo.New("Bot " + "NjUwMTQ1OTY0MDQwMTI2NDcx.XgxFgQ.IlWIFSBHI2WHMFQlqeWx9_IzpiQ")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	manager := cmds.NewManager(logger, []string{}, []string{}, CommandErrorFunc)
	manager.AddPrefix("t.")

	manager.RegisterCommand(cmds.DefaultHelp)

	client.AddHandler(manager.Handle)

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
