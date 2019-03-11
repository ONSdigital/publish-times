package input

import (
	"bufio"
	"github.com/ONSdigital/publish-times/collections"
	"github.com/ONSdigital/publish-times/console"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
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

type Err struct {
	Message string
}

func (e Err) Error() string {
	return e.Message
}

func Listen() (err error) {
	sc := bufio.NewScanner(os.Stdin)

	console.WriteHeader()
	console.WriteHelpMenu()
	console.WritePrompt()

	for sc.Scan() {
		input := sc.Text()

		if input == quit {
			break
		}
		if err = processCommand(input); err != nil {
			break
		}
	}

	invalidInput, ok := err.(Err)
	if ok {
		console.Warn(invalidInput.Message)
	} else {
		return err
	}

	return nil
}

func processCommand(input string) error {
	if input == clearTerminal {
		console.Clear()
		return nil
	}

	if input == help {
		console.WriteHelpMenu()
		return nil
	}

	if strings.HasPrefix(input, list) {
		published, err := collections.GetAll()
		if err != nil {
			return err
		}
		console.WriteFiles(published)
		return nil
	}

	isRange, err := regexp.MatchString(rangeRegex, input)
	if err != nil {
		return err
	}

	if isRange {
		return listRange(input)
	}

	if strings.HasPrefix(input, showPublishTime) {
		return showPublishSummary(input)
	}
	return nil
}

func listRange(input string) error {
	parsed := strings.Replace(input, rangePrefix, "", -1)
	parsed = strings.Replace(parsed, rangeSuffix, "", -1)

	values := strings.Split(parsed, ",")
	start, err := strconv.Atoi(strings.TrimSpace(values[0]))
	if err != nil {
		return Err{Message: err.Error()}
	}

	end, err := strconv.Atoi(strings.TrimSpace(values[1]))
	if err != nil {
		return Err{Message: err.Error()}
	}

	return collections.Range(start, end)
}

func showPublishSummary(input string) error {
	input = strings.Replace(input, showPublishTime, "", -1)
	args := strings.Split(input, ",")

	publishInfos := make([]*collections.PublishSummary, 0)

	all, err := collections.GetAll()
	if err != nil {
		return err
	}

	for _, i := range args {
		index, err := strconv.Atoi(strings.TrimSpace(i))
		if err != nil {
			return Err{Message: "Invalid index value"}
		}

		if index < 0 {
			return Err{Message: "Index must be greater than 0"}
		}

		if index >= len(all) {
			return Err{Message: "Index cannot be greater than the number of published collections"}
		}

		collection := all[index]
		summary, err := collections.NewPublishSummary(index, collection)
		if err != nil {
			return err
		}

		publishInfos = append(publishInfos, summary)
	}
	console.WritePublishSummaries(publishInfos)
	return nil
}
