package main

import (
	"bufio"
	"bytes"
	"io"
	"sort"
)

// WordKey is the index into Dictionary.wordsSorted
type WordKey int

type WordDoc struct {
	length    int
	word      string
	letterSet LetterSet
}

type WordAnagramIndexes struct {
	start int // inclusive
	end   int // inclusive
}

type Dictionary struct {
	letterSets               []LetterSet
	anagramsByLetterSetIndex []WordAnagramIndexes
	wordsSorted              []WordDoc
}

func NewDictionary(r io.Reader, minWordLength int) *Dictionary {
	d := &Dictionary{}

	words := []WordDoc{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) >= minWordLength {
			word := WordDoc{length: len(line), word: line}
			populateLetterSetFromString(&word.letterSet, line)
			words = append(words, word)
		}
	}

	sort.Slice(words, func(a, b int) bool {
		lenA := words[a].length
		lenB := words[b].length
		if lenA != lenB {
			return lenB < lenA // longest first
		}
		lettersA := words[a].letterSet.letters[:]
		lettersB := words[b].letterSet.letters[:]
		return bytes.Compare(lettersA, lettersB) < 0
	})

	d.wordsSorted = words
	d.anagramsByLetterSetIndex = make([]WordAnagramIndexes, 0, len(words))
	d.letterSets = make([]LetterSet, 0, len(words))

	numLetterSets := 0
	for i := 0; i < len(words); i++ {
		word := words[i]
		if i == 0 || !word.letterSet.IsEqual(&d.letterSets[numLetterSets-1]) {
			numLetterSets++
			d.letterSets = append(d.letterSets, word.letterSet)
			indexes := WordAnagramIndexes{start: i, end: i}
			d.anagramsByLetterSetIndex = append(d.anagramsByLetterSetIndex, indexes)
		} else {
			d.anagramsByLetterSetIndex[numLetterSets-1].end++
		}
	}

	return d
}

// func (d *Dictionary) print() {
// 	for i := 0; i < len(dict.anagramsByLetterSetIndex); i++ {
// 		indexes := dict.anagramsByLetterSetIndex[i]
// 		start := indexes.start
// 		end := indexes.end
// 		words := []string{}
// 		for j := start; j <= end; j++ {
// 			words = append(words, dict.wordsSorted[j].word)
// 		}
// 		fmt.Println(strings.Join(words, " "))

// 	}
// }
