package puzzlestore

import "math/rand"

// Puzzles keeps all puzzles and generates new on demand.
type Puzzles struct {
	random *rand.Rand
	chain  []string
}

// NewPuzzles creates a new Puzzles.
func NewPuzzles(seed int64) *Puzzles {
	return &Puzzles{rand.New(rand.NewSource(seed)), make([]string, 0)}
}

// Get returns a puzzle for a given index.
func (p *Puzzles) Get(index int) string {
	if index < len(p.chain) {
		return p.chain[index]
	}
	noOfPuzzlesInChain := len(p.chain)
	noOfPuzzlesToGenerate := index + 1 - noOfPuzzlesInChain
	for i := 0; i < noOfPuzzlesToGenerate; i++ {
		puzzle := p.shuffle("PUSSGURKA")
		p.chain = append(p.chain, puzzle)
	}
	return p.chain[index]
}

//
// Shuffling the characters of a word
//
func (p *Puzzles) shuffle(word string) string {
	runes := []rune(word)
	perm := p.random.Perm(len(runes))
	shuffledRunes := make([]rune, len(runes))
	for i := range perm {
		p := perm[i]
		shuffledRunes[i] = runes[p]
	}
	return string(shuffledRunes)
}
