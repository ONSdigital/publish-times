package main

import (
	"bufio"
	"encoding/json"
	"github.com/ONSdigital/go-ns/log"
	"github.com/daiLlew/publish-times/console"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	layout         = "2006-01-02T15:04:05.000Z"
	zebedeeRootEnv = "zebedee_root"
	zebedeeDir     = "zebedee"
	publishLogDir  = "publish-log"

	// commands
	quit            = "q"
	help            = "h"
	list            = "ls"
	showPublishTime = "pt "
)

var publishLogPath string
var publishes = make([]string, 0)

type publishTimes struct {
	PublishStartDate string `json:"publishStartDate"`
	PublishEndDate   string `json:"publishEndDate"`
}

func main() {
	console.Init()
	log.HumanReadable = true
	publishLogPath = path.Join(os.Getenv(zebedeeRootEnv), zebedeeDir, publishLogDir)

	runApp()
}

func runApp() {
	sc := bufio.NewScanner(os.Stdin)

	console.WriteHeader()
	console.WriteHelpMenu()
	console.WritePrompt()

	for sc.Scan() {
		input := sc.Text()

		processCommand(input)
		console.WritePrompt()
	}
}

func processCommand(input string) {
	if input == quit {
		exit()
	}

	if input == help {
		console.WriteHelpMenu()
		return
	}

	if input == list {
		listPublishLog()
		return
	}

	if strings.HasPrefix(input, showPublishTime) {
		args := strings.Split(input, showPublishTime)

		index, err := strconv.Atoi(strings.TrimSpace(args[1]))
		if err != nil {
			exitErr(err)
		}
		calculatePublishTime(index)
		return
	}

}

func listPublishLog() {
	files, err := filepath.Glob(publishLogPath + "/*.json")
	if err != nil {
		exitErr(err)
	}

	publishes = files
	output := make([]string, 0)
	for _, f := range publishes {
		file, err := filepath.Rel(publishLogPath, f)
		if err != nil {
			exitErr(err)
		}

		output = append(output, file)
	}

	console.WriteFiles(output)
}

func calculatePublishTime(index int) {
	filePath := publishes[index]
	fileName, err := filepath.Rel(publishLogPath, filePath)
	if err != nil {
		exitErr(err)
	}

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		exitErr(err)
	}

	var times publishTimes
	err = json.Unmarshal(b, &times)
	if err != nil {
		exitErr(err)
	}

	publishTime, err := times.calculatePublishTimes()
	if err != nil {
		exitErr(err)
	}

	console.WritePublishTime(fileName, publishTime)
}

func exit() {
	console.WriteExit()
	os.Exit(0)
}

func exitErr(err error) {
	log.Error(err, nil)
	os.Exit(1)
}

func (p publishTimes) calculatePublishTimes() (float64, error) {
	start, err := time.Parse(layout, p.PublishStartDate)
	if err != nil {
		return 0, err
	}

	end, err := time.Parse(layout, p.PublishEndDate)
	if err != nil {
		return 0, err
	}

	return end.Sub(start).Seconds(), nil
}
