package summary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	layout = "2006-01-02T15:04:05.000Z"
)

type Details struct {
	Name             string
	Time             float64
	FileCount        int
	Size             string
	PublishStartDate string `json:"publishStartDate"`
	PublishEndDate   string `json:"publishEndDate"`
}

func New(fileInfo os.FileInfo, publishLogPath string) (*Details, error) {
	collectionJSONPath := filepath.Join(publishLogPath, fileInfo.Name())

	start, end, err := getPublishStartEndTimes(collectionJSONPath)
	if err != nil {
		return nil, err
	}

	publishTime, err := collectionPublishTime(start, end)
	if err != nil {
		return nil, err
	}

	collectionDir := strings.Replace(collectionJSONPath, filepath.Ext(collectionJSONPath), "", 1)

	fileCount, err := getPublishedFileCount(collectionDir)
	if err != nil {
		return nil, err
	}

	size, err := getPublishSize(collectionDir)
	if err != nil {
		fmt.Println(3)
		return nil, err
	}

	return &Details{
		Name:      fileInfo.Name(),
		Time:      publishTime,
		FileCount: fileCount,
		Size:      size,
	}, nil
}

func getPublishStartEndTimes(path string) (start string, end string, err error) {
	var collection struct {
		PublishStartDate string `json:"publishStartDate"`
		PublishEndDate   string `json:"publishEndDate"`
	}

	var b []byte
	b, err = ioutil.ReadFile(path)
	if err != nil {
		return start, end, err
	}

	err = json.Unmarshal(b, &collection)
	if err != nil {
		return start, end, err
	}

	return collection.PublishStartDate, collection.PublishEndDate, nil
}

func getPublishedFileCount(collectionDir string) (int, error) {
	files := exec.Command("find", collectionDir, "-type", "f")
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
	outBuf := bytes.NewBuffer([]byte{})
	cmd := exec.Command("du", "-sh", collectionDir)
	cmd.Stdout = outBuf

	if err := cmd.Run(); err != nil {
		return "", err
	}
	result := outBuf.String()
	return strings.TrimSpace(strings.Replace(result, collectionDir, "", 1)), nil
}

func collectionPublishTime(start string, end string) (float64, error) {
	started, err := time.Parse(layout, start)
	if err != nil {
		return 0, err
	}

	completed, err := time.Parse(layout, end)
	if err != nil {
		return 0, err
	}

	return completed.Sub(started).Seconds(), nil
}
