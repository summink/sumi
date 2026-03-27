package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/InkShaStudio/go-command"
	"github.com/inksha/sumi/internal/utils/common"
	"github.com/inksha/sumi/internal/utils/ufs"
)

const templateRepo = "https://github.com/summink/sumi-plugin-project-template.git"

func createPluginProjectWithOptions(name, author, description, version, license, tags, platforms string) {
	if name == "" {
		common.Exit("plugin name is required")
	}

	fullName := pluginPrefix + name
	targetDir := fullName

	if ufs.Exists(targetDir) {
		common.Exit(fmt.Sprintf("directory %s already exists", targetDir))
	}

	// Clone template repository
	fmt.Printf("Creating plugin project %s...\n", fullName)
	cloneCmd := exec.Command("git", "clone", "--depth=1", templateRepo, targetDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		common.Exit(fmt.Sprintf("failed to clone template: %s", err.Error()))
	}

	// Remove .git directory to start fresh
	gitDir := path.Join(targetDir, ".git")
	if err := os.RemoveAll(gitDir); err != nil {
		common.Exit(fmt.Sprintf("failed to remove .git: %s", err.Error()))
	}

	// Replace template placeholders
	replaceInFile(path.Join(targetDir, "README.md"), map[string]string{
		"sumi-plugin-project-template":         fullName,
		"A template for creating sumi plugins": description,
		"My Sumi Plugin":                       fullName,
		"myplugin":                             name,
	})

	replaceInFile(path.Join(targetDir, "go.mod"), map[string]string{
		"github.com/summink/sumi-plugin-project-template": fmt.Sprintf("github.com/%s/%s", orgName, fullName),
	})

	replaceInFile(path.Join(targetDir, "release-please-config.json"), map[string]string{
		"sumi-plugin-project-template": fullName,
	})

	// Update manifest.json
	updateManifest(path.Join(targetDir, "manifest.json"), fullName, author, description, version, license, tags, platforms)

	// Initialize new git repository
	initCmd := exec.Command("git", "init")
	initCmd.Dir = targetDir
	initCmd.Run()

	fmt.Printf("\n✓ Plugin project created: %s\n", targetDir)
	fmt.Println("\nNext steps:")
	fmt.Printf("  cd %s\n", targetDir)
	fmt.Println("  go mod tidy")
	fmt.Println("  git add . && git commit -m \"Initial commit\"")
}

func replaceInFile(filePath string, replacements map[string]string) {
	if !ufs.Exists(filePath) {
		return
	}

	data, err := ufs.ReadFileByByte(filePath)
	if err != nil {
		return
	}

	content := string(data)
	for old, new := range replacements {
		content = strings.ReplaceAll(content, old, new)
	}

	ufs.WriteFileByByte(filePath, []byte(content))
}

func updateManifest(filePath, name, author, description, version, license, tags, platforms string) {
	if !ufs.Exists(filePath) {
		return
	}

	data, err := ufs.ReadFileByByte(filePath)
	if err != nil {
		return
	}

	var manifest map[string]any
	if err := json.Unmarshal(data, &manifest); err != nil {
		return
	}

	// Update fields, preserve others like $schema
	manifest["name"] = name
	manifest["version"] = version
	manifest["repo"] = fmt.Sprintf("https://github.com/%s/%s", orgName, name)
	manifest["doc"] = manifest["repo"]

	if author != "" {
		manifest["author"] = author
	}
	if description != "" {
		manifest["description"] = description
	}
	if license != "" {
		manifest["license"] = license
	}
	if tags != "" {
		manifest["tags"] = strings.Split(tags, ",")
	}
	if platforms != "" {
		platformList := []map[string]string{}
		for _, p := range strings.Split(platforms, ",") {
			parts := strings.Split(strings.TrimSpace(p), "/")
			if len(parts) == 2 {
				platformList = append(platformList, map[string]string{
					"os":   parts[0],
					"arch": parts[1],
				})
			}
		}
		if len(platformList) > 0 {
			manifest["platforms"] = platformList
		}
	}

	newData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return
	}

	ufs.WriteFileByByte(filePath, newData)
}

func newPlugin() *command.SCommand {
	name := command.NewCommandArg[string]("name").ChangeDescription("plugin name (without sumi-plugin- prefix)")
	author := command.NewCommandFlag[string]("author").ChangeDescription("plugin author name")
	description := command.NewCommandFlag[string]("description").ChangeDescription("plugin description")
	version := command.NewCommandFlag[string]("version").ChangeDescription("plugin version (default: 0.0.1)")
	license := command.NewCommandFlag[string]("license").ChangeDescription("plugin license (default: MIT)")
	tags := command.NewCommandFlag[string]("tags").ChangeDescription("plugin tags, comma separated (e.g. plugin,tool)")
	platforms := command.NewCommandFlag[string]("platforms").ChangeDescription("supported platforms, comma separated (e.g. windows/amd64,linux/amd64)")

	cmd := command.NewCommand("new").
		ChangeDescription("Create a new plugin project from template").
		AddArgs(name).
		AddFlags(author, description, version, license, tags, platforms).
		RegisterHandler(func(cmd *command.SCommand) {
			v := version.Value
			if v == "" {
				v = "0.0.1"
			}
			l := license.Value
			if l == "" {
				l = "MIT"
			}
			createPluginProjectWithOptions(name.Value, author.Value, description.Value, v, l, tags.Value, platforms.Value)
		})

	return cmd
}
