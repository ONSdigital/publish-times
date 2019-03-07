package main

import (
	"bufio"
	"flag"
	"github.com/ONSdigital/log.go/log"
	"github.com/ONSdigital/publish-times/console"
	"github.com/ONSdigital/publish-times/summary"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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
	rangeRegex      = `range\(\s*\d+\s*,\s*\d+\s*\)`

	rangePrefix = "range("
	rangeSuffix = ")"
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
	loadAllPublishedCollections()

	console.WriteHeader()
	console.WriteHelpMenu()
	console.WritePrompt()

	for sc.Scan() {
		input := sc.Text()

		err := processCommand(input)
		if err != nil {
			handleInputErr(err)
		}
		console.WritePrompt()
	}
}

func processCommand(input string) error {
	if input == quit {
		exit()
	}

	if input == clearTerminal {
		console.Clear()
		return nil
	}

	if input == help {
		console.WriteHelpMenu()
		return nil
	}

	if strings.HasPrefix(input, list) {
		listCollections()
		return nil
	}

	isRange, err := regexp.MatchString(rangeRegex, input)
	if err != nil {
		return err
	}

	if isRange {
		return rangeCollections(input)
	}

	if strings.HasPrefix(input, showPublishTime) {
		return showPublishSummary(input)
	}
	return nil
}

func handleInputErr(err error) {
	invalidInput, ok := err.(InputErr)
	if ok {
		console.Warn(invalidInput.Message)
	} else {
		exitErr(err)
	}
}

func listCollections() {
	loadAllPublishedCollections()
	console.WriteFiles(publishes)
}

func rangeCollections(input string) error {
	parsed := strings.Replace(input, rangePrefix, "", -1)
	parsed = strings.Replace(parsed, rangeSuffix, "", -1)

	values := strings.Split(parsed, ",")
	start, err := strconv.Atoi(strings.TrimSpace(values[0]))
	if err != nil {
		return err
	}

	end, err := strconv.Atoi(strings.TrimSpace(values[1]))
	if err != nil {
		return err
	}

	if start < 0 {
		return InputErr{"Start index cannot be less than 0"}
	}

	if start > end {
		return InputErr{"Start index cannot be greater than end index"}
	}

	if end > len(publishes) {
		return InputErr{"End index greater than total number of published collections"}
	}

	console.WriteRange(start, end, publishes)
	return nil
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
		summary, err := summary.New(index, collection, publishLogPath)
		if err != nil {
			return err
		}

		publishInfos = append(publishInfos, summary)
	}
	console.WritePublishSummaries(publishInfos)
	return nil
}

func loadAllPublishedCollections() {
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

	if len(publishes) > 100 {
		// TODO for now return the last 100
		// Consider adding some pagination feature in future
		publishes = publishes[:100]
	}
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
