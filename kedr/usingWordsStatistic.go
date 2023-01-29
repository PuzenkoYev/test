package kedr

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type StatisticOperations interface {
	ParseFile()
	ParseFileAndShowResults()
	AddWord(word string)
	ShowResults()
}

type UsingWordsStatistic struct {
	Path  string
	Words map[string]int
	Wg    *sync.WaitGroup
}

func (u *UsingWordsStatistic) AddWord(word string) {
	rg := regexp.MustCompile(`[[:punct:]]`)
	word = rg.ReplaceAllString(word, "")

	if len(word) <= 9 {
		u.Words[strconv.Itoa(len(word))]++
	} else {
		u.Words[">9"]++
	}
}

func (u *UsingWordsStatistic) ParseFile() {
	if isFileNotExist(u.Path) {
		return
	}

	res, _ := isFileBinary(u.Path)
	if res {
		return
	}
	file, _ := os.Open(u.Path)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		u.AddWord(scanner.Text())
	}
}

func (u *UsingWordsStatistic) ParseFileAndShowResults() {
	defer u.Wg.Done()

	u.ParseFile()
	u.ShowResults()
}

func (u *UsingWordsStatistic) ShowResults() {
	res, _ := isFileBinary(u.Path)
	if isFileNotExist(u.Path) {
		return
	}
	if res {
		fmt.Println(u.Path, " file is binary!")
		return
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "%s:	", u.Path)
	for _, k := range getSortedKeys(u.Words) {
		fmt.Fprintf(&builder, "[%s]=[%d],  ", k, u.Words[k])
	}

	fmt.Println(builder.String())
}

func getSortedKeys(words map[string]int) []string {
	var keys []string
	for k := range words {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func isFileBinary(name string) (bool, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return false, err
	}

	for _, b := range data {
		if b > 127 {
			return true, nil
		}
	}

	return false, nil
}

func isFileNotExist(name string) bool {
	_, err := os.Stat(name)
	if !os.IsNotExist(err) {
		return false
	}

	return true
}
