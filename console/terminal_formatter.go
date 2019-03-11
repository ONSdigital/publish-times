package console

import (
	"fmt"
	"github.com/ONSdigital/publish-times/collections"
	. "github.com/logrusorgru/aurora"
	"os"
	"os/exec"
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
		{"pt <i> | pt <i>, <j>... <n>", "Get the publish time for the collection(s) with the specified index/indices."},
		{"range(i, j)", "List the collections from i to j"},
	}
}

type PublishInfo struct {
	Name      string
	Time      float64
	FileCount int
	Size      string
}

type Option struct {
	Value       string
	Description string
}

func Clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
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

func Warn(i interface{}) {
	fmt.Printf("%s %s\n", Key("#>"), Red(i))
}

func WriteFiles(files []os.FileInfo) {
	fmt.Fprintf(w, "%s\t %s\t %s\t\n", Key("Index"), ValAplha("Collection"), ValBeta("Date"))

	for i, f := range files {
		lastMod := f.ModTime().Format("Mon Jan _2 15:04:05 2006")
		fmt.Fprintf(w, "- %d\t %s\t %s\t\n", Key(i), ValAplha(f.Name()), ValBeta(lastMod))
	}
	w.Flush()
}

func WriteRange(start int, end int, files []os.FileInfo) {
	fmt.Fprintf(w, "%s\t %s\t %s\t\n", Key("Index"), ValAplha("Collection"), ValBeta("Date"))

	for i := start; i <= end; i++ {
		f := files[i]
		lastMod := f.ModTime().Format("Mon Jan _2 15:04:05 2006")
		fmt.Fprintf(w, "- %d\t %s\t %s\t\n", Key(i), ValAplha(f.Name()), ValBeta(lastMod))
	}
	w.Flush()
}

func WritePublishSummaries(infos []*collections.PublishSummary) {
	fmt.Fprintf(w, "%s %s\t%s\t%s\t%s\t%s\t\n",
		Key("- "),
		ValAplha("Index"),
		ValBeta("Collection"),
		ValAplha("Time (seconds)"),
		ValBeta("Count"),
		ValAplha("Size"),
	)

	for _, info := range infos {
		fmt.Fprintf(w, "%s %d\t%s\t%s\t%d\t%s\t\n",
			Key("- "),
			ValAplha(info.Index),
			ValBeta(info.Name),
			ValAplha(fmt.Sprintf("%f", info.Time)),
			ValBeta(info.FileCount),
			ValAplha(info.Size),
		)
	}

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
