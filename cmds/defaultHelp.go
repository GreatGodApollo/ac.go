package cmds

import (
	"fmt"
	"github.com/GreatGodApollo/acgo/embeds"
	"github.com/GreatGodApollo/acgo/permissions"
	"sort"
)

var embedColor int

func DefaultHelp(ec int) *Command {
	embedColor = ec

	return &Command{
		Name:            "help",
		Description:     "Get some help with the bot.",
		OwnerOnly:       false,
		Hidden:          false,
		UserPermissions: 0,
		BotPermissions:  permissions.PermissionMessagesSend | permissions.PermissionMessagesEmbedLinks,
		Type:            CommandTypeEverywhere,
		Run:             helpCommandFunc,
		ProcessArgs:     helpArgsFunc,
	}
}

// A HelpCommandArgs is passed into a CommandContext. It provides the necessary information for a help command to run.
type helpCommandArgs struct {
	// The name of the command the user is searching for
	Command string

	// The rest of the arguments provided
	Rest []string
}

// HelpArgsFunc is a CommandArgFunc
// It returns the proper HelpCommandArgs struct given the args provided
// It returns an empty struct if no args are provided
func helpArgsFunc(args []string) interface{} {
	if len(args) == 1 {
		return helpCommandArgs{Command: args[0], Rest: nil}
	} else if len(args) > 1 {
		return helpCommandArgs{Command: args[0], Rest: args[1:]}
	}
	return helpCommandArgs{}
}

// HelpCommandFunc is a CommandRunFunc.
// It supplies the user a list of commands in the CommandManager it is assigned to.
// It returns an error if any occurred.
//
// Usage: {prefix}help [command]
func helpCommandFunc(ctx CommandContext, args []string) error {
	argStruct := ctx.Args.(helpCommandArgs)
	if len(args) > 0 {
		if command, has, _ := ctx.Manager.GetCommand(argStruct.Command); has {
			if command.Hidden {
				return nil
			}

			var (
				ownerOnlyString string
				typeString      string
			)

			if command.OwnerOnly {
				ownerOnlyString = "Yes"
			} else {
				ownerOnlyString = "No"
			}

			switch command.Type {
			case CommandTypePM:
				{
					typeString = "Private"
				}
			case CommandTypeGuild:
				{
					typeString = "Guild-only"
				}
			case CommandTypeEverywhere:
				{
					typeString = "Anywhere"
				}
			}

			var alList string
			for i, a := range command.Aliases {
				if i == len(command.Aliases)-1 {
					alList += fmt.Sprintf("`%s`", a)
				} else {
					alList += fmt.Sprintf("`%s` ", a)
				}
			}
			if alList == "" {
				alList = "No Aliases"
			}

			e := embeds.NewEmbed().
				SetTitle(fmt.Sprintf("Help for `%s`!", command.Name)).
				SetColor(embedColor).
				SetDescription(command.Description).
				AddInlineField("Owner Only?", ownerOnlyString).
				AddInlineField("Usage?", typeString).
				AddField("Aliases", alList)

			_, err := ctx.ReplyEmbed(e.MessageEmbed)
			return err
		} else {
			e := embeds.NewEmbed().
				SetTitle("Command does not exist.").
				SetColor(0xFF0000).
				SetDescription(fmt.Sprintf("Please use `%shelp` for a list of commands.", ctx.Manager.Prefixes[0]))
			_, err := ctx.ReplyEmbed(e.MessageEmbed)
			return err
		}
	}
	m := ctx.Manager.Commands

	keys := make([]string, 0, len(*m))
	for _, k := range *m {
		n := k.Name
		keys = append(keys, n)
	}
	sort.Strings(keys)

	var list string
	for _, k := range keys {
		cmd, _, _ := ctx.Manager.GetCommand(k)
		if !cmd.Hidden {
			list += fmt.Sprintf("**%s** - `%s`\n", cmd.Name, cmd.Description)
		}
	}

	var footer string

	if len(*m) == 1 {
		footer = "There is 1 command."
	} else {
		footer = fmt.Sprintf("There are %d commands.", len(*m))
	}

	embed := embeds.NewEmbed().
		SetTitle("Commands:").
		SetDescription(list).
		SetColor(embedColor).
		SetFooter(footer)

	_, err := ctx.ReplyEmbed(embed.MessageEmbed)
	return err
}
