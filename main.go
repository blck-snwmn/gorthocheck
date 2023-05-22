package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

//go:embed patterns.json
var embeddedFiles []byte

type WordPattern struct {
	Expected string   `json:"expected"`
	Patterns []string `json:"patterns"`
}

func readPatterns() (map[string]string, error) {
	var wordPatterns []WordPattern
	err := json.Unmarshal(embeddedFiles, &wordPatterns)
	if err != nil {
		return nil, err
	}

	patterns := make(map[string]string)
	for _, wp := range wordPatterns {
		for _, pattern := range wp.Patterns {
			patterns[pattern] = wp.Expected
		}
	}

	return patterns, nil
}

func checkFile(filename string, patterns map[string]string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		for pattern, expected := range patterns {
			if strings.Contains(line, pattern) {
				fmt.Printf("Found typo '%s' at line %d, did you mean '%s'?\n", pattern, lineNumber, expected)
			}
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide a file to check.")
		os.Exit(1)
	}

	filename := os.Args[1]
	patterns, err := readPatterns()
	if err != nil {
		fmt.Println("Error reading patterns:", err)
		os.Exit(1)
	}

	err = checkFile(filename, patterns)
	if err != nil {
		fmt.Println("Error checking file:", err)
		os.Exit(1)
	}
}
