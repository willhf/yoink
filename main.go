package main

import (
	"flag"
	"fmt"
	"os"
)

// TODO: make the input letter distribution a configurable input
var inputLetterDistribution = LetterSet{
	letters: [NUM_LETTERS]byte{12, 3, 5, 6, 18, 3, 6, 4, 11, 2, 2, 7, 4, 9, 10, 4, 2, 10, 8, 9, 5, 2, 2, 2, 2, 2},
}

func main() {
	dictionaryPath := flag.String("dictionary", "", "the location of the dictionary file")
	logTurns := flag.Bool("log-turns", true, "log each turn")
	seed := flag.Int("seed", 1234567891, "the seed for the random number generator")
	minWordLength := flag.Int("min-word-length", 4, "the minimum word length")
	flag.Parse()

	dictionaryFile, err := os.Open(*dictionaryPath)
	if err != nil {
		fmt.Println("Error opening dictionary file:", err)
		os.Exit(1)
	}
	defer dictionaryFile.Close()

	dict := NewDictionary(dictionaryFile, *minWordLength)
	if dict == nil {
		fmt.Println("nil dict")
	}

	flipOrder := inputLetterDistribution.toFlipOrder(*seed)
	game := NewGame(dict, flipOrder, []string{"Alice", "Bob", "Charlie"}, *logTurns)
	game.play()
}
