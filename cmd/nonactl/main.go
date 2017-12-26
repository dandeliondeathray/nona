package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	switch os.Args[1] {
	case "newround":
		seed, err := parseNewRound(os.Args[2:])
		if err != nil {
			log.Fatalf("Bad flags: %s", err)
		}
		newRound(seed)
		break
	default:
		log.Fatalf("Unknown command '%s'", os.Args[1])
	}
}

func parseNewRound(args []string) (int64, error) {
	if len(args) == 0 {
		return 0, fmt.Errorf("seed required")
	}
	if len(args) > 1 {
		return 0, fmt.Errorf("too many arguments")
	}

	seed, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("'%s' could not be parsed as an int64", args[0])
	}

	return seed, nil
}

func newRound(seed int64) {
	log.Printf("Seed: %d", seed)
}
