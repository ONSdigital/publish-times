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
	fmt.Fprintf(w, "%s %s\n", Key("#>"), ValAplha("Options:"))
	for _, item := range helpMenu {
		fmt.Fprintf(w, "- %s\t%s\t\n", Key(item.Value), ValAplha(item.Description))
	}
	w.Flush()
}

func WriteHeader() {
	fmt.Println()
	fmt.Printf(" %s\n", Key("#####################"))
	fmt.Printf(" %s %s %s\n", Key("###"), Bold(ValAplha("Publish Times")), Key("###"))
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
	fmt.Fprintf(w, "%s\t %s\t %s\t\n", Key("Index"), ValAplha("Collection"), ValBeta("Date"))

	for i, f := range files {
		lastMod := f.ModTime().Format("Mon Jan _2 15:04:05 2006")
		fmt.Fprintf(w, "- %d\t %s\t %s\t\n", Key(i), ValAplha(f.Name()), ValBeta(lastMod))
	}
	w.Flush()
}

func WritePublishTime(file string, publishTime float64, fileCount int, totalSize string) {
	fmt.Fprintf(w, "%s %s\t%s\t%s\t%s\t\n",
		Key("- "),
		ValAplha("Collection"),
		ValBeta("Time (seconds)"),
		ValAplha("Count"),
		ValBeta("Size"),
	)

	// fmt.Fprintf(w, "%s %s\t%s\t%d\t%s\t\n", Key("- "), ValAplha(file), ValBeta(fmt.Sprintf("%f", publishTime)), ValAplha(fileCount), ValBeta(totalSize))
	fmt.Fprintf(w, "%s %s\t%s\t%d\t%s\t\n", Key("- "), ValAplha(file), ValBeta(fmt.Sprintf("%f", publishTime)), ValAplha(fileCount), ValBeta(totalSize))
	w.Flush()
}

func Key(i interface{}) Value {
	return Bold(Green(i))
}

func ValAplha(i interface{}) Value {
	return Cyan(i)
}

func ValBeta(i interface{}) Value {
	return Magenta(i)
}

func WritePrompt() {
	Write("#> ")
}

func WriteExit() {
	fmt.Printf("%s %s\n", Key("#>"), ValAplha("goodbye!"))
}
