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

	IOC     string = "ioc"
	TRIGRAM string = "trigram"
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

// doHillclimb performs hillclimb attack
func doHillclimb(cipherText string, baseIOC float64) (string, float64) {

	// Hillclimb based on IOC score
	currentPlugboard, _ := getBestPlugboard(cipherText, baseIOC, defaultPlugboard, IOC)
	currentConfig.plugboard = currentPlugboard

	// Get base trigram score
	plainText := runEnigma(cipherText, currentConfig)
	currentTrigramScore := calculateTrigram(plainText)

	// Hillclimb based on Trigram score
	bestPlugboard, bestTrigramScore := getBestPlugboard(cipherText, currentTrigramScore, currentPlugboard, TRIGRAM)

	return bestPlugboard, bestTrigramScore
}

// getBestPlugboard returns the best plugboard configuration based on
// the `Index of Coincidence` or `Trigram` score
func getBestPlugboard(cipherText string, bestScore float64, bestPlugboard string, scoreType string) (string, float64) {
	seen := make(map[string]bool)
	seen[bestPlugboard] = true

	for i := 0; i < 26; i++ {
		currentScore := bestScore
		currentBestPlugboard := bestPlugboard

		// fmt.Printf("\nChecking for: %v", string("ABCDEFGHIJKLMNOPQRSTUVWXYZ"[i]))

		for j := i + 1; j < 26; j++ {
			var plugboards []string

			// fmt.Printf("\nDefault plugboard: %v", defaultPlugboard)
			// fmt.Printf("\nCurrent plugboard: %v", bestPlugboard)

			if defaultPlugboard[i] != bestPlugboard[i] && defaultPlugboard[j] != bestPlugboard[j] { // both `i` and `j` are already swapped
				// revert the positions
				revertedPlugboard := swap(rune(defaultPlugboard[i]), rune(bestPlugboard[i]), bestPlugboard)
				revertedPlugboard = swap(rune(defaultPlugboard[j]), rune(bestPlugboard[j]), revertedPlugboard)

				// make the combinations
				// plugboards = append(plugboards, swap(rune(defaultPlugboard[i]), rune(bestPlugboard[i]), revertedPlugboard))
				plugboards = append(plugboards, swap(rune(defaultPlugboard[i]), rune(bestPlugboard[j]), revertedPlugboard))
				// plugboards = append(plugboards, swap(rune(defaultPlugboard[i]), rune(defaultPlugboard[j]), revertedPlugboard))
				// plugboards = append(plugboards, swap(rune(defaultPlugboard[j]), rune(bestPlugboard[i]), revertedPlugboard))
				// plugboards = append(plugboards, swap(rune(defaultPlugboard[j]), rune(bestPlugboard[j]), revertedPlugboard))
			} else if defaultPlugboard[i] != bestPlugboard[i] { // only `i` is already swapped
				// revert the positions
				revertedPlugboard := swap(rune(defaultPlugboard[i]), rune(bestPlugboard[i]), bestPlugboard)

				// make the combinations
				plugboards = append(plugboards, swap(rune(defaultPlugboard[i]), rune(defaultPlugboard[j]), revertedPlugboard))
			} else if defaultPlugboard[j] != bestPlugboard[j] { // only `j` is already swapped
				// revert the positions
				revertedPlugboard := swap(rune(defaultPlugboard[j]), rune(bestPlugboard[j]), bestPlugboard)

				// make the combinations
				plugboards = append(plugboards, swap(rune(defaultPlugboard[i]), rune(defaultPlugboard[j]), revertedPlugboard))
				plugboards = append(plugboards, swap(rune(defaultPlugboard[i]), rune(bestPlugboard[j]), revertedPlugboard))
			} else { // both `i` and `j` are in their original positions
				plugboards = append(plugboards, swap(rune(defaultPlugboard[i]), rune(defaultPlugboard[j]), bestPlugboard))
			}

			// fmt.Printf("\nPlugboards: %v", plugboards)

			for _, plugboard := range plugboards {

				// if plugboard combination is already seen, then ignore
				if _, ok := seen[plugboard]; ok {
					continue
				}

				seen[plugboard] = true

				plainText := runEnigmaWithPlugboard(cipherText, plugboard)
				localScore := float64(0)

				if scoreType == IOC {
					localScore = calculateIOC(plainText)
				} else {
					localScore = calculateTrigram(plainText)
				}

				if localScore > currentScore {
					currentScore = localScore
					currentBestPlugboard = plugboard
				}
			}
		}

		if currentScore > bestScore {
			bestScore = currentScore
			bestPlugboard = currentBestPlugboard
		}
	}

	return bestPlugboard, bestScore
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
	for i, rotor := range config.rotors {
		if i != len(config.rotors)-1 {
			fmt.Printf("%v ", rotor)
		} else {
			fmt.Printf("%v", rotor)
		}
	}

	fmt.Printf("\n")
	for i, position := range config.positions {
		if i != len(config.positions)-1 {
			fmt.Printf("%v ", position)
		} else {
			fmt.Printf("%v", position)
		}
	}

	formattedPlugboard := formatPlugboard(config.plugboard)

	fmt.Printf("\n")
	for i, plugboardItem := range formattedPlugboard {
		if i != len(formattedPlugboard)-1 {
			fmt.Printf("%v ", plugboardItem)
		} else {
			fmt.Printf("%v", plugboardItem)
		}
	}

	fmt.Printf("\n")
}
