package cmd

import (
	"fmt"
	"os"

	"github.com/InkShaStudio/go-command"
	"github.com/inksha/sumi/internal/plugins"
	"github.com/inksha/sumi/internal/sumi_hash"
	"github.com/inksha/sumi/internal/sumi_rand"
	"github.com/inksha/sumi/internal/sumi_template"
)

func Execute() {
	sumi := command.NewCommand("sumi").
		ChangeDescription("Always want to see summer in you eyes.").
		RegisterHandler(func(cmd *command.SCommand) {
			fmt.Println("Hello, Sumi!")
		})

	sumi.AddSubCommand(
		sumi_hash.RegisterCommand(),
		sumi_template.RegisterCommand(),
		sumi_rand.RegisterCommand(),
		plugins.RegisterCommand(),
	)

	cmd := command.RegisterCommand(sumi)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
