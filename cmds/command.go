package cmds

import (
	"errors"
	"github.com/GreatGodApollo/acgo/permissions"
)

type Command struct {
	Name        string
	Aliases     []string
	Description string

	OwnerOnly bool
	Hidden    bool

	Type CommandType

	Execute CommandFunc

	SubCommands map[string]*Command
	SubOnly     bool

	UserPerms permissions.Permission
	BotPerms  permissions.Permission
}

func (cmd *Command) OnCommand(ctx *Context) error {
	if cmd.SubOnly && len(ctx.Args) > 0 {
		if cmd.SubCommands != nil {
			sub := cmd.GetSubCommand(ctx.Args[0])
			if sub != nil {
				ctx.Args = ctx.Args[1:]
				return sub.OnCommand(ctx)
			} else {
				return errors.New("unknown sub")
			}
		} else {
			return errors.New("no subs")
		}
	} else if cmd.SubOnly {
		return errors.New("must provide sub")
	} else {
		if len(ctx.Args) > 0 {
			if sub := cmd.GetSubCommand(ctx.Args[0]); sub != nil {
				return sub.OnCommand(ctx)
			}
		}
		return cmd.Execute(ctx)
	}
}

func (cmd *Command) RegisterSubCommand(c *Command) *Command {
	if c != nil {
		if c.Name != "" {
			cmd.addSubCommand(c.Name, c)
		}
		if c.Aliases != nil {
			for _, v := range c.Aliases {
				cmd.addSubCommand(v, c)
			}
		}
	}
	return cmd
}

func (cmd *Command) UnregisterSubCommand(c *Command) *Command {
	if cmd != nil {
		if cmd.Name != "" {
			cmd.removeSubCommand(cmd.Name)
		}
		if cmd.Aliases != nil {
			for _, c := range cmd.Aliases {
				cmd.removeSubCommand(c)
			}
		}
	}
	return cmd
}

func (cmd *Command) GetSubCommand(name string) *Command {
	val, ok := cmd.SubCommands[name]
	if ok {
		return val
	}
	return nil
}

func (cmd *Command) addSubCommand(name string, c *Command) *Command {
	cmd.SubCommands[name] = c
	return cmd
}

func (cmd *Command) removeSubCommand(name string) *Command {
	delete(cmd.SubCommands, name)
	return cmd
}

type CommandFunc func(*Context) error

type CommandType int

const (
	Direct CommandType = iota

	Guild

	Everywhere
)
