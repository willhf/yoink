package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
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

func (d *Dictionary) findStealWord(existingWord *WordDoc, lettersInPlay *LetterSet) *WordDoc {
	existingWordLetterSet := &existingWord.letterSet

	var allLetters LetterSet
	allLetters.addLetterSet(lettersInPlay)
	allLetters.addLetterSet(existingWordLetterSet)

	for idx, ls := range d.letterSets {
		if existingWordLetterSet.IsSubsetOf(&ls) && ls.IsSubsetOf(&allLetters) {
			wordIndexes := d.anagramsByLetterSetIndex[idx]
			for i := wordIndexes.start; i <= wordIndexes.end; i++ {
				candidate := d.wordsSorted[i]
				// NOTE! this strings.Contains is not precisely correct.
				// The rules say that the nature of the word must change:
				// 'engine' can't become 'engines', which this substring check prevents,
				// but this substring check prevents some valid transformations like
				// stealing 'quit' to become 'equity'
				if !strings.Contains(candidate.word, existingWord.word) && candidate.length > existingWord.length {
					return &d.wordsSorted[i]
				}
			}
		}
	}

	return nil
}

func (d *Dictionary) findMakeWord(lettersInPlay *LetterSet) *WordDoc {
	for idx, ls := range d.letterSets {
		if ls.IsSubsetOf(lettersInPlay) {
			wordIndexes := d.anagramsByLetterSetIndex[idx]
			return &d.wordsSorted[wordIndexes.start]
		}
	}
	return nil
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

func (d *Dictionary) print() {
	for i := 0; i < len(d.anagramsByLetterSetIndex); i++ {
		indexes := d.anagramsByLetterSetIndex[i]
		start := indexes.start
		end := indexes.end
		words := []string{}
		for j := start; j <= end; j++ {
			words = append(words, d.wordsSorted[j].word)
		}
		fmt.Println(strings.Join(words, " "))

	}
}
