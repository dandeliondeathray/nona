package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var url string

func main() {

	flag.StringVar(&url, "url", "http://localhost:8080", "URL for the nona control layer.")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		os.Exit(1)
	}

	switch args[0] {
	case "newround":
		seed, err := parseNewRound(args[1:])
		if err != nil {
			log.Fatalf("Bad flags: %s", err)
		}
		newRound(seed)
		break
	default:
		log.Fatalf("Unknown command '%s'", args[0])
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
	url := fmt.Sprintf("%s/newround/%d", url, seed)
	resp, err := http.Post(url, "application/text", nil)
	if err != nil {
		log.Println("New round request failed, because:", err)
		return
	}
	if resp.StatusCode != 200 {
		log.Printf("Unexpected status code %d", resp.StatusCode)
	}
}
