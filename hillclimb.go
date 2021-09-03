// enigma "LOTSOFPEOPLEWOR" --rotors "Gamma VI IV III" --rings="1 1 1 16" --position "D A B Q" --plugboard "MS KU FY AG BN PQ HJ DI ER LW" --reflector "C-thin"
// Use `github.com/becgabri/enigma`

package main

import (
	"fmt"
	"os"
)

var defaultConfig = struct {
	Reflector string
	Rings     []int
	Positions []string
	Rotors    []string
}{
	Reflector: "C-thin",
	Rings:     []int{1, 1, 1, 16},
	Positions: []string{"", "", "B", "Q"},
	Rotors:    []string{"", "", "IV", "III"},
}

var bestConfig = struct {
	Reflector string
	Rings     []int
	Positions []string
	Rotors    []string
}{
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

	fmt.Printf("\n [+] Read cipher text :\n %v", string(cipherBytes))

	// config := make([]enigma.RotorConfig, len(defaultConfig.Rotors))
	// for index, rotor := range defaultConfig.Rotors {
	// 	ring := defaultConfig.Rings[index]
	// 	start := defaultConfig.Positions[index][0]
	// 	config[index] = enigma.RotorConfig{ID: rotor, Start: start, Ring: ring}
	// }

	// plugboards := strings.Split("MS KU FY AG BN PQ HJ DI ER LW", " ")
	// e := enigma.NewEnigma(config, defaultConfig.Reflector, plugboards)
	// encoded := e.EncodeString(string(cipherBytes))

	for _, rotorPositionOne := range rotorsPositionOne {
		for _, rotorPositionTwo := range rotorsPositionTwo {
			defaultConfig.Positions[0] = rotorPositionOne
			defaultConfig.Positions[1] = rotorPositionTwo
		}
	}

}
