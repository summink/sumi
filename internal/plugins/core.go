package plugins

var manifestFile = "manifest.json"
var orgName = "summink"

type PluginPlatform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type PluginManifest struct {
	Name        string           `json:"name"`
	Version     string           `json:"version"`
	Description string           `json:"description"`
	Author      string           `json:"author"`
	License     string           `json:"license"`
	Repo        string           `json:"repo"`
	Doc         string           `json:"doc"`
	Platforms   []PluginPlatform `json:"platforms"`
	Tags        []string         `json:"tags"`
}

type PluginConfig struct {
	PluginDir string `json:"pluginDir"`
}
