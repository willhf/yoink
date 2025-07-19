package main

import (
	"flag"
	"fmt"
	"os"
)

var inputLetterDistribution = LetterSet{
	letters: [NUM_LETTERS]byte{12, 3, 5, 6, 18, 3, 6, 4, 11, 2, 2, 7, 4, 9, 10, 4, 2, 10, 8, 9, 5, 2, 2, 2, 2, 2},
}

func main() {
	dictionaryPath := flag.String("dictionary", "", "the location of the dictionary file")
	flag.Parse()

	dictionaryFile, err := os.Open(*dictionaryPath)
	if err != nil {
		fmt.Println("Error opening dictionary file:", err)
		os.Exit(1)
	}
	defer dictionaryFile.Close()

	minWordLength := 4
	dict := NewDictionary(dictionaryFile, minWordLength)
	if dict == nil {
		fmt.Println("nil dict")
	}

	seed := 1234567890
	flipOrder := inputLetterDistribution.toFlipOrder(seed)
	game := NewGame(dict, flipOrder, []string{"Player 1", "Player 2", "Player 3"})
	game.play()
}
