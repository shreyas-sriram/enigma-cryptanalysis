// enigma "LOTSOFPEOPLEWOR" --rotors "Gamma VI IV III" --rings="1 1 1 16" --position "D A B Q" --plugboard "MS KU FY AG BN PQ HJ DI ER LW" --reflector "C-thin"
// Program uses `github.com/emedvedev/enigma` as an external API to test enigma configurations

package main

import (
	"fmt"
	"math"
	"os"
)

// enigma stores the enigma configuration
type enigma struct {
	reflector string
	rings     []int
	positions []string
	rotors    []string
	plugboard string
}

var currentConfig = enigma{
	reflector: "C-thin",
	rings:     []int{1, 1, 1, 16},
	positions: []string{"_", "_", "B", "Q"},
	rotors:    []string{"_", "_", "IV", "III"},
	plugboard: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
}

var bestConfig = enigma{
	reflector: "C-thin",
	rings:     []int{1, 1, 1, 16},
	positions: []string{"", "", "B", "Q"},
	rotors:    []string{"", "", "IV", "III"},
	plugboard: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
}

const (
	TRIGRAMS_FILENAME = "english_trigrams.txt"
)

var (
	rotorsPositionOne = []string{"Beta", "Gamma", "I", "II", "V", "VI"}
	rotorsPositionTwo = []string{"I", "II", "V", "VI"}
)

func main() {
	fmt.Printf("\n --- enigma cryptanalysis --- \n")

	if len(os.Args) != 2 {
		fmt.Printf("\n [-] Cipher text file not provided, exiting\n")
		os.Exit(0)
	}

	cipherText := readFile(os.Args[1])

	fmt.Printf("\n [+] Read cipher text :\n %v", cipherText)

	bestTrigramScore := math.Inf(-1)

	initializeTrigrams(TRIGRAMS_FILENAME)

	for _, rotorPositionOne := range rotorsPositionOne { // rotor at 1st place
		for _, rotorPositionTwo := range rotorsPositionTwo { // rotor at 2nd place
			if rotorPositionOne == rotorPositionTwo {
				continue
			}

			currentConfig.rotors[0] = rotorPositionOne
			currentConfig.rotors[1] = rotorPositionTwo

			fmt.Printf("\n\n [+] Trying enigma configuration:")
			printConfig(currentConfig)

			for i := 0; i < 26; i++ { // starting position of rotor at 1st place
				for j := 0; j < 26; j++ { // starting position of rotor at 2nd place
					currentConfig.positions[0] = string(rune(i + 65))
					currentConfig.positions[1] = string(rune(j + 65))

					currentPlugboard := getBestPlugboard(cipherText)
					currentConfig.plugboard = currentPlugboard

					plainText := runEnigma(cipherText, currentConfig)

					currentTrigramScore := calculateTrigram(plainText)

					// fmt.Printf("\n Received Trigram score: %v", currentTrigramScore)

					if currentTrigramScore > bestTrigramScore {
						bestTrigramScore = currentTrigramScore
						copyStruct(&bestConfig, &currentConfig)
					}
				}
			}
		}
	}

	// fmt.Printf("\n\n\n Best IOC score: %v", bestTrigramScore)

	fmt.Printf("\n\n [+] Best configuration from analysis:")
	printExpected(bestConfig)

	// plainText := runEnigma(cipherText, bestConfig)
	// fmt.Println(plainText)
}
