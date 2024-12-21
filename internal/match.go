package internal

import (
	"errors"
	"fmt"
)

type Match struct {
	player1      uint8
	player2      uint8
	isPlayer2Set bool
}

func createMatch(player1 uint8, player2 uint8) (Match, error) {
	if player1 == player2 {
		return Match{}, errors.New("same player is not allowed")
	}
	if player1 < player2 {
		return Match{player1: player1, player2: player2, isPlayer2Set: true}, nil
	} else {
		return Match{player1: player2, player2: player1, isPlayer2Set: true}, nil
	}
}

func createPartialMatch(player1 uint8) Match {
	return Match{player1: player1, isPlayer2Set: false}
}

func (m Match) GetPlayers() []uint8 {
	if m.isPlayer2Set {
		return []uint8{m.player1, m.player2}
	} else {
		return []uint8{m.player1}
	}
}

func (m Match) String(players *[]Player) string {
	if m.isPlayer2Set {
		return fmt.Sprintf("%s vs %s", (*players)[m.player1].Name, (*players)[m.player2].Name)
	} else {
		return fmt.Sprintf("%s vs ...", (*players)[m.player1].Name)
	}
}

func canMatchBeAdded(matches []Match, match Match) bool {
	for _, m := range matches {
		if m.player1 == match.player1 || (match.isPlayer2Set && m.player1 == match.player2) {
			return false
		}
		if m.isPlayer2Set && (m.player2 == match.player1 || (match.isPlayer2Set && m.player2 == match.player2)) {
			return false
		}
	}
	return true
}

func (m *Match) replacePlayer(oldPlayer uint8, newPlayer uint8) error {
	if isInSlice(newPlayer, m.GetPlayers()) {
		return errors.New("new player is already in the match")
	}
	if m.player1 == oldPlayer {
		m.player1 = newPlayer
		return nil
	}
	if m.isPlayer2Set && m.player2 == oldPlayer {
		m.player2 = newPlayer
		return nil
	}
	return errors.New("old player is not in the match")
}
