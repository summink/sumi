package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/InkShaStudio/go-command"
	"github.com/google/go-github/v84/github"
	"github.com/inksha/sumi/internal/utils/api"
	"github.com/inksha/sumi/internal/utils/common"
	"github.com/inksha/sumi/internal/utils/ufs"
)

func findPlugin() []PluginManifest {
	plugins := []PluginManifest{}

	if !ufs.Exists(cfg.PluginDir) {
		ufs.MkDir(cfg.PluginDir, true)
	}

	list := ufs.ListDir(cfg.PluginDir)

	if len(list) == 0 {
		println("not install plugins")
		return plugins
	}

	for _, item := range list {
		manifest := path.Join(item, manifestFile)
		if ufs.Exists(manifest) {
			if value, err := ufs.ReadFileByByte(manifest); err == nil {
				var manifest PluginManifest
				if err := json.Unmarshal(value, &manifest); err == nil {
					plugins = append(plugins, manifest)
				}
			}
		}

	}

	return plugins
}

func findPluginByGithub() []PluginManifest {
	client := github.NewClient(nil)

	opts := &github.RepositoryListByOrgOptions{Type: "public"}

	repos, _, err := client.Repositories.ListByOrg(context.Background(), orgName, opts)

	if err != nil {
		common.Exit(err.Error())
	}

	plugins := []PluginManifest{}

	for _, repo := range repos {

		if strings.Contains(repo.GetName(), "sumi-plugin-") && !*repo.IsTemplate {

			url := strings.ReplaceAll(*repo.ContentsURL, "{+path}", manifestFile)

			reps, _ := api.Get(url)

			var manifest PluginManifest

			if url, ok := reps["download_url"].(string); ok {
				rawManifest, _ := api.GetRaw(url)
				if err := json.Unmarshal(rawManifest, &manifest); err == nil {
					plugins = append(plugins, manifest)
				}
			}
		}
	}

	return plugins
}

func list() *command.SCommand {
	online := command.NewCommandFlag[bool]("online").ChangeDescription("get plugins from online")

	cmd := command.NewCommand("list").
		ChangeDescription("list plugins").
		AddFlags(online).
		RegisterHandler(func(cmd *command.SCommand) {
			plugins := []PluginManifest{}

			if online.Value {
				plugins = findPluginByGithub()
			} else {
				plugins = findPlugin()
			}

			maxNameLen := 0

			for _, plugin := range plugins {
				maxNameLen = max(maxNameLen, len(plugin.Name))
			}

			for _, plugin := range plugins {
				println(fmt.Sprintf("%s %s", plugin.Name+strings.Repeat(" ", maxNameLen-len(plugin.Name)), plugin.Version))

				if online.Value {
					localDir := path.Join(cfg.PluginDir, "."+plugin.Name)

					if !ufs.Exists(localDir) {
						ufs.MkDir(localDir, true)
					}

					if !ufs.Exists(path.Join(localDir, manifestFile)) {
						if value, err := json.Marshal(plugin); err == nil {
							ufs.WriteFileByByte(path.Join(localDir, manifestFile), value)
						}
					}
				}
			}
		})

	return cmd
}
