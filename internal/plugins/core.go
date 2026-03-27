package plugins

import sumicc "github.com/summink/sumi-common-command"

var manifestFile = "manifest.json"
var orgName = "summink"
var pluginPrefix = "sumi-plugin-"

type PluginDetails struct {
	Name     string                `json:"name"`
	Manifest sumicc.PluginManifest `json:"manifest"`
	Execute  string                `json:"execute"`
}

type PluginConfig struct {
	PluginDir string `json:"pluginDir"`
}
