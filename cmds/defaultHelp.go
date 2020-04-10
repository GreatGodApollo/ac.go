package cmds

func DefaultHelp() *Command {
	return &Command{
		Name:        "help",
		Description: "Get some help with commands!",
		OwnerOnly:   false,
		Hidden:      false,
		Type:        Everywhere,
		Execute:     helpExecute,
		SubOnly:     false,
		UserPerms:   0,
		BotPerms:    0,
	}
}

func helpExecute(ctx *Context) error {

}
