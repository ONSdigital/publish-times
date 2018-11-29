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
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	helpMenu = []Option{
		{"q", "Quit."},
		{"clear", "Clear the terminal"},
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
	fmt.Println()
	fmt.Printf(" %s\n", Key("#####################"))
	fmt.Printf(" %s %s %s\n", Key("###"), Bold(Val("Publish Times")), Key("###"))
	fmt.Printf(" %s\n", Key("#####################"))
	fmt.Println()
}

func Write(value string) {
	fmt.Printf("%s", Key(value))
}

func WriteNewLine(value string) {
	fmt.Printf("%s\n", Key(value))
}

func WriteFiles(files []os.FileInfo) {
	fmt.Fprintf(w, "%s\t %s\t %s\t\n", Key("index"), Val("filename"), Val2("date published"))

	for i, f := range files {
		lastMod := f.ModTime().Format("Mon Jan _2 15:04:05 2006")
		fmt.Fprintf(w, "- %d\t %s\t %s\t\n", Key(i), Val(f.Name()), Val2(lastMod))
	}
	w.Flush()
}

func WritePublishTime(file string, publishTime float64) {
	fmt.Fprintf(w, "%s %s\t%s\t\n", Key("#>"), Val("Filename"), Val2("Seconds"))
	fmt.Fprintf(w, "%s %s\t%s\t\n", Key("#>"), Val(file), Val2(fmt.Sprintf("%f", publishTime)))
	w.Flush()
}

func Key(i interface{}) Value {
	return Bold(Green(i))
}

func Val(i interface{}) Value {
	return Cyan(i)
}

func Val2(i interface{}) Value {
	return Magenta(i)
}

func WritePrompt() {
	Write("#> ")
}

func WriteExit() {
	fmt.Printf("%s %s\n", Key("#>"), Val("goodbye!"))
}
