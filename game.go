package main

import (
	"fmt"
	"sort"
)

type Player struct {
	words []*WordDoc
	name  string
}

func (p *Player) score() int {
	score := 0
	for _, word := range p.words {
		score += word.length
	}
	return score
}

func (p *Player) addWord(word *WordDoc) {
	p.words = append(p.words, word)
	p.sortWordsByLengthDescending()
}

func (p *Player) removeWord(word *WordDoc) {
	for i, w := range p.words {
		if w == word {
			p.words = append(p.words[:i], p.words[i+1:]...)
			break
		}
	}
	p.sortWordsByLengthDescending()
}

func (p *Player) sortWordsByLengthDescending() {
	sort.Slice(p.words, func(i, j int) bool {
		return p.words[i].length > p.words[j].length
	})
}

func (p *Player) getWordsString() string {
	wordsString := ""
	for _, word := range p.words {
		wordsString += word.word + " "
	}
	return wordsString
}

type Game struct {
	dictionary        *Dictionary
	lettersToFlip     []byte
	numLettersFlipped int
	players           []*Player
	lettersInPlay     LetterSet
	logTurns          bool
}

func NewGame(dictionary *Dictionary, lettersToFlip []byte, playerNames []string, logTurns bool) *Game {
	players := []*Player{}
	for _, name := range playerNames {
		players = append(players, &Player{name: name})
	}

	return &Game{dictionary: dictionary, lettersToFlip: lettersToFlip, numLettersFlipped: 0, players: players, logTurns: logTurns}
}

func (g *Game) nextPlayer(playerIndex int) int {
	return (playerIndex + 1) % len(g.players)
}

func (g *Game) getOtherPlayers(playerIndex int) []int {
	otherPlayers := []int{}
	for p := g.nextPlayer(playerIndex); p != playerIndex; p = g.nextPlayer(p) {
		otherPlayers = append(otherPlayers, p)
	}
	return otherPlayers
}

var turnIndentationPrefix = "   "

func (g *Game) playTurn(playerIndex int) (actionTaken bool) {
	player := g.players[playerIndex]
	otherPlayers := g.getOtherPlayers(playerIndex)
	for _, otherPlayerIndex := range otherPlayers {
		otherPlayer := g.players[otherPlayerIndex]
		candidateWords := otherPlayer.words
		for _, candidateWord := range candidateWords {
			stealWord := g.dictionary.findStealWord(candidateWord, &g.lettersInPlay)
			if stealWord != nil {
				// todo maybe make this faster
				lettersInPlayRemoved := stealWord.letterSet.diff(&candidateWord.letterSet)
				g.lettersInPlay.removeLetterSet(&lettersInPlayRemoved)

				if g.logTurns {
					fmt.Printf("%sSTEAL! %s steals '%s' from %s using new letters '%s' to create '%s'\n",
						turnIndentationPrefix, player.name, candidateWord.word, otherPlayer.name, lettersInPlayRemoved.String(), stealWord.word)
				}
				player.addWord(stealWord)
				otherPlayer.removeWord(candidateWord)
				return true
			}
		}
	}

	makeWord := g.dictionary.findNewWord(&g.lettersInPlay)
	if makeWord != nil {
		if g.logTurns {
			fmt.Printf("%sNEW WORD! %s makes word '%s'\n",
				turnIndentationPrefix, player.name, makeWord.word)
		}
		player.addWord(makeWord)
		g.lettersInPlay.removeLetterSet(&makeWord.letterSet)
		return true
	}

	return false
}

func (g *Game) play() {
	for turnNumber := 0; turnNumber < len(g.lettersToFlip); turnNumber++ {
		letter := g.lettersToFlip[turnNumber]
		g.lettersInPlay.addLetter(letter)
		if g.logTurns {
			fmt.Printf("turn %d: flipped '%s', letters in play: '%s'\n", turnNumber, string(letter), g.lettersInPlay.String())
		}
		playersWhoTookNoAction := map[int]struct{}{}
		for playerIndex := g.nextPlayer(turnNumber); len(playersWhoTookNoAction) < len(g.players); playerIndex = g.nextPlayer(playerIndex) {
			actionTaken := g.playTurn(playerIndex)
			if !actionTaken {
				playersWhoTookNoAction[playerIndex] = struct{}{}
			} else {
				// if any player took an action, reset the state because a player who didnt take an action
				// possibly could now
				playersWhoTookNoAction = map[int]struct{}{}
			}
		}
		if g.logTurns {
			for _, player := range g.players {
				fmt.Printf("%s%s: %s\n", turnIndentationPrefix, player.name, player.getWordsString())
			}
		}
	}

	for _, player := range g.players {
		fmt.Println(player.name, "score:", player.score(), "words:", player.getWordsString())
	}
}
