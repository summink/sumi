package plugins

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/InkShaStudio/go-command"
	"github.com/google/go-github/v84/github"
	"github.com/inksha/sumi/internal/utils/api"
	"github.com/inksha/sumi/internal/utils/common"
	"github.com/inksha/sumi/internal/utils/ufs"
)

func cleanupOnFailure(dir string) {
	if dir != "" && ufs.Exists(dir) {
		os.RemoveAll(dir)
	}
}

func downloadPlugin(repo string, version string, system string, arch string) {
	client := github.NewClient(nil)
	releases, _, err := client.Repositories.ListReleases(context.Background(), orgName, repo, &github.ListOptions{})

	if err != nil {
		common.Exit("failed to list releases: " + err.Error())
	}

	assetName := system + "-" + arch

	if system == "windows" {
		assetName += ".exe"
	}

	outputDir := path.Join(cfg.PluginDir, strings.ReplaceAll(repo, pluginPrefix, ""))

topReleases:
	for _, release := range releases {
		for _, asset := range release.Assets {

			if strings.EqualFold(*asset.Name, assetName) || strings.HasSuffix(*asset.Name, assetName) {

				assetURL := asset.GetURL()
				assetInfo, err := api.Get(assetURL)
				if err != nil {
					common.Exit("failed to get asset info: " + err.Error())
				}

				browserDownloadURL, ok := assetInfo["browser_download_url"].(string)
				if !ok {
					common.Exit("failed to get download URL: invalid response format")
				}

				assetData, err := api.GetRaw(browserDownloadURL)
				if err != nil {
					common.Exit("failed to download asset: " + err.Error())
				}

				if len(assetData) == 0 {
					common.Exit("failed to download asset: empty response")
				}

				output := path.Join(outputDir, assetName)

				if err := ufs.MkDirIfNotExist(outputDir, true); err != nil {
					common.Exit("failed to create plugin directory: " + err.Error())
				}

				if err := ufs.WriteFileByByte(output, assetData); err != nil {
					cleanupOnFailure(outputDir)
					common.Exit("failed to write plugin file: " + err.Error())
				}

				if err := os.Chmod(output, 0755); err != nil {
					cleanupOnFailure(outputDir)
					common.Exit("failed to set plugin permissions: " + err.Error())
				}

				fmt.Printf("install plugin %s%s success to %s\n", repo, release.GetTagName(), output)

				break topReleases
			}
		}
	}
}

func install() *command.SCommand {
	name := command.NewCommandArg[string]("name").ChangeDescription("install plugin name")
	version := command.NewCommandArg[string]("version").ChangeDescription("install plugin version")

	cmd := command.NewCommand("install").
		ChangeDescription("Install a plugin").
		AddArgs(name, version).
		RegisterHandler(func(cmd *command.SCommand) {
			system := runtime.GOOS
			arch := runtime.GOARCH

			if ufs.Exists(path.Join(cfg.PluginDir, name.Value)) {
				println("Plugin " + name.Value + " already installed!")
				return
			}

			if version.Value != "" {
				println("Version parameter is not currently supported!")
			}

			downloadPlugin(pluginPrefix+name.Value, version.Value, system, arch)
		})

	return cmd
}
