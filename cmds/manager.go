package cmds

import (
	"github.com/GreatGodApollo/acgo/permissions"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

type Manager struct {
	// The list of prefixes the bot should respond to
	Prefixes []string

	// The list of IDs the bot should consider to be an owner
	Owners []string

	// Should the Manager ignore bots?
	IgnoreBots bool

	// The logger for the bot
	Logger log.Logger

	// The map of Commands in the Manager
	Commands map[string]*Command

	// The function to run when the manager errors
	ErrorFunc ManagerOnError
}

type ManagerOnError func(cmdm *Manager, ctx Context, err error)

func (cmdm *Manager) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot && cmdm.IgnoreBots {
		return
	}

	var prefix string
	var contains bool
	var err error
	var ctx Context
	for i := 0; i < len(cmdm.Prefixes); i++ {
		prefix = cmdm.Prefixes[i]
		if strings.HasPrefix(m.Content, prefix) {
			contains = true
			break
		}
	}

	if !contains {
		return
	}

	cmd := strings.Split(strings.TrimPrefix(m.Content, prefix), " ")

	ctx.Manager = cmdm
	ctx.Message = m.Message
	ctx.Member = m.Member
	ctx.Session = s
	ctx.User = m.Author
	ctx.Guild, _ = s.Guild(m.GuildID)
	ctx.Args = cmd[1:]

	ctx.Channel, _ = s.Channel(m.ChannelID)

	if command := cmdm.GetCommand(cmd[0]); command != nil {
		var inDm bool
		if ctx.Channel.Type == discordgo.ChannelTypeDM {
			inDm = true
		}

		// Check UserPermissions
		if command.Type != Direct && !inDm && !permissions.Check(s, m.GuildID, m.Author.ID, command.UserPerms) {
			if permissions.Check(s, m.GuildID, s.State.User.ID, permissions.PermissionMessagesEmbedLinks) {
				embed := &discordgo.MessageEmbed{
					Title:       "Insufficient Permissions!",
					Description: "You don't have the required permissions to run this command!",
					Color:       0xff0000,
				}

				if !command.Hidden {
					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
				}
			} else {
				if !command.Hidden {
					_, err = s.ChannelMessageSend(m.ChannelID, ":x: You don't have the correct permissions to run this command! :x:")
				}
			}
			if err != nil {
				cmdm.ErrorFunc(cmdm, ctx, err)
			}
			cmdm.Logger.Printf("P: FALSE C: %s[%s] U: %s#%s[%s] M: %s", ctx.Channel.Name, ctx.Channel.ID, ctx.User.Username, ctx.User.Discriminator, ctx.User.ID, m.Content)
			return
		}

		// Check BotPermissions
		if command.Type != Direct && !inDm && !permissions.Check(s, m.GuildID, s.State.User.ID, command.BotPerms) {
			if permissions.Check(s, m.GuildID, s.State.User.ID, permissions.PermissionMessagesEmbedLinks) {
				embed := &discordgo.MessageEmbed{
					Title:       "Insufficient Permissions!",
					Description: "I don't have the correct permissions to run this command!",
					Color:       0xff0000,
				}

				if !command.Hidden {
					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
				}
			} else {
				if !command.Hidden {
					_, err = s.ChannelMessageSend(m.ChannelID, ":x: I don't have the correct permissions to run this command! :x:")
				}
			}

			if err != nil {
				cmdm.ErrorFunc(cmdm, ctx, err)
			}
			cmdm.Logger.Printf("P: FALSE C: %s[%s] U: %s#%s[%s] M: %s", ctx.Channel.Name, ctx.Channel.ID, ctx.User.Username, ctx.User.Discriminator, ctx.User.ID, m.Content)
			return
		}

		// Check if it's the right channel type
		if inDm && command.Type == Guild {
			embed := &discordgo.MessageEmbed{
				Title:       "Invalid Channel!",
				Description: "You cannot run this command in a private message.",
				Color:       0xff0000,
			}

			if !command.Hidden {
				_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
			}

			if err != nil {
				cmdm.ErrorFunc(cmdm, ctx, err)
			}
			cmdm.Logger.Printf("P: FALSE C: %s[%s] U: %s#%s[%s] M: %s", ctx.Channel.Name, ctx.Channel.ID, ctx.User.Username, ctx.User.Discriminator, ctx.User.ID, m.Content)
			return
		} else if !inDm && command.Type == Direct {
			embed := &discordgo.MessageEmbed{
				Title:       "Invalid Channel!",
				Description: "You cannot run this command in a guild.",
				Color:       0xff0000,
			}

			if !command.Hidden {
				_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
			}

			if err != nil {
				cmdm.ErrorFunc(cmdm, ctx, err)
			}
			cmdm.Logger.Printf("P: FALSE C: %s[%s] U: %s#%s[%s] M: %s", ctx.Channel.Name, ctx.Channel.ID, ctx.User.Username, ctx.User.Discriminator, ctx.User.ID, m.Content)
			return
		}

		// Is it an owner command but you don't have owner?
		if command.OwnerOnly && !cmdm.IsOwner(m.Author.ID) {
			embed := &discordgo.MessageEmbed{
				Title:       "You can't run that command!",
				Description: "Sorry, only bot owners can run that command",
				Color:       0xff0000,
			}

			if !command.Hidden {
				_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
			}

			if err != nil {
				cmdm.ErrorFunc(cmdm, ctx, err)
			}
			cmdm.Logger.Printf("P: FALSE C: %s[%s] U: %s#%s[%s] M: %s", ctx.Channel.Name, ctx.Channel.ID, ctx.User.Username, ctx.User.Discriminator, ctx.User.ID, m.Content)
			return
		}

		// They actually had permissions
		cmdm.Logger.Printf("P: TRUE C: %s[%s] U: %s#%s[%s] M: %s", ctx.Channel.Name, ctx.Channel.ID, ctx.User.Username, ctx.User.Discriminator, ctx.User.ID, m.Content)

		err = command.OnCommand(&ctx)
		if err != nil {
			cmdm.ErrorFunc(cmdm, ctx, err)
		}
	}

}

func removeFromSlice(slice []string, i int) []string {
	slice[len(slice)-1], slice[i] = slice[i], slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func NewManager(logger log.Logger, prefixes, owners []string, errorFunc ManagerOnError) *Manager {
	return &Manager{
		Prefixes:  prefixes,
		Owners:    owners,
		Logger:    logger,
		ErrorFunc: errorFunc,
		Commands:  make(map[string]*Command),
	}
}

func (cmdm *Manager) AddPrefix(prefix string) *Manager {
	cmdm.Prefixes = append(cmdm.Prefixes, prefix)
	return cmdm
}

func (cmdm *Manager) RemovePrefix(prefix string) *Manager {
	for i, v := range cmdm.Prefixes {
		if v == prefix {
			cmdm.Prefixes = removeFromSlice(cmdm.Prefixes, i)
		}
	}
	return cmdm
}

func (cmdm *Manager) SetPrefixes(prefixes []string) *Manager {
	cmdm.Prefixes = prefixes
	return cmdm
}

func (cmdm *Manager) GetPrefixes() []string {
	return cmdm.Prefixes
}

func (cmdm *Manager) RegisterCommand(cmd *Command) *Manager {
	if cmd.Name != "" {
		cmdm.addCommand(cmd.Name, cmd)
		if cmd.Aliases != nil {
			for _, v := range cmd.Aliases {
				cmdm.addCommand(v, cmd)
			}
		}
	}
	return cmdm
}

func (cmdm *Manager) UnregisterCommand(cmd *Command) *Manager {
	if cmd != nil {
		if cmd.Name != "" {
			cmdm.removeCommand(cmd.Name)
		}
		if cmd.Aliases != nil {
			for _, c := range cmd.Aliases {
				cmdm.removeCommand(c)
			}
		}
	}
	return cmdm
}

func (cmdm *Manager) addCommand(name string, cmd *Command) *Manager {
	cmdm.Commands[name] = cmd
	return cmdm
}

func (cmdm *Manager) removeCommand(name string) *Manager {
	delete(cmdm.Commands, name)
	return cmdm
}

func (cmdm *Manager) GetCommand(name string) *Command {
	val, ok := cmdm.Commands[name]
	if ok {
		return val
	}
	return nil
}

func (cmdm *Manager) IsOwner(id string) bool {
	for _, o := range cmdm.Owners {
		if id == o {
			return true
		}
	}
	return false
}
