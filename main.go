package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"github.com/ONSdigital/log.go/log"
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
	layout           = "2006-01-02T15:04:05.000Z"
	publishLogDirEnv = "PUBLISH_LOG_DIR"
	zebedeeRootEnv   = "zebedee_root"
	zebedeeDir       = "zebedee"
	publishLogDir    = "publish-log"

	// commands
	quit            = "q"
	clearTerminal   = "clear"
	help            = "h"
	list            = "ls"
	showPublishTime = "pt "
)

var publishLogPath string
var publishes = make([]os.FileInfo, 0)

type publishedCollection struct {
	PublishStartDate string `json:"publishStartDate"`
	PublishEndDate   string `json:"publishEndDate"`
}

func main() {
	console.Init()
	log.Namespace = "publish-times"

	publishLogFlag := flag.String("publishLogDir", "", "override - use this value instead of the env config")
	flag.Parse()

	if len(*publishLogFlag) > 0 {
		log.Event(nil, "setting publish log path from cmd flag", nil)
		publishLogPath = *publishLogFlag
	} else {
		publishLogPath = os.Getenv(publishLogDirEnv)
		log.Event(nil, "setting publish log path from os env var", log.Data{publishLogDirEnv: publishLogPath})
	}

	log.Event(nil, "configuration", log.Data{"publishLogPath": publishLogPath})
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
		showPublishDetails(input)
	}
}

func showPublishDetails(input string) {
	args := strings.Split(input, showPublishTime)

	index, err := strconv.Atoi(strings.TrimSpace(args[1]))
	if err != nil {
		exitErr(err)
	}

	collection := publishes[index]
	publishTime := calculatePublishTime(collection)
	fileCount, err := getPublishedFileCount(collection.Name())
	if err != nil {
		exitErr(err)
	}

	size, err := getPublishSize(collection.Name())
	if err != nil {
		exitErr(err)
	}
	console.WritePublishTime(collection.Name(), publishTime, fileCount, size)
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

	if len(publishes) > 15 {
		// TODO for now return the last 15
		// Consider adding some pagination feature in future
		publishes = publishes[:15]
	}
}

func calculatePublishTime(fileInfo os.FileInfo) float64 {
	filename := filepath.Join(publishLogPath, fileInfo.Name())

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		exitErr(err)
	}

	var publishedCol publishedCollection
	err = json.Unmarshal(b, &publishedCol)
	if err != nil {
		exitErr(err)
	}

	publishTime, err := publishedCol.calculatePublishTimes()
	if err != nil {
		exitErr(err)
	}

	return publishTime
}

func getPublishedFileCount(collectionName string) (int, error) {
	files := exec.Command("find", getCollectionDir(collectionName), "-type", "f")
	files.Stderr = os.Stderr
	fOut, err := files.StdoutPipe()
	if err != nil {
		return 0, err
	}

	buff := bytes.NewBuffer([]byte{})

	defer fOut.Close()
	count := exec.Command("wc", "-l")
	count.Stdin = fOut
	count.Stdout = buff

	files.Start()
	count.Start()
	files.Wait()
	count.Wait()

	val := strings.TrimSpace(buff.String())
	return strconv.Atoi(val)
}

func getPublishSize(collectionDir string) (string, error) {
	dir := getCollectionDir(collectionDir)

	outBuf := bytes.NewBuffer([]byte{})
	cmd := exec.Command("du", "-sh", dir)
	cmd.Stdout = outBuf

	if err := cmd.Run(); err != nil {
		return "", err
	}
	result := outBuf.String()
	return strings.TrimSpace(strings.Replace(result, dir, "", 1)), nil
}

func getCollectionDir(collectionName string) string {
	collectionJSON := path.Join(publishLogPath, collectionName)
	return strings.Replace(collectionJSON, filepath.Ext(collectionJSON), "", 1)
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

func (p publishedCollection) calculatePublishTimes() (float64, error) {
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
