package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/emedvedev/enigma"
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
	ioc := float32(0)

	for ascii := 65; ascii < 91; ascii++ {
		fi := strings.Count(text, string(rune(ascii)))
		ioc += float32(fi * (fi - 1))
	}

	n := float32(len(text))
	ioc = float32(ioc / (n * (n - 1)))

	return ioc
}

func getScore(cipherText string) float32 {
	config := make([]enigma.RotorConfig, len(defaultConfig.Rotors))
	for index, rotor := range defaultConfig.Rotors {
		ring := defaultConfig.Rings[index]
		start := defaultConfig.Positions[index][0]
		config[index] = enigma.RotorConfig{ID: rotor, Start: start, Ring: ring}
	}

	plugboards := strings.Split("MS KU FY AG BN PQ HJ DI ER LW", " ") // TODO : has to be permutated
	e := enigma.NewEnigma(config, defaultConfig.Reflector, plugboards)
	plainText := e.EncodeString(string(cipherText))

	score := calculateIOC(plainText)

	return score
}

func printConfig(config Enigma) {
	fmt.Printf("\n\n Enigma configuration:")

	fmt.Printf("\n\n\t Rotor configuration:\n\t\t")
	for _, rotor := range config.Rotors {
		fmt.Printf("%v ", rotor)
	}

	fmt.Printf("\n\n\t Rotor position configuration:\n\t\t")
	for _, position := range config.Positions {
		fmt.Printf("%v ", position)
	}

	fmt.Printf("\n\n\n")
}

func copyStruct(dst *Enigma, src *Enigma) {
	dst.Reflector = src.Reflector

	dst.Rings = make([]int, len(src.Rings))
	copy(dst.Rings, src.Rings)

	dst.Positions = make([]string, len(src.Positions))
	copy(dst.Positions, src.Positions)

	dst.Rotors = make([]string, len(src.Rotors))
	copy(dst.Rotors, src.Rotors)
}
