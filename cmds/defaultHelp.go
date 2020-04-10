package cmds

import "github.com/GreatGodApollo/acgo/permissions"

var DefaultHelp = &Command{
	Name:        "help",
	Description: "Get some help with commands!",
	OwnerOnly:   false,
	Hidden:      false,
	Type:        Everywhere,
	Execute:     helpExecute,
	SubOnly:     false,
	UserPerms:   0,
	BotPerms:    permissions.PermissionMessagesSend,
}

func helpExecute(ctx Context) error {
	_, err := ctx.Reply("WIP")
	return err
}
