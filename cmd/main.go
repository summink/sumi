package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/InkShaStudio/go-command"
	"github.com/inksha/sumi/internal/plugins"
	"github.com/inksha/sumi/internal/utils/common"
)

const (
	NAME        = "sumi"
	DESCRIPTION = "Always want to see summer in you eyes."
)

func Execute() {
	localPlugins := plugins.ListPlugin()

	args := os.Args[1:]

	cmd := command.NewCommand(NAME).ChangeDescription(DESCRIPTION)

	internalCommands := map[string]func(){
		"plugin": func() {
			cmd.AddSubCommand(plugins.RegisterCommand())

			if err := command.RegisterCommand(cmd).Execute(); err != nil {
				common.Exit(err.Error())
			}
		},
	}

	if len(args) == 0 {
		println(DESCRIPTION)

		return
	}

	name := args[0]
	localCommands := []string{""}

	for cmd, handler := range internalCommands {
		localCommands = append(localCommands, cmd)

		if strings.EqualFold(cmd, name) {
			handler()
			return
		}
	}

	for _, details := range localPlugins {
		localCommands = append(localCommands, details.Name)

		if details.Name == name {
			cmd := exec.Command(details.Execute, args[1:]...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				common.Exit(err.Error())
			}
			println(string(output))
			break
		}
	}

	commands := strings.Join(localCommands, fmt.Sprintf("\n  - %s ", NAME))

	println(
		fmt.Sprintf(
			"\n%s\n%s\n\n%s%s\n",
			NAME,
			DESCRIPTION,
			"Usage: ",
			commands,
		),
	)

}
