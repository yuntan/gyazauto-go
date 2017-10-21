package main

import (
	"github.com/op/go-logging"

	"os"
)

const (
	colorLogFormat = "%{time:2006/01/02 15:04:05.000000} %{shortfunc:-15.15s} | %{color:bold}%{level:.4s}%{color:reset} %{color}%{message}%{color:reset}"
)

var (
	log = logging.MustGetLogger("gyazauto")
)

func setupLogger() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	var format logging.Formatter
	format = logging.MustStringFormatter(colorLogFormat)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	if verbose {
		backendLeveled.SetLevel(logging.DEBUG, "")
	} else {
		backendLeveled.SetLevel(logging.INFO, "")
	}
	log.SetBackend(backendLeveled)
}
