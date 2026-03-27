package plugins

import (
	"os"
	"path"

	"github.com/InkShaStudio/go-command"
)

var cfg PluginConfig

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	cfg = PluginConfig{
		PluginDir: path.Join(homeDir, ".sumi", "plugins"),
	}
}

func RegisterCommand() *command.SCommand {

	cmd := command.
		NewCommand("plugin").
		ChangeDescription("Plugin commands").
		AddSubCommand(
			list(),
			install(),
			uninstall(),
			update(),
			info(),
			newPlugin(),
		)

	return cmd
}
