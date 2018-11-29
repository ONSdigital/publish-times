package main

import (
	"bufio"
	"encoding/json"
	"github.com/ONSdigital/go-ns/log"
	"github.com/daiLlew/publish-times/console"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
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
	clearTerminal   = "clear"
	help            = "h"
	list            = "ls"
	showPublishTime = "pt "
)

var publishLogPath string
var publishes = make([]os.FileInfo, 0)

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
	loadCollectionFiles()

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

	if input == clearTerminal {
		clear()
		return
	}

	if input == help {
		console.WriteHelpMenu()
		return
	}

	if input == list {
		loadCollectionFiles()
		console.WriteFiles(publishes)
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

func loadCollectionFiles() {
	infos, err := ioutil.ReadDir(publishLogPath)
	if err != nil {
		exitErr(err)
	}

	publishes = make([]os.FileInfo, 0)
	for _, i := range infos {
		if !i.IsDir() && filepath.Ext(i.Name()) == ".json" {
			publishes = append(publishes, i)
		}
	}

	sort.SliceStable(publishes, func(i, j int) bool {
		return publishes[i].ModTime().After(publishes[j].ModTime())
	})
}

func calculatePublishTime(index int) {
	fileInfo := publishes[index]
	filename := filepath.Join(publishLogPath, fileInfo.Name())

	b, err := ioutil.ReadFile(filename)
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

	console.WritePublishTime(fileInfo.Name(), publishTime)
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
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
