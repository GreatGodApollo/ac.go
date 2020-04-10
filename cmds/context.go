package cmds

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"io"
)

// A Context is passed to a CommandRunFunc. It contains the information needed for a command to execute.
type CommandContext struct {
	// The connection to Discord.
	Session *discordgo.Session

	// The event that fired the CommandHandler.
	Event *discordgo.MessageCreate

	// The CommandManager that handled this command.
	Manager *Manager

	// The custom args struct for this command
	Args interface{}

	// The Message that fired this event.
	Message *discordgo.Message

	// The User that fired this event.
	User *discordgo.User

	// The Channel the event was fired in.
	Channel *discordgo.Channel

	// The guild the Channel belongs to.
	Guild *discordgo.Guild

	// The User's guild member.
	Member *discordgo.Member
}

// Reply sends a message to the channel a CommandContext was initiated for.
func (ctx *CommandContext) Reply(message string) (*discordgo.Message, error) {
	return ctx.Session.ChannelMessageSend(ctx.Channel.ID, message)
}

// ReplyEmbed sends an embed to the channel a CommandContext was initiated for.
func (ctx *CommandContext) ReplyEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID, embed)
}

// ReplyFile sends a file to the channel a CommandContext was initiated for.
func (ctx *CommandContext) ReplyFile(filename string, file io.Reader) (*discordgo.Message, error) {
	return ctx.Session.ChannelFileSend(ctx.Channel.ID, filename, file)
}

// PurgeMessages purges 'x' number of messages from the Channel a CommandContext was initiated for.
func (ctx *CommandContext) PurgeMessages(num int) error {
	if num >= 1 && num <= 100 {
		msgs, err := ctx.Session.ChannelMessages(ctx.Channel.ID, num, "", "", "")
		if err != nil {
			return err
		}
		var ids []string
		for _, msg := range msgs {
			ids = append(ids, msg.ID)
		}
		return ctx.Session.ChannelMessagesBulkDelete(ctx.Channel.ID, ids)
	} else if num > 1 && num > 100 {
		return errors.New("too many messages")
	} else if num == 0 {
		return errors.New("must supply a number")
	}
	return nil
}
