package main

import (
"flag"
"github.com/ONSdigital/log.go/log"
"github.com/ONSdigital/publish-times/console"
"github.com/ONSdigital/publish-times/input"
"os"
)

var (
	publishLogDirEnv = "PUBLISH_LOG_DIR"
	publishLogPath   string
	publishes        = make([]os.FileInfo, 0)
)

func main() {
	console.Init()
	log.Namespace = "publish-times"

	publishLogFlag := flag.String("publishLogDir", "", "override - use this value instead of the env config")
	debug := flag.Bool("debug", false, "log out configuration values on start up")
	flag.Parse()

	if len(*publishLogFlag) > 0 {
		publishLogPath = *publishLogFlag
	} else {
		publishLogPath = os.Getenv(publishLogDirEnv)
	}

	if *debug {
		log.Event(nil, "configuration", log.Data{"publishLogPath": publishLogPath})
	}

	if err := input.Listen(); err != nil {
		log.Event(nil, "application error", log.Error(err))
		os.Exit(1)
	}
}
