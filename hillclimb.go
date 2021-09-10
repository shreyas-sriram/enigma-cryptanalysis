/*
 *
 * Program uses `github.com/emedvedev/enigma` as an external API to simulate enigma configurations
 *
 */

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

// currentConfig stores the enigma configurations through the
// hillclimb attack iterations
var currentConfig = enigma{
	reflector: "C-thin",
	rings:     []int{1, 1, 1, 16},
	positions: []string{"_", "_", "B", "Q"},
	rotors:    []string{"_", "_", "IV", "III"},
	plugboard: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
}

// currentConfig stores the best enigma configuration
var bestConfig = enigma{
	reflector: "C-thin",
	rings:     []int{1, 1, 1, 16},
	positions: []string{"_", "_", "B", "Q"},
	rotors:    []string{"_", "_", "IV", "III"},
	plugboard: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
}

var (
	rotorsPositionOne = [...]string{"I", "II", "V", "VI", "Beta", "Gamma"}
	rotorsPositionTwo = [...]string{"I", "II", "V", "VI"}
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("\n [-] Cipher text file not provided, exiting\n")
		os.Exit(0)
	}

	cipherText := readFile(os.Args[1])

	// fmt.Printf("\n [+] Read cipher text :\n %v", cipherText)

	bestTrigramScore := math.Inf(-1)
	totalIOC := float64(0)
	totalIterations := float64(0)

	initializeTrigrams()

	for _, rotorPositionOne := range rotorsPositionOne { // rotor at 1st place
		for _, rotorPositionTwo := range rotorsPositionTwo { // rotor at 2nd place

			if rotorPositionOne == rotorPositionTwo {
				continue
			}

			currentConfig.rotors[0] = rotorPositionOne
			currentConfig.rotors[1] = rotorPositionTwo

			// fmt.Printf("\n\n [+] Trying enigma configuration:")
			// printConfig(currentConfig)

			for i := 0; i < 26; i++ { // starting position of rotor at 1st place
				for j := 0; j < 26; j++ { // starting position of rotor at 2nd place
					currentConfig.positions[0] = string(rune(i + 65))
					currentConfig.positions[1] = string(rune(j + 65))

					// Optimization - brings down runtime from 30 minutes to 1 minute or less
					plainText := runEnigmaWithPlugboard(cipherText, defaultPlugboard)
					currentIOC := calculateIOC(plainText)

					if currentIOC <= totalIOC/totalIterations {
						continue
					} else {
						totalIOC += currentIOC
						totalIterations++
					}

					// Hill-climb attack
					currentPlugboard := getBestPlugboard(cipherText, currentIOC)
					currentConfig.plugboard = currentPlugboard

					plainText = runEnigma(cipherText, currentConfig)
					currentTrigramScore := calculateTrigram(plainText)

					// fmt.Println(currentIOC)
					// fmt.Println(currentPlugboard)
					// fmt.Println(plainText)
					// fmt.Println(currentTrigramScore)
					// os.Exit(0)

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

	// fmt.Printf("\n\n [+] Best configuration from analysis:")
	printExpected(bestConfig)

	// plainText := runEnigma(cipherText, bestConfig)
	// fmt.Println(plainText)
}
