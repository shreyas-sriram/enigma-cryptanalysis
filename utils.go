package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/emedvedev/enigma"
)

var trigrams = make(map[string]float64)

const (
	defaultPlugboard string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// readFile read a file and returns the contents as bytes
func readFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("\nError opening file: %v", filename)
		return nil, fmt.Errorf("error opening file")
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("\nError reading file: %v", filename)
		return nil, fmt.Errorf("error reading file")
	}

	return fileBytes, nil
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

func swap(a, b rune, plugboard string) string {
	fmt.Printf("\n Swap %v and %v", string(a), string(b))

	for i, letter := range plugboard {
		if letter == a {
			plugboard = plugboard[0:i] + string(b) + plugboard[i+1:]
		} else if letter == b {
			plugboard = plugboard[0:i] + string(a) + plugboard[i+1:]
		}
	}

	return plugboard
}

func getBestPlugboard(cipherText string) string {
	bestPlugboard := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	plainText := runEnigma(cipherText, bestPlugboard)
	bestIOC := calculateIOC(plainText)

	for i := 0; i < 26; i++ {

		currentIOC := bestIOC
		currentBestPlugboard := bestPlugboard

		for j := i + 1; j < 26; j++ {
			var plugboards []string

			fmt.Printf("\nCurrent plugboard: %v", bestPlugboard)

			if defaultPlugboard[j] == bestPlugboard[j] {
				plugboards = append(plugboards, swap(rune(defaultPlugboard[i]), rune(bestPlugboard[j]), bestPlugboard))
				fmt.Printf("\nPlugboard: %v", plugboards)
			} else {
				revertedPlugboard := swap(rune(defaultPlugboard[j]), rune(bestPlugboard[j]), bestPlugboard)
				fmt.Printf("\nSwapped plugboard: %v", revertedPlugboard)

				revertedPlugboard = swap(rune(defaultPlugboard[i]), rune(bestPlugboard[i]), revertedPlugboard)
				fmt.Printf("\nSwapped plugboard: %v", revertedPlugboard)

				plugboards = append(plugboards, swap(rune(defaultPlugboard[i]), rune(bestPlugboard[i]), revertedPlugboard))
				fmt.Printf("\nPlugboard: %v", plugboards)

				plugboards = append(plugboards, swap(rune(defaultPlugboard[i]), rune(bestPlugboard[j]), revertedPlugboard))
				fmt.Printf("\nPlugboards: %v", plugboards)
			}

			for _, plugboard := range plugboards {
				plainText := runEnigma(cipherText, plugboard)
				localIOC := calculateIOC(plainText)

				fmt.Printf("\nlocalIOC: %v", localIOC)

				if localIOC > currentIOC {
					currentIOC = localIOC
					currentBestPlugboard = plugboard
				}
			}
		}

		fmt.Printf("\ncurrentIOC: %v", currentIOC)
		fmt.Printf("\nbestIOC: %v", bestIOC)

		if currentIOC > bestIOC {
			bestIOC = currentIOC
			bestPlugboard = currentBestPlugboard
		}
	}

	fmt.Printf("\nbestIOC: %v", bestIOC)
	fmt.Printf("\nbestPlugboard: %v", bestPlugboard)

	return bestPlugboard
}

// runEnigma configures and run the Enigma with the given settings
// and returns the output
func runEnigma(input, plugboard string) string {
	config := make([]enigma.RotorConfig, len(defaultConfig.Rotors))
	for index, rotor := range defaultConfig.Rotors {
		ring := defaultConfig.Rings[index]
		start := defaultConfig.Positions[index][0]
		config[index] = enigma.RotorConfig{ID: rotor, Start: start, Ring: ring}
	}

	formattedPlugboard := formatPlugboard(plugboard)
	fmt.Printf("\nRawplugboard: %v", defaultPlugboard)
	fmt.Printf("\nRawplugboard: %v", plugboard)
	fmt.Printf("\nFormatted plugboard: %v", formattedPlugboard)

	e := enigma.NewEnigma(config, defaultConfig.Reflector, formattedPlugboard)
	output := e.EncodeString(string(input))

	return output
}

// func getScore(cipherText string, plugboard string) float64 {
// 	// config := make([]enigma.RotorConfig, len(defaultConfig.Rotors))
// 	// for index, rotor := range defaultConfig.Rotors {
// 	// 	ring := defaultConfig.Rings[index]
// 	// 	start := defaultConfig.Positions[index][0]
// 	// 	config[index] = enigma.RotorConfig{ID: rotor, Start: start, Ring: ring}
// 	// }

// 	// formattedPlugboard := formatPlugboard(plugboard)

// 	// // plugboards := strings.Split("MS KU FY AG BN PQ HJ DI ER LW", " ") // TODO : has to be permutated
// 	// e := enigma.NewEnigma(config, defaultConfig.Reflector, formattedPlugboard)
// 	// plainText := e.EncodeString(string(cipherText))

// 	score := calculateIOC(plainText)

// 	return score
// }

// printConfig pretty prints the configurtation of the Enigma
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

// copyStruct makes a deep copy of `Enigma` struct into another
// struct
func copyStruct(dst *Enigma, src *Enigma) {
	dst.Reflector = src.Reflector

	dst.Rings = make([]int, len(src.Rings))
	copy(dst.Rings, src.Rings)

	dst.Positions = make([]string, len(src.Positions))
	copy(dst.Positions, src.Positions)

	dst.Rotors = make([]string, len(src.Rotors))
	copy(dst.Rotors, src.Rotors)
}

// initializeTrigrams read and initializes the trigram scores from a file
// into a map structure
func initializeTrigrams(filename string) {
	fileBytes, err := readFile(filename)
	if err != nil {
		if err != nil {
			fmt.Printf("\nError reading file: %v", filename)
			os.Exit(0)
		}
	}

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

// calculateTrigram calculates the trigram score for a given text
func calculateTrigram(text string) float64 {
	score := float64(0)

	for i := 0; i < len(text)-2; i++ {
		_trigram := text[i : i+3]
		score += trigrams[_trigram]
	}

	return score
}
