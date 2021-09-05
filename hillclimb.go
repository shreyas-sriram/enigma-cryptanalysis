// enigma "LOTSOFPEOPLEWOR" --rotors "Gamma VI IV III" --rings="1 1 1 16" --position "D A B Q" --plugboard "MS KU FY AG BN PQ HJ DI ER LW" --reflector "C-thin"
// Use `github.com/becgabri/enigma`

package main

import (
	"fmt"
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

var rotorsPositionOne = []string{"Beta", "Gamma", "I", "II", "V", "VI"}
var rotorsPositionTwo = []string{"I", "II", "V", "VI"}

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

	best := float32(0)

	for _, rotorPositionOne := range rotorsPositionOne {
		for _, rotorPositionTwo := range rotorsPositionTwo {
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

					ioc := getScore(cipherText)

					fmt.Printf("\n Received IOC score: %v", ioc)
					fmt.Printf("\n\n ----------------------")

					if ioc > best {
						best = ioc
						copyStruct(&bestConfig, &defaultConfig)
					}
				}
			}
		}
	}

	fmt.Printf("\n\n\n Best IOC score: %v", best)
	printConfig(bestConfig)
}
