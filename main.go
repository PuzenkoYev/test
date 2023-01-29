package main

import (
	"bufio"
	"os"
	"strings"
	"sync"
	"test/kedr"
)

var wg sync.WaitGroup

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	line := strings.Split(scanner.Text(), " ")

	for i := 0; i < len(line); i++ {
		result := kedr.UsingWordsStatistic{Path: line[i], Words: make(map[string]int), Wg: &wg}
		wg.Add(1)
		go result.ParseFileAndShowResults()
	}

	wg.Wait()
}
