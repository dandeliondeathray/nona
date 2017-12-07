package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dandeliondeathray/nona/service/chain"
	"github.com/dandeliondeathray/nona/service/plumber"
)

func main() {
	schemasPath := os.Getenv("SCHEMA_PATH")
	brokerEnv := os.Getenv("KAFKA_BROKERS")
	if brokerEnv == "" {
		log.Fatalf("No KAFKA_BROKERS set!")
	}
	dictionaryPath := os.Getenv("DICTIONARY_PATH")
	if dictionaryPath == "" {
		log.Fatalf("No DICTIONARY_PATH set!")
	}
	brokers := strings.Split(brokerEnv, ",")
	log.Println("Kafka brokers:", brokers)

	portEnv := os.Getenv("NONA_CHAIN_PORT")
	port := 8080
	if portEnv != "" {
		portInt, err := strconv.Atoi(portEnv)
		if err != nil {
			log.Fatalf("Port is not an int: %v", portEnv)
		}
		port = portInt
	}

	codecs, err := plumber.LoadCodecsFromPath(schemasPath)
	if err != nil {
		log.Fatalf("Could not load codecs from path %s", schemasPath)
	}

	dictionary, err := readLinesFromDictionary(dictionaryPath)
	if err != nil {
		log.Fatalf("Could not read dictionary, because: %s", err)
	}
	teams := chain.NewTeams(dictionary)

	service := chain.NewService(teams)
	service.Start()

	plumber := plumber.NewPlumber(service, codecs)
	if err = plumber.Start(brokers); err != nil {
		log.Fatalf("Could not start plumber: %s", err)
	}

	service.ListenAndServe(port)

	chBlock := make(chan bool)
	<-chBlock
}

func readLinesFromDictionary(path string) ([]string, error) {
	inFile, err := os.Open(path)
	if err != nil {
		return []string{}, err
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	dictionary := make([]string, 0, 1000)
	for scanner.Scan() {
		dictionary = append(dictionary, scanner.Text())
	}
	return dictionary, nil
}
