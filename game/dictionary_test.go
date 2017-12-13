package game_test

import (
	"testing"

	"github.com/dandeliondeathray/nona/game"
)

func TestDictionary_LoadFromFile_KnownWordsArePresent(t *testing.T) {
	path := "dictionary.txt"
	words, err := game.LoadDictionaryFromFile(path)
	if err != nil {
		t.Fatalf("Could read dictionary at %s because: %v", path, err)
	}

	foundKnownWord := false
	for _, w := range words {
		if w == "PUSSGURKA" {
			foundKnownWord = true
		}
	}
	if !foundKnownWord {
		t.Fatalf("Did not find known word PUSSGURKA in the dictionary")
	}
}
