package game

import (
	"bufio"
	"fmt"
	"os"
)

// LoadDictionaryFromFile returns a list of words, normalized, from a file.
func LoadDictionaryFromFile(path string) ([]string, error) {
	dictionaryFile, err := os.Open(path)
	if err != nil {
		return []string{}, fmt.Errorf("Could not read dictionary at path %s", path)
	}
	dictionary := []string{}
	scanner := bufio.NewScanner(dictionaryFile)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Println will add back the final '\n'
		dictionary = append(dictionary, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return []string{}, fmt.Errorf("Error when reading dictionary: %v", err)
	}
	return dictionary, nil
}
