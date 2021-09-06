// enigma "LOTSOFPEOPLEWOR" --rotors "Gamma VI IV III" --rings="1 1 1 16" --position "D A B Q" --plugboard "MS KU FY AG BN PQ HJ DI ER LW" --reflector "C-thin"
// Use `github.com/becgabri/enigma`

package main

import (
	"fmt"
	"math"
	"os"
)

type Enigma struct {
	Reflector string
	Rings     []int
	Positions []string
	Rotors    []string
}

var defaultConfig = Enigma{
	Reflector: "C-thin",
	Rings:     []int{1, 1, 1, 16},
	Positions: []string{"", "", "B", "Q"},
	Rotors:    []string{"", "", "IV", "III"},
}

var bestConfig = Enigma{
	Reflector: "C-thin",
	Rings:     []int{1, 1, 1, 16},
	Positions: []string{"", "", "B", "Q"},
	Rotors:    []string{"", "", "IV", "III"},
}

const (
	TRIGRAMS_FILENAME = "english_trigrams.txt"
)

var (
	// rotorsPositionOne = []string{"Beta", "Gamma", "I", "II", "V", "VI"}
	// rotorsPositionTwo = []string{"I", "II", "V", "VI"}
	rotorsPositionOne = []string{"Gamma"}
	rotorsPositionTwo = []string{"VI"}
)

func main() {
	fmt.Printf("\n --- Enigma cryptanalysis --- \n")

	if len(os.Args) != 2 {
		fmt.Printf("\n [-] Cipher text file not provided, exiting\n")
		os.Exit(0)
	}

	cipherBytes, err := readFile(os.Args[1])
	if err != nil {
		fmt.Printf("\n [-] Error reading cipher text file: %v", err)
		os.Exit(0)
	}

	cipherText := string(cipherBytes)

	fmt.Printf("\n [+] Read cipher text :\n %v", cipherText)

	bestTrigramScore := math.Inf(-1)

	initializeTrigrams(TRIGRAMS_FILENAME)

	for _, rotorPositionOne := range rotorsPositionOne { // rotor at 1st place
		for _, rotorPositionTwo := range rotorsPositionTwo { // rotor at 2nd place
			if rotorPositionOne == rotorPositionTwo {
				continue
			}

			defaultConfig.Rotors[0] = rotorPositionOne
			defaultConfig.Rotors[1] = rotorPositionTwo

			for i := 0; i < 26; i++ { // starting position of rotor at 1st place
				for j := 0; j < 26; j++ { // starting position of rotor at 2nd place
					defaultConfig.Positions[0] = string(rune(i + 65))
					defaultConfig.Positions[1] = string(rune(j + 65))

					printConfig(defaultConfig)

					currentPlugboard := getBestPlugboard(cipherText)

					plainText := runEnigma(cipherText, currentPlugboard)

					currentTrigramScore := calculateTrigram(plainText)

					fmt.Printf("\n Received Trigram score: %v", currentTrigramScore)
					fmt.Printf("\n\n ----------------------")

					if currentTrigramScore > bestTrigramScore {
						bestTrigramScore = currentTrigramScore
						copyStruct(&bestConfig, &defaultConfig)
					}
				}
			}
		}
	}

	// res := formatPlugboard("UGCDEFBHIJKLMNOPQRSTAVWXYZ")
	// fmt.Printf("\n\n\n Plugboard: %+v", res)

	fmt.Printf("\n\n\n Best IOC score: %v", bestTrigramScore)
	printConfig(bestConfig)
}
