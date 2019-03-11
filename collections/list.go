package collections

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

var (
	PublishLogPath string
)

func GetAll() ([]os.FileInfo, error) {
	infos, err := ioutil.ReadDir(PublishLogPath)
	if err != nil {
		return nil, err
	}

	publishes := make([]os.FileInfo, 0)
	for _, i := range infos {
		if !i.IsDir() && filepath.Ext(i.Name()) == ".json" {
			publishes = append(publishes, i)
		}
	}

	sort.SliceStable(publishes, func(i, j int) bool {
		return publishes[i].ModTime().After(publishes[j].ModTime())
	})
	return publishes, nil
}
