package main

import (
	"github.com/atotto/clipboard"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/browser"
	"gopkg.in/fsnotify.v1"

	"flag"
	"fmt"
	"os"
	"time"
)

var (
	// to set it to git tag, `go install -ldflags "-X main.version=$(git describe)"`
	version = "XXX"

	// cmd flags
	verbose      bool
	openOnUpload bool

	chNotify = make(chan fsnotify.Event)
	watcher  *fsnotify.Watcher
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&openOnUpload, "open", false, "Open browser when uploaded to Gyazo")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "gyazauto %s\n\n", version)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	setupLogger()

	if clipboard.Unsupported {
		log.Warning("Clipboard feature may not works")
	}

	loadConfig()

	if config.AccessToken == "" {
		log.Info("Not authorized. Starting authorization process...")
		authorize()
	}

	if len(config.WatchDirs) == 0 {
		log.Info("You don't specify watch dir. Exiting...")
		os.Exit(0)
	}

	log.Info("You have specified %d dirs", len(config.WatchDirs))
	for _, s := range config.WatchDirs {
		log.Info("  %s", s)
	}
	setupWatcher()
	watcherLoop()
	defer watcher.Close()
}

func setupWatcher() {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to setting up filesystem watcher: %v", err)
		return
	}

	log.Info("Adding files to watcher...")
	for _, s := range config.WatchDirs {
		path, err := homedir.Expand(s)
		if err != nil {
			log.Error("Failed to expand ~: %v", err)
			continue
		}
		info, err := os.Stat(path)
		if err != nil {
			log.Error("Cannot read directory %s: %v", path, err)
			continue
		}
		if !info.IsDir() {
			log.Error("%s is not a directory", path)
			continue
		}
		err = watcher.Add(path)
		if err != nil {
			log.Error("Cannot add %s to watcher: %v", path, err)
			continue
		}
		log.Debug("%s added to watcher", path)
	}
}

func watcherLoop() {
	// event loop for watcher
	for {
		select {
		case event := <-watcher.Events:
			log.Debug("%+v", event)
			// NOTE WRITE event comes several times. Only accept CREATE event
			// and wait several seconds for the file is ready to upload.
			if event.Op != fsnotify.Create {
				continue
			}

			go func() {
				time.Sleep(5 * time.Second)
				url := upload(event.Name)
				if url == "" {
					return
				}

				fmt.Println(url) // stdout

				if openOnUpload {
					if err := browser.OpenURL(url); err != nil {
						log.Error("Failed to open browser: %v", err)
					}
				}

				// copy to clipboard
				if err := clipboard.WriteAll(url); err != nil {
					log.Error("Failed to write to clipboard: %v", err)
					return
				}
				log.Info("Copied Gyazo URL to clipboard")
			}()

		case err := <-watcher.Errors:
			log.Error("Filesystem watcher unexpectedly crashed: %v", err)
		}
	}
}
