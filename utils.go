package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// readFile read a file and returns the contents as bytes
func readFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("\nError opening file")
		return nil, fmt.Errorf("error opening file")
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("\nError reading file")
		return nil, fmt.Errorf("error reading file")
	}

	return fileBytes, nil
}

// calculateIOC calculates the Index of Coincidence for a given text
func calculateIOC(text string) float32 {
	ioc := float32(0.0)

	for ascii := 65; ascii < 91; ascii++ {
		fi := strings.Count(text, string(rune(ascii)))
		ioc += float32(fi * (fi - 1))
	}

	n := float32(len(text))
	ioc = float32(ioc / (n * (n - 1)))

	return ioc
}
