package main

import (
	"bufio"
	"flag"
	"github.com/ONSdigital/log.go/log"
	"github.com/daiLlew/publish-times/console"
	"github.com/daiLlew/publish-times/summary"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	publishLogDirEnv = "PUBLISH_LOG_DIR"

	// commands
	quit            = "q"
	clearTerminal   = "clear"
	help            = "h"
	list            = "ls"
	showPublishTime = "pt "
)

var (
	publishLogPath string
	publishes      = make([]os.FileInfo, 0)
)

type InputErr struct {
	Message string
}

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

	begin()
}

func begin() {
	sc := bufio.NewScanner(os.Stdin)
	loadPublishedCollections()

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
		loadPublishedCollections()
		console.WriteFiles(publishes)
		return
	}

	if strings.HasPrefix(input, showPublishTime) {
		if err := showPublishSummary(input); err != nil {
			invalidInput, ok := err.(InputErr)
			if ok {
				console.Warn(invalidInput.Message)
			} else {
				exitErr(err)
			}
		}
	}
}

func showPublishSummary(input string) error {
	input = strings.Replace(input, showPublishTime, "", -1)
	args := strings.Split(input, ",")

	publishInfos := make([]*summary.Details, 0)

	for _, i := range args {
		index, err := strconv.Atoi(strings.TrimSpace(i))
		if err != nil {
			return InputErr{Message: "Invalid index value"}
		}

		if index < 0 {
			return InputErr{Message: "Index must be greater than 0"}
		}

		if index > len(publishes) {
			return InputErr{Message: "Index cannot be greater than the number of published collections"}
		}

		collection := publishes[index]
		summary, err := summary.New(collection, publishLogPath)
		if err != nil {
			return err
		}

		publishInfos = append(publishInfos, summary)
	}
	console.WritePublishSummaries(publishInfos)
	return nil
}

func loadPublishedCollections() {
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

	if len(publishes) > 15 {
		// TODO for now return the last 15
		// Consider adding some pagination feature in future
		publishes = publishes[:15]
	}
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
	log.Event(nil, "application error", log.Error(err))
	os.Exit(1)
}

func (e InputErr) Error() string {
	return e.Message
}
