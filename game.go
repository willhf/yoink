package main

type Player struct {
	wordsByLengthDescending []WordKey
	name                    string
}

type Game struct {
	dictionary        *Dictionary
	lettersToFlip     []byte
	numLettersFlipped int
	players           []*Player
}

func NewGame(dictionary *Dictionary, lettersToFlip []byte, playerNames []string) *Game {
	players := []*Player{}
	for _, name := range playerNames {
		players = append(players, &Player{name: name})
	}

	return &Game{dictionary: dictionary, lettersToFlip: lettersToFlip, numLettersFlipped: 0, players: players}
}

func (g *Game) play() {

}
