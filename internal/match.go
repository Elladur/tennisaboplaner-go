package internal

import (
	"errors"
	"fmt"
)

type Match struct {
	Player1      uint8
	Player2      uint8
	IsPlayer2Set bool
}

func CreateMatch(player1 uint8, player2 uint8) (Match, error) {
	if player1 == player2 {
		return Match{}, errors.New("same player is not allowed")
	}
	if player1 < player2 {
		return Match{Player1: player1, Player2: player2, IsPlayer2Set: true}, nil
	} else {
		return Match{Player1: player2, Player2: player1, IsPlayer2Set: true}, nil
	}
}

func CreatePartialMatch(player1 uint8) Match {
	return Match{Player1: player1, IsPlayer2Set: false}
}

func (m Match) GetPlayers() []uint8 {
	if m.IsPlayer2Set {
		return []uint8{m.Player1, m.Player2}
	} else {
		return []uint8{m.Player1}
	}
}

func (m Match) String(players *[]Player) string {
	if m.IsPlayer2Set {
		return fmt.Sprintf("%s vs %s", (*players)[m.Player1].Name, (*players)[m.Player2].Name)
	} else {
		return fmt.Sprintf("%s vs ...", (*players)[m.Player1].Name)
	}
}

func CanMatchBeAdded(matches []Match, match Match) bool {
	for _, m := range matches {
		if m.Player1 == match.Player1 || (match.IsPlayer2Set && m.Player1 == match.Player2) {
			return false
		}
		if m.IsPlayer2Set && (m.Player2 == match.Player1 || (match.IsPlayer2Set && m.Player2 == match.Player2)) {
			return false
		}
	}
	return true
}

func (m *Match) ReplacePlayer(oldPlayer uint8, newPlayer uint8) error {
	for _, player := range m.GetPlayers() {
		if player == newPlayer {
			return errors.New("new player is already in the match")
		}
	}
	if m.Player1 == oldPlayer {
		m.Player1 = newPlayer
		return nil
	}
	if m.IsPlayer2Set && m.Player2 == oldPlayer {
		m.Player2 = newPlayer
		return nil
	}
	return errors.New("old player is not in the match")
}
