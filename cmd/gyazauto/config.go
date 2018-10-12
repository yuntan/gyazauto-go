package main

import (
	"github.com/kirsle/configdir"
	"gopkg.in/yaml.v2"

	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	confName = "config.yaml"
)

var (
	confDir  = configdir.LocalConfig(appNameLower)
	confPath = filepath.Join(confDir, confName)
)

var config struct {
	WatchDirs   []string `yaml:"watch_dirs,omitempty"`
	AccessToken string   `yaml:"access_token,omitempty"`
}

func loadConfig() {
	log.Debug("loadConfig")

	err := configdir.MakePath(confDir)
	if err != nil {
		log.Fatalf("cannnot create config dir: %v", err)
		return
	}

	if _, err = os.Stat(confPath); os.IsNotExist(err) {
		log.Debug("loadConfig: %s not found", confPath)
		return
	}

	log.Info("Load config from %s", confPath)
	b, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatalf("Failed to load config from %s, please recheck permission: %v",
			confPath, err)
		return
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		log.Fatalf("Failed to load config from %s, please recheck syntax: %v",
			confPath, err)
		return
	}

	log.Debug("loadConfig: OK")
}

func saveConfig() {
	log.Debug("saveConfig")

	err := configdir.MakePath(confDir)
	if err != nil {
		log.Fatalf("Cannnot create config dir: %v", err)
		return
	}

	b, err := yaml.Marshal(config)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = ioutil.WriteFile(confPath, b, 0600)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Debug("saveConfig: OK")
}
