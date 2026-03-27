package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/InkShaStudio/go-command"
	"github.com/google/go-github/v84/github"
	"github.com/inksha/sumi/internal/utils/api"
	"github.com/inksha/sumi/internal/utils/common"
	"github.com/inksha/sumi/internal/utils/ufs"
	sumicc "github.com/summink/sumi-common-command"
)

func findPlugin() []sumicc.PluginManifest {
	plugins := []sumicc.PluginManifest{}

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
				var manifest sumicc.PluginManifest
				if err := json.Unmarshal(value, &manifest); err == nil {
					plugins = append(plugins, manifest)
				}
			}
		}

	}

	return plugins
}

func findPluginByGithub() []sumicc.PluginManifest {
	client := github.NewClient(nil)

	opts := &github.RepositoryListByOrgOptions{Type: "public"}

	repos, _, err := client.Repositories.ListByOrg(context.Background(), orgName, opts)

	if err != nil {
		common.Exit(err.Error())
	}

	plugins := []sumicc.PluginManifest{}

	for _, repo := range repos {

		if strings.Contains(repo.GetName(), pluginPrefix) && !*repo.IsTemplate {

			url := strings.ReplaceAll(*repo.ContentsURL, "{+path}", manifestFile)

			reps, err := api.Get(url)
			if err != nil {
				fmt.Printf("Warning: failed to get manifest for %s: %v\n", repo.GetName(), err)
				continue
			}

			var manifest sumicc.PluginManifest

			downloadURL, ok := reps["download_url"].(string)
			if !ok {
				fmt.Printf("Warning: invalid manifest URL for %s\n", repo.GetName())
				continue
			}

			rawManifest, err := api.GetRaw(downloadURL)
			if err != nil {
				fmt.Printf("Warning: failed to download manifest for %s: %v\n", repo.GetName(), err)
				continue
			}

			if err := json.Unmarshal(rawManifest, &manifest); err == nil {
				plugins = append(plugins, manifest)
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
			plugins := []sumicc.PluginManifest{}

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

func ListPlugin() []PluginDetails {
	plugins := findPlugin()

	localPlugins := []PluginDetails{}

	system := runtime.GOOS
	arch := runtime.GOARCH
	execute := system + "-" + arch

	if system == "windows" {
		execute += ".exe"
	}

	for _, plugin := range plugins {
		name := strings.ReplaceAll(plugin.Name, pluginPrefix, "")

		localPlugins = append(localPlugins, PluginDetails{
			Name:     name,
			Manifest: plugin,
			Execute:  path.Join(cfg.PluginDir, name, execute),
		})
	}

	return localPlugins
}
