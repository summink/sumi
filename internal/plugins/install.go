package plugins

import (
	"context"
	"encoding/json"
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

func getInstalledVersion(pluginName string) string {
	manifestPath := path.Join(cfg.PluginDir, pluginName, manifestFile)
	if !ufs.Exists(manifestPath) {
		return ""
	}

	data, err := ufs.ReadFileByByte(manifestPath)
	if err != nil {
		return ""
	}

	var manifest PluginManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return ""
	}

	return manifest.Version
}

func downloadPlugin(repo string, version string, system string, arch string) {
	client := github.NewClient(nil)
	releases, _, err := client.Repositories.ListReleases(context.Background(), orgName, repo, &github.ListOptions{})

	if err != nil {
		common.Exit("failed to list releases: " + err.Error())
	}

	if len(releases) == 0 {
		common.Exit("no releases found for " + repo)
	}

	assetName := system + "-" + arch

	if system == "windows" {
		assetName += ".exe"
	}

	outputDir := path.Join(cfg.PluginDir, strings.ReplaceAll(repo, pluginPrefix, ""))

	var targetRelease *github.RepositoryRelease

	if version == "" {
		targetRelease = releases[0]
	} else {
		normalizedVersion := version
		if !strings.HasPrefix(version, "v") {
			normalizedVersion = "v" + version
		}

		for _, release := range releases {
			tag := release.GetTagName()
			if tag == version || tag == normalizedVersion {
				targetRelease = release
				break
			}
		}

		if targetRelease == nil {
			common.Exit(fmt.Sprintf("version %s not found for %s", version, repo))
		}
	}

	for _, asset := range targetRelease.Assets {
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

			fmt.Printf("install plugin %s@%s success to %s\n", repo, targetRelease.GetTagName(), output)
			return
		}
	}

	common.Exit(fmt.Sprintf("no asset found for %s-%s in release %s", system, arch, targetRelease.GetTagName()))
}

func install() *command.SCommand {
	name := command.NewCommandArg[string]("name").ChangeDescription("install plugin name")
	version := command.NewCommandArg[string]("version").ChangeDescription("install plugin version (optional, e.g. v1.0.0 or 1.0.0)")
	force := command.NewCommandFlag[bool]("force").ChangeDescription("force install, overwrite existing plugin")

	cmd := command.NewCommand("install").
		ChangeDescription("Install a plugin").
		AddArgs(name, version).
		AddFlags(force).
		RegisterHandler(func(cmd *command.SCommand) {
			system := runtime.GOOS
			arch := runtime.GOARCH

			pluginDir := path.Join(cfg.PluginDir, name.Value)
			installedVersion := getInstalledVersion(name.Value)

			if ufs.Exists(pluginDir) {
				if installedVersion != "" {
					fmt.Printf("Plugin %s@%s is already installed.\n", name.Value, installedVersion)
				} else {
					fmt.Printf("Plugin %s is already installed.\n", name.Value)
				}

				if !force.Value {
					fmt.Println("Use --force to overwrite.")
					return
				}

				fmt.Println("Force mode enabled, reinstalling...")
				os.RemoveAll(pluginDir)
			}

			downloadPlugin(pluginPrefix+name.Value, version.Value, system, arch)
		})

	return cmd
}
