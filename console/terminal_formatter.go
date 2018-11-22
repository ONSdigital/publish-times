package console

import (
	"fmt"
	. "github.com/logrusorgru/aurora"
	"os"
	"text/tabwriter"
)

var (
	helpMenu []Option
	w        *tabwriter.Writer
)

func Init() {
	w = new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	helpMenu = []Option{
		{"q", "Quit."},
		{"h", "Display the options menu."},
		{"ls", "List the collection json files in the publish log dir with an index."},
		{"pt <index>", "Get the publish time for the collection with the specified index."},
	}
}

type Option struct {
	Value       string
	Description string
}

func WriteHelpMenu() {
	fmt.Fprintf(w, "%s %s\n", Key("#>"), Val("Options:"))
	for _, item := range helpMenu {
		fmt.Fprintf(w, "- %s\t%s\t\n", Key(item.Value), Val(item.Description))
	}
	w.Flush()
}

func WriteHeader() {
	fmt.Printf(" %s\n", Key("#####################"))
	fmt.Printf(" %s %s %s\n", Key("###"), Bold(Val("Publish Times")), Key("###"))
	fmt.Printf(" %s\n", Key("#####################"))
}

func Write(value string) {
	fmt.Printf("%s", Key(value))
}

func WriteNewLine(value string) {
	fmt.Printf("%s\n", Key(value))
}

func WriteFiles(files []string) {
	fmt.Fprintf(w, "%s %s\n", Key("#>"), Val("Collection json files:"))

	for i, f := range files {
		fmt.Fprintf(w, "- %d\t%s\n", Key(i), Val(f))
	}
	w.Flush()
}

func WritePublishTime(file string, publishTime float64) {
	fmt.Fprintf(w, "%s %s\n", Key("#>"), Val("Publish time:"))
	fmt.Fprintf(w, "- %s\t%s\n", Key("file"), Val(file))
	fmt.Fprintf(w, "- %s\t%s\n", Key("time"), Val(fmt.Sprintf("%f seconds", publishTime)))
	w.Flush()
}

func Key(i interface{}) Value {
	return Bold(Green(i))
}

func Val(i interface{}) Value {
	return Cyan(i)
}

func WritePrompt() {
	Write("#> ")
}

func WriteExit() {
	fmt.Printf("%s %s\n", Key("#>"), Val("goodbye!"))
}
