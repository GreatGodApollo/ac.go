package cmds

import "github.com/GreatGodApollo/acgo/permissions"

// A CommandFunc is ran whenever a CommandManager gets a message supposed to run the given command.
type CommandFunc func(CommandContext, []string) error
type CommandArgFunc func([]string) interface{}

// A Command represents any given command contained in a bot.
type Command struct {
	// The name of the command (What it will be triggered by).
	Name string

	// Command aliases
	Aliases []string

	// The command's description.
	Description string

	// If the command is only able to be ran by an owner.
	OwnerOnly bool

	// If the command is hidden from help.
	Hidden bool

	// The permissions the user is required to have to execute the command.
	UserPermissions permissions.Permission

	// The permissions the bot is required to have to execute the command.
	BotPermissions permissions.Permission

	// The CommandType designates where the command can be ran.
	Type CommandType

	// The function that will be executed whenever a message fits the criteria to execute the command.
	Run CommandFunc

	// The function that will be ran to process arguments
	ProcessArgs CommandArgFunc
}

// A CommandType represents the locations commands can be used.
type CommandType int

const (
	// A Command that is only supposed to run in a personal message
	CommandTypePM CommandType = iota

	// A command that is only supposed to run in a Guild
	CommandTypeGuild

	// A Command that is able to run anywhere
	CommandTypeEverywhere
)
