package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"

	enigma_api "github.com/emedvedev/enigma"
)

var trigrams = make(map[string]float64)

const (
	defaultPlugboard string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	TRIGRAMS_FILENAME string = "english_trigrams.txt"
)

// initializeTrigrams read and initializes the trigram scores from a file
// into a map structure
func initializeTrigrams() {
	fileBytes := readFile(TRIGRAMS_FILENAME)

	file := string(fileBytes)
	fileLines := strings.Split(file, "\n")

	total := float64(0)

	for _, line := range fileLines {
		parts := strings.Split(line, " ")

		frequency, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Printf("\nError converting trigram frequency: %v", parts[1])
			os.Exit(0)
		}

		trigrams[parts[0]] = float64(frequency)
		total += float64(frequency)
	}

	for k, v := range trigrams {
		trigrams[k] = math.Log(v / total)
	}
}

// readFile read a file and returns the contents as string
func readFile(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("\nError opening file: %v", filename)
		os.Exit(0)
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("\nError reading file: %v", filename)
		os.Exit(0)
	}

	return string(fileBytes)
}

// getBestPlugboard returns the best plugboard configuration based on
// the `Index of Coincidence`
func getBestPlugboard(cipherText string, bestIOC float64) string {
	bestPlugboard := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for i := 0; i < 26; i++ {

		currentIOC := bestIOC
		currentBestPlugboard := bestPlugboard

		for j := i + 1; j < 26; j++ {
			var plugboards []string

			// fmt.Printf("\nCurrent plugboard: %v", bestPlugboard)

			if defaultPlugboard[j] == bestPlugboard[j] {
				plugboards = append(plugboards, swap(rune(defaultPlugboard[i]), rune(bestPlugboard[j]), bestPlugboard))
				// fmt.Printf("\nPlugboard: %v", plugboards)
			} else {
				revertedPlugboard := swap(rune(defaultPlugboard[j]), rune(bestPlugboard[j]), bestPlugboard)
				// fmt.Printf("\nSwapped plugboard: %v", revertedPlugboard)

				if defaultPlugboard[i] != bestPlugboard[i] {
					revertedPlugboard = swap(rune(defaultPlugboard[i]), rune(bestPlugboard[i]), revertedPlugboard)
					// fmt.Printf("\nSwapped plugboard: %v", revertedPlugboard)
				}

				plugboards = append(plugboards, swap(rune(defaultPlugboard[i]), rune(bestPlugboard[i]), revertedPlugboard))
				// fmt.Printf("\nPlugboard: %v", plugboards)

				plugboards = append(plugboards, swap(rune(defaultPlugboard[i]), rune(bestPlugboard[j]), revertedPlugboard))
				// fmt.Printf("\nPlugboards: %v", plugboards)
			}

			for _, plugboard := range plugboards {
				plainText := runEnigmaWithPlugboard(cipherText, plugboard)
				localIOC := calculateIOC(plainText)

				// fmt.Printf("\nlocalIOC: %v", localIOC)

				if localIOC > currentIOC {
					currentIOC = localIOC
					currentBestPlugboard = plugboard
				}
			}
		}

		// fmt.Printf("\ncurrentIOC: %v", currentIOC)
		// fmt.Printf("\nbestIOC: %v", bestIOC)

		if currentIOC > bestIOC {
			bestIOC = currentIOC
			bestPlugboard = currentBestPlugboard
		}
	}

	// fmt.Printf("\nbestIOC: %v", bestIOC)
	// fmt.Printf("\nbestPlugboard: %v", bestPlugboard)

	return bestPlugboard
}

// runEnigma configures and run the enigma with the given settings
// and returns the output
func runEnigma(input string, config enigma) string {
	newConfig := make([]enigma_api.RotorConfig, len(config.rotors))
	for index, rotor := range config.rotors {
		ring := config.rings[index]
		start := config.positions[index][0]
		newConfig[index] = enigma_api.RotorConfig{ID: rotor, Start: start, Ring: ring}
	}

	formattedPlugboard := formatPlugboard(config.plugboard)
	// fmt.Printf("\nRawplugboard: %v", defaultPlugboard)
	// fmt.Printf("\nRawplugboard: %v", plugboard)
	// fmt.Printf("\nFormatted plugboard: %v", formattedPlugboard)

	e := enigma_api.NewEnigma(newConfig, config.reflector, formattedPlugboard)
	output := e.EncodeString(input)

	return output
}

// runEnigmaWithPlugboard is a wrapper for `runEnigma()` and takes in
// plugboard configuration instead of `enigma` struct
func runEnigmaWithPlugboard(input, plugboard string) string {
	var config enigma
	copyStruct(&config, &currentConfig)

	config.plugboard = plugboard

	return runEnigma(input, config)
}

// formatPlugboard converts to raw plugboard string into the format
// expected by `github.com/emedvedev/enigma` (string array)
func formatPlugboard(rawPlugboard string) []string {
	formattedPlugboard := make([]string, 0)

	for i, letter := range rawPlugboard {
		if string(rune(i+65)) < string(letter) {
			formattedPlugboard = append(formattedPlugboard, string(rune(i+65))+string(letter))
		}
	}

	return formattedPlugboard
}

// swap swaps two characters in a given string and returns the new string
func swap(a, b rune, plugboard string) string {
	// fmt.Printf("\n Swap %v and %v", string(a), string(b))

	for i, letter := range plugboard {
		if letter == a {
			plugboard = plugboard[0:i] + string(b) + plugboard[i+1:]
		} else if letter == b {
			plugboard = plugboard[0:i] + string(a) + plugboard[i+1:]
		}
	}

	return plugboard
}

// copyStruct makes a deep copy of `Enigma` struct into another
// struct
func copyStruct(dst *enigma, src *enigma) {
	dst.reflector = src.reflector

	dst.rings = make([]int, len(src.rings))
	copy(dst.rings, src.rings)

	dst.positions = make([]string, len(src.positions))
	copy(dst.positions, src.positions)

	dst.rotors = make([]string, len(src.rotors))
	copy(dst.rotors, src.rotors)

	dst.plugboard = src.plugboard
}

// calculateIOC calculates the `Index of Coincidence` for a given text
func calculateIOC(text string) float64 {
	ioc := float64(0)

	for ascii := 65; ascii < 91; ascii++ {
		fi := strings.Count(text, string(rune(ascii)))
		ioc += float64(fi * (fi - 1))
	}

	n := float64(len(text))
	ioc = float64(ioc / (n * (n - 1)))

	return ioc
}

// calculateTrigram calculates the trigram score for a given text
func calculateTrigram(text string) float64 {
	score := float64(0)

	for i := 0; i < len(text)-2; i++ {
		_trigram := text[i : i+3]
		score += trigrams[_trigram]
	}

	return score
}

// printConfig pretty prints the configuration of the enigma
func printConfig(config enigma) {
	fmt.Printf("\n\n Enigma configuration:")

	fmt.Printf("\n\n\t Rotor configuration:\n\t\t")
	for _, rotor := range config.rotors {
		fmt.Printf("%v ", rotor)
	}

	fmt.Printf("\n\n\t Rotor position configuration:\n\t\t")
	for _, position := range config.positions {
		fmt.Printf("%v ", position)
	}

	fmt.Printf("\n\n\t Plugboard configuration:\n")
	fmt.Printf("\t\t%v", config.plugboard)
	fmt.Printf("\n\t\t%v", formatPlugboard(config.plugboard))

	fmt.Printf("\n\n\n")
}

// printExpected prints the configuration of the enigma in
// the format required by the assignment
func printExpected(config enigma) {
	for _, rotor := range config.rotors {
		fmt.Printf("%v ", rotor)
	}

	fmt.Printf("\n")
	for _, position := range config.positions {
		fmt.Printf("%v ", position)
	}

	fmt.Printf("\n")
	for _, plugboardItem := range formatPlugboard(config.plugboard) {
		fmt.Printf("%v ", plugboardItem)
	}
}
