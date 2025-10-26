package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Exclusion defines a rule where a Giver cannot give a gift to a Receiver.
type Exclusion struct {
	Giver    string `json:"giver"`
	Receiver string `json:"receiver"`
}

// Config holds the full structure of the JSON config file.
type Config struct {
	Participants []string    `json:"participants"`
	Exclusions   []Exclusion `json:"exclusions"`
}

// Pairing holds all {giver: receiver} pairings for a single valid solution.
type Pairing map[string]string

// loadConfigFromFile opens, reads, and parses the JSON config file.
// It now returns a single Config struct containing both participants and exclusions.
func loadConfigFromFile(filename string) (Config, error) {
	var config Config // Initialize an empty config struct

	// 1. Read the file's contents
	byteValue, err := os.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("could not read file %s: %w", filename, err)
	}

	// 2. Unmarshal the JSON data into our Config struct
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return config, fmt.Errorf("could not parse JSON in %s: %w", filename, err)
	}

	return config, nil
}

// RunSecretSanta orchestrates the process.
// (No changes to this function)
func RunSecretSanta(participants []string, exclusions []Exclusion) {
	// 1. SETUP
	exclusionMap := createExclusionMap(participants, exclusions)
	allValidPairings := make([]Pairing, 0)
	currentPairing := make(Pairing)
	availableReceivers := make(map[string]bool)
	for _, p := range participants {
		availableReceivers[p] = true
	}

	// 2. SOLVE
	findPairingsRecursive(participants, 0, currentPairing, availableReceivers, exclusionMap, &allValidPairings)

	// 3. PRINT COUNT
	count := len(allValidPairings)
	fmt.Printf("Found %d possible unique pairings.\n", count)

	// 4. SELECT AND PRINT ONE PAIRING
	if count > 0 {
		rand.Seed(time.Now().UnixNano())
		randomIndex := rand.Intn(count)
		selectedPairing := allValidPairings[randomIndex]

		fmt.Println("\n--- Selected Pairing ---")
		for _, giver := range participants {
			receiver := selectedPairing[giver]
			fmt.Printf("%s 🎁 --> %s\n", giver, receiver)
		}
	} else {
		fmt.Println("No valid pairings could be found with these rules!")
	}
}

// findPairingsRecursive is the core backtracking algorithm.
// (No changes to this function)
func findPairingsRecursive(
	allGivers []string,
	giverIndex int,
	currentPairing Pairing,
	availableReceivers map[string]bool,
	exclusionMap map[string]map[string]bool,
	allValidPairings *[]Pairing,
) {
	if giverIndex == len(allGivers) {
		solutionCopy := make(Pairing)
		for k, v := range currentPairing {
			solutionCopy[k] = v
		}
		*allValidPairings = append(*allValidPairings, solutionCopy)
		return
	}

	currentGiver := allGivers[giverIndex]
	forbiddenReceivers := exclusionMap[currentGiver]

	for potentialReceiver, isAvailable := range availableReceivers {
		if !isAvailable {
			continue
		}

		isSelf := (currentGiver == potentialReceiver)
		isExcluded := forbiddenReceivers[potentialReceiver]

		if !isSelf && !isExcluded {
			currentPairing[currentGiver] = potentialReceiver
			availableReceivers[potentialReceiver] = false

			findPairingsRecursive(allGivers, giverIndex+1, currentPairing, availableReceivers, exclusionMap, allValidPairings)

			availableReceivers[potentialReceiver] = true
			delete(currentPairing, currentGiver)
		}
	}
}

// createExclusionMap builds a map for fast O(1) lookups.
// (No changes to this function)
func createExclusionMap(participants []string, exclusions []Exclusion) map[string]map[string]bool {
	exMap := make(map[string]map[string]bool)
	for _, p := range participants {
		exMap[p] = make(map[string]bool)
	}
	for _, ex := range exclusions {
		if _, ok := exMap[ex.Giver]; ok {
			exMap[ex.Giver][ex.Receiver] = true
		}
	}
	return exMap
}

// --- Main function to run the example ---
func main() {
	// --- 1. Get filename from command-line arguments ---
	if len(os.Args) < 2 {
		fmt.Println("Error: Please provide the config JSON file as an argument.")
		fmt.Println("Usage: go run . config.json")
		os.Exit(1)
	}
	filename := os.Args[1]

	// --- 2. Load Config (participants AND exclusions) ---
	config, err := loadConfigFromFile(filename)
	if err != nil {
		fmt.Printf("Error loading config file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded %d participants and %d exclusion rules from %s\n\n",
		len(config.Participants), len(config.Exclusions), filename)

	// --- 3. Run the generator! ---
	// Pass the loaded data directly to the function.
	RunSecretSanta(config.Participants, config.Exclusions)
}
