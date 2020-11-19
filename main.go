package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {

	targets := parseFile("./targets.txt")
	fmt.Printf("已载入 %d 个目标\n", len(targets))

	ch := make(chan Target)
	defer func() {
		close(ch)
	}()

	for _, target := range targets {
		go getCertHost(ch, target)
	}

	for x := 0; x < len(targets); x++ {
		fmt.Println(<-ch)
	}

}

func parseFile(path string) []Target {
	bytes, err := ioutil.ReadFile(path)
	check(err)
	lines := strings.Split(string(bytes), "\n")
	var targets []Target
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		targets = append(targets, ParserTarget(line))

	}
	return targets
}
