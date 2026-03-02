package plugins

import (
	"os"
	"path"

	"github.com/InkShaStudio/go-command"
)

var cfg = PluginConfig{
	PluginDir: path.Join(os.Getenv("USERPROFILE"), ".sumi", "plugins"),
}

func RegisterCommand() *command.SCommand {

	cmd := command.
		NewCommand("plugin").
		ChangeDescription("Plugin commands").
		AddSubCommand(
			list(),
		)

	return cmd
}
