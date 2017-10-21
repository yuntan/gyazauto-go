package main

import (
	"github.com/shibukawa/configdir"
	"gopkg.in/yaml.v2"

	"path/filepath"
)

const (
	configName = "config.yaml"
)

var (
	confDir = configdir.New("gyazauto", "gyazauto")
)

var config struct {
	WatchDirs   []string `yaml:"watch_dirs,omitempty"`
	AccessToken string   `yaml:"access_token,omitempty"`
}

// func getYamlPath() string {
// 	// TODO AppData on windows
// 	path, err := homedir.Expand("~/.config/gyazauto/config.yaml")
// 	if err != nil {
// 		log.Fatal(err)
// 		return ""
// 	}
// 	return path
// }

func loadConfig() {
	folder := confDir.QueryFolderContainsFile(configName)
	if folder == nil {
		log.Debug("loadConfig: %s not found", configName)
		return
	}

	b, err := folder.ReadFile(configName)
	if err != nil {
		log.Debug("loadConfig: %v", err)
		log.Error("Failed to load config from %s, please recheck permission",
			filepath.Join(folder.Path, configName))
		return
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		log.Debug("loadConfig: %v", err)
		log.Error("Failed to load config from %s, please recheck syntax",
			filepath.Join(folder.Path, configName))
		return
	}

	log.Info("Load config from %s", filepath.Join(folder.Path, configName))
}

func saveConfig() {
	// err := os.MkdirAll(filepath.Dir(path), 0700)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	b, err := yaml.Marshal(config)
	if err != nil {
		log.Fatal(err)
		return
	}

	// store to user folder
	// Windows: %LOCALAPPDATA%
	// Linux: ~/.config
	folders := confDir.QueryFolders(configdir.Global)
	err = folders[0].WriteFile(configName, b)
	if err != nil {
		log.Fatal(err)
		return
	}
}
