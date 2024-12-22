package internal

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Match represents a match between two players
// player1 is always the player with the lower ID
type Match struct {
	player1      uint8
	player2      uint8
	isPlayer2Set bool
}

// MarshalJSON marshals a Match to JSON
func (m Match) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Player1":      int(m.player1),
		"Player2":      int(m.player2),
		"IsPlayer2Set": m.isPlayer2Set,
	})
}

// UnmarshalJSON unmarshals a Match from JSON
func (m *Match) UnmarshalJSON(data []byte) error {
	var match struct {
		Player1      int
		Player2      int
		IsPlayer2Set bool
	}
	err := json.Unmarshal(data, &match)
	if err != nil {
		return err
	}
	if match.IsPlayer2Set {
		val, err := createMatch(uint8(match.Player1), uint8(match.Player2))
		if err != nil {
			return err
		}
		m.player1 = val.player1
		m.player2 = val.player2
		m.isPlayer2Set = true
		return nil
	}
	val := createPartialMatch(uint8(match.Player1))
	m.player1 = val.player1
	m.player2 = val.player2
	m.isPlayer2Set = true
	m = &val
	return nil
}

func createMatch(player1 uint8, player2 uint8) (Match, error) {
	if player1 == player2 {
		return Match{}, errors.New("same player is not allowed")
	}
	if player1 < player2 {
		return Match{player1: player1, player2: player2, isPlayer2Set: true}, nil
	}
	return Match{player1: player2, player2: player1, isPlayer2Set: true}, nil
}

func createPartialMatch(player1 uint8) Match {
	return Match{player1: player1, isPlayer2Set: false}
}

func (m Match) getPlayers() []uint8 {
	if m.isPlayer2Set {
		return []uint8{m.player1, m.player2}
	}
	return []uint8{m.player1}
}

func (m Match) String(players *[]Player) string {
	if m.isPlayer2Set {
		return fmt.Sprintf("%s vs %s", (*players)[m.player1].Name, (*players)[m.player2].Name)
	}
	return fmt.Sprintf("%s vs ...", (*players)[m.player1].Name)
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
	if isInSlice(newPlayer, m.getPlayers()) {
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
