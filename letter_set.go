package main

import (
	"bytes"
)

const NUM_LETTERS = 26

type LetterSet struct {
	letters [NUM_LETTERS]byte
}

func NewLetterSetFromWord(s string) *LetterSet {
	ls := &LetterSet{}
	populateLetterSetFromString(ls, s)
	return ls
}

func populateLetterSetFromString(letterSet *LetterSet, s string) *LetterSet {
	for _, letter := range s {
		letterSet.letters[letter-'a']++
	}
	return letterSet
}

func (ls *LetterSet) IsSubsetOf(other *LetterSet) bool {
	for i := 0; i < NUM_LETTERS; i++ {
		if ls.letters[i] > other.letters[i] {
			return false
		}
	}
	return true
}

func (ls *LetterSet) String() string {
	buf := bytes.Buffer{}
	for i := 0; i < NUM_LETTERS; i++ {
		for j := 0; j < int(ls.letters[i]); j++ {
			buf.WriteString(string(rune(i + 'a')))
		}
	}
	return buf.String()
}

func (ls *LetterSet) addLetter(letter byte) {
	ls.letters[letter-'a']++
}

func (ls *LetterSet) addLetterSet(other *LetterSet) {
	for i := 0; i < NUM_LETTERS; i++ {
		ls.letters[i] += other.letters[i]
	}
}

func (ls *LetterSet) diff(other *LetterSet) LetterSet {
	diff := LetterSet{}
	for i := 0; i < NUM_LETTERS; i++ {
		diff.letters[i] = ls.letters[i] - other.letters[i]
	}
	return diff
}

func (ls *LetterSet) removeLetterSet(other *LetterSet) {
	for i := 0; i < NUM_LETTERS; i++ {
		ls.letters[i] -= other.letters[i]
	}
}

func (ls *LetterSet) IsEqual(other *LetterSet) bool {
	return bytes.Equal(ls.letters[:], other.letters[:])
}

func (ls *LetterSet) toFlipOrder(seed int) []byte {
	charsToFlip := []byte{}

	for letterIndex, numInstancesOfLetter := range ls.letters {
		letter := 'a' + byte(letterIndex)
		for i := 0; i < int(numInstancesOfLetter); i++ {
			charsToFlip = append(charsToFlip, letter)
		}
	}

	for i := 0; i < len(charsToFlip); i++ {
		j := pseudoNoise(i, seed) % len(charsToFlip)
		charsToFlip[i], charsToFlip[j] = charsToFlip[j], charsToFlip[i]
	}

	return charsToFlip
}

// Simple deterministic hash function
func pseudoNoise(i, seed int) int {
	x := uint32(i*7349 + seed*3797)
	x ^= x >> 13
	x *= 0x5bd1e995
	x ^= x >> 15
	return int(x)
}
