package cmds

import (
	"github.com/bwmarrin/discordgo"
	"io"
)

type Context struct {
	Session *discordgo.Session

	Event *discordgo.MessageCreate

	Manager *Manager

	Args []string

	Guild *discordgo.Guild

	Channel *discordgo.Channel

	Message *discordgo.Message

	Member *discordgo.Member

	User *discordgo.User
}

func (ctx Context) Reply(message string) (*discordgo.Message, error) {
	return ctx.Session.ChannelMessageSend(ctx.Channel.ID, message)
}

func (ctx *Context) ReplyEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID, embed)
}

func (ctx *Context) ReplyFile(filename string, file io.Reader) (*discordgo.Message, error) {
	return ctx.Session.ChannelFileSend(ctx.Channel.ID, filename, file)
}
