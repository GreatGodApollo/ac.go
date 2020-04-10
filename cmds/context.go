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

func NewContext() Context {
	return Context{}
}

func (ctx Context) Reply(message string) (*discordgo.Message, error) {
	return ctx.Session.ChannelMessageSend(ctx.Channel.ID, message)
}

func (ctx Context) ReplyEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID, embed)
}

func (ctx Context) ReplyFile(filename string, file io.Reader) (*discordgo.Message, error) {
	return ctx.Session.ChannelFileSend(ctx.Channel.ID, filename, file)
}

func (ctx Context) SetSession(session *discordgo.Session) Context {
	ctx.Session = session
	return ctx
}

func (ctx Context) SetEvent(event *discordgo.MessageCreate) Context {
	ctx.Event = event
	return ctx
}

func (ctx Context) SetManager(manager *Manager) Context {
	ctx.Manager = manager
	return ctx
}

func (ctx Context) SetArgs(args []string) Context {
	ctx.Args = args
	return ctx
}

func (ctx Context) SetGuild(guild *discordgo.Guild) Context {
	ctx.Guild = guild
	return ctx
}

func (ctx Context) SetChannel(channel *discordgo.Channel) Context {
	ctx.Channel = channel
	return ctx
}

func (ctx Context) SetMessage(message *discordgo.Message) Context {
	ctx.Message =  message
	return ctx
}

func (ctx Context) SetMember(member *discordgo.Member) Context {
	ctx.Member = member
	return ctx
}

func (ctx Context) SetUser(user *discordgo.User) Context {
	ctx.User = user
	return ctx
}
