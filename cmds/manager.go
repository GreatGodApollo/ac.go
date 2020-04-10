/*
 * Vi - A Discord Bot written in Go
 * Copyright (C) 2019  Brett Bender
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// The Commands package both contains the Manager framework and the bot commands.
// Everything is pretty modular and can be adapted to your own use cases.
package cmds

import (
	"github.com/GreatGodApollo/acgo/permissions"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"strings"
)

// CommandHandler works as the Manager's message listener.
// It returns nothing.
func (cmdm *Manager) CommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot && cmdm.IgnoreBots {
		return
	}

	var prefix string
	var contains bool
	var err error
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

	channel, _ := s.Channel(m.ChannelID)

	if command, exist, _ := cmdm.GetCommand(cmd[0]); exist {
		var inDm bool
		if channel.Type == discordgo.ChannelTypeDM {
			inDm = true
		}

		// Check UserPermissions
		if command.Type != CommandTypePM && !inDm && !permissions.Check(s, m.GuildID, m.Author.ID, command.UserPermissions) {
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
				cmdm.OnErrorFunc(cmdm, Context{}, err)
			}
			cmdm.Logger.Debugf("P: FALSE C: %s[%s] U: %s#%s[%s] M: %s", channel.Name, m.ChannelID, m.Author.Username, m.Author.Discriminator, m.Author.ID, m.Content)
			return
		}

		// Check BotPermissions
		if command.Type != CommandTypePM && !inDm && !permissions.Check(s, m.GuildID, s.State.User.ID, command.BotPermissions) {
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
				cmdm.OnErrorFunc(cmdm, Context{}, err)
			}
			cmdm.Logger.Debugf("P: FALSE C: %s[%s] U: %s#%s[%s] M: %s", channel.Name, m.ChannelID, m.Author.Username, m.Author.Discriminator, m.Author.ID, m.Content)
			return
		}

		// Check if it's the right channel type
		if channel.Type == discordgo.ChannelTypeDM && command.Type == CommandTypeGuild {
			embed := &discordgo.MessageEmbed{
				Title:       "Invalid Channel!",
				Description: "You cannot run this command in a private message.",
				Color:       0xff0000,
			}

			if !command.Hidden {
				_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
			}

			if err != nil {
				cmdm.OnErrorFunc(cmdm, Context{}, err)
			}
			cmdm.Logger.Debugf("P: FALSE C: %s[%s] U: %s#%s[%s] M: %s", channel.Name, m.ChannelID, m.Author.Username, m.Author.Discriminator, m.Author.ID, m.Content)
			return
		} else if channel.Type == discordgo.ChannelTypeGuildText && command.Type == CommandTypePM {
			embed := &discordgo.MessageEmbed{
				Title:       "Invalid Channel!",
				Description: "You cannot run this command in a guild.",
				Color:       0xff0000,
			}

			if !command.Hidden {
				_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
			}

			if err != nil {
				cmdm.OnErrorFunc(cmdm, Context{}, err)
			}
			cmdm.Logger.Debugf("P: FALSE C: %s[%s] U: %s#%s[%s] M: %s", channel.Name, m.ChannelID, m.Author.Username, m.Author.Discriminator, m.Author.ID, m.Content)
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
				cmdm.OnErrorFunc(cmdm, Context{}, err)
			}
			cmdm.Logger.Debugf("P: FALSE C: %s[%s] U: %s#%s[%s] M: %s", channel.Name, m.ChannelID, m.Author.Username, m.Author.Discriminator, m.Author.ID, m.Content)
			return
		}

		// They actually had permissions
		cmdm.Logger.Debugf("P: TRUE C: %s[%s] U: %s#%s[%s] M: %s", channel.Name, m.ChannelID, m.Author.Username, m.Author.Discriminator, m.Author.ID, m.Content)
		guild, _ := s.Guild(m.GuildID)
		member, _ := s.State.Member(m.GuildID, m.Author.ID)

		ctx := Context{
			Session: s,
			Event:   m,
			Manager: cmdm,
			Args:    cmd[1:],
			Message: m.Message,
			User:    m.Author,
			Channel: channel,
			Guild:   guild,
			Member:  member,
		}

		err = command.Run(ctx, cmd[1:])
		if err != nil {
			cmdm.OnErrorFunc(cmdm, ctx, err)
		}
	}
}

// AddPrefix adds a new prefix to the Manager's prefix list.
// It returns nothing.
func (cmdm *Manager) AddPrefix(prefix string) {
	cmdm.Prefixes = append(cmdm.Prefixes, prefix)
}

// RemovePrefix removes a prefix from the Manager's prefix list.
// It returns nothing.
func (cmdm *Manager) RemovePrefix(prefix string) {
	for i, v := range cmdm.Prefixes {
		if v == prefix {
			cmdm.Prefixes = append(cmdm.Prefixes[:i], cmdm.Prefixes[i+1:]...)
			break
		}
	}
}

// SetPrefixes sets the Manager's prefix list.
// It returns nothing.
func (cmdm *Manager) SetPrefixes(prefixes []string) {
	cmdm.Prefixes = prefixes
}

// GetPrefixes gets the Manager's prefix list.
// It returns a string array.
func (cmdm *Manager) GetPrefixes() []string {
	return cmdm.Prefixes
}

// AddNewCommand adds a new command to the Manager's command list.
// It returns nothing.
func (cmdm *Manager) AddNewCommand(name string, aliases []string, desc string, owneronly, hidden bool, userperms, botperms permissions.Permission,
	cmdType CommandType, run CommandFunc) {
	var cmd *Command
	if _, exists, _ := cmdm.GetCommand(name); !exists {
		cmd = &Command{
			name, aliases, desc, owneronly, hidden, userperms, botperms, cmdType, run, nil,
		}
	}
	*cmdm.Commands = append(*cmdm.Commands, cmd)
}

// AddCommand adds an existent command to the Manager's command list.
// It returns nothing.
func (cmdm *Manager) AddCommand(cmd *Command) {
	if _, exists, _ := cmdm.GetCommand(cmd.Name); !exists {
		*cmdm.Commands = append(*cmdm.Commands, cmd)
	}
}

func (cmdm *Manager) GetCommand(name string) (cmd *Command, exists bool, index int) {
	for i, c := range *cmdm.Commands {
		if c.Name == name {
			return c, true, i
		}
		for _, a := range c.Aliases {
			if a == name {
				return c, true, i
			}
		}
	}
	return nil, false, 0
}

// RemoveCommand removes a command from the Manager's command list.
// It returns nothing.
func (cmdm *Manager) RemoveCommand(name string) {
	if _, exists, index := cmdm.GetCommand(name); exists {
		*cmdm.Commands = RemoveCommandFromSlice(*cmdm.Commands, index)
	}
}

// IsOwner checks if a user ID is is in the owner list.
// It returns a bool.
func (cmdm *Manager) IsOwner(id string) bool {
	for _, o := range cmdm.Owners {
		if id == o {
			return true
		}
	}
	return false
}

// NewManager instantiates a new Manager.
// It returns a Manager.
func NewManager(l *logrus.Logger, ignoreBots bool, errorFunc ManagerOnErrorFunc) Manager {
	return Manager{
		Prefixes:    []string{},
		Owners:      []string{},
		Commands:    &[]*Command{},
		Logger:      l,
		IgnoreBots:  ignoreBots,
		OnErrorFunc: errorFunc,
	}
}

// A Manager represents a set of prefixes, owners and commands, with some extra utility to create a command handler.
type Manager struct {
	// The array of prefixes a Manager will respond to.
	Prefixes []string

	// The array of IDs that will be considered a bot owner.
	Owners []string

	// The bot instance Logger.
	Logger *logrus.Logger

	// The map of Commands in the Manager.
	Commands *[]*Command

	// If the Manager ignores bots or not.
	IgnoreBots bool

	// The function that will be ran when the Manager encounters an error.
	OnErrorFunc ManagerOnErrorFunc
}

// A ManagerOnErrorFunc is a function that will run whenever the Manager encounters an error.
type ManagerOnErrorFunc func(cmdm *Manager, ctx Context, err error)

func RemoveCommandFromSlice(s []*Command, i int) []*Command {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
