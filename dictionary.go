package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"sort"
	"strings"
)

type WordDoc struct {
	length    int
	word      string
	letterSet LetterSet
}

type IndexRange struct {
	start int // inclusive
	end   int // inclusive
}

type Dictionary struct {
	minWordLength            int
	letterSets               []LetterSet
	anagramsByLetterSetIndex []IndexRange
	wordsSorted              []WordDoc
	letterSetsByWordLength   map[int]IndexRange
}

func (d *Dictionary) findStealWord(existingWord *WordDoc, lettersInPlay *LetterSet) *WordDoc {
	existingWordLetterSet := &existingWord.letterSet

	var allLetters LetterSet
	allLetters.addLetterSet(lettersInPlay)
	allLetters.addLetterSet(existingWordLetterSet)

	// no reason to examine any word that is longer than the total number of letters we have
	start := d.letterSetsByWordLength[allLetters.numLetters()].start

	// no reason to examine any word that is shorter than existingWord
	end := d.letterSetsByWordLength[existingWord.length].end

	for idx := start; idx <= end; idx++ {
		ls := d.letterSets[idx]
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

func (d *Dictionary) findNewWord(lettersInPlay *LetterSet) *WordDoc {
	numLettersInPlay := lettersInPlay.numLetters()
	if numLettersInPlay < d.minWordLength {
		return nil
	}
	// no reason to examine any word that is longer than the total number of letters we have
	start := d.letterSetsByWordLength[numLettersInPlay].start

	for idx := start; idx < len(d.letterSets); idx++ {
		ls := d.letterSets[idx]
		if ls.IsSubsetOf(lettersInPlay) {
			wordIndexes := d.anagramsByLetterSetIndex[idx]
			return &d.wordsSorted[wordIndexes.start]
		}
	}
	return nil
}

func newDictionaryFromFile(path string, minWordLength int) (*Dictionary, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return NewDictionary(file, minWordLength), nil
}

func NewDictionary(r io.Reader, minWordLength int) *Dictionary {
	d := &Dictionary{minWordLength: minWordLength}

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
	d.anagramsByLetterSetIndex = make([]IndexRange, 0, len(words))
	d.letterSets = make([]LetterSet, 0, len(words))

	numLetterSets := 0
	for i := 0; i < len(words); i++ {
		word := words[i]
		if i == 0 || !word.letterSet.IsEqual(&d.letterSets[numLetterSets-1]) {
			numLetterSets++
			d.letterSets = append(d.letterSets, word.letterSet)
			indexes := IndexRange{start: i, end: i}
			d.anagramsByLetterSetIndex = append(d.anagramsByLetterSetIndex, indexes)
		} else {
			d.anagramsByLetterSetIndex[numLetterSets-1].end++
		}
	}

	d.letterSetsByWordLength = make(map[int]IndexRange)
	for i := 0; i < len(d.letterSets); i++ {
		numLetters := d.letterSets[i].numLetters()
		rng := d.letterSetsByWordLength[numLetters]
		if rng.start == 0 {
			rng.start = i
		}
		rng.end = i
		d.letterSetsByWordLength[numLetters] = rng
	}

	return d
}
