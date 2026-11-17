package internal

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Match represents a match between two players
// player1 is always the player with the lower ID
type Match struct {
	player1      int
	player2      int
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
		val, err := createMatch(match.Player1, match.Player2)
		if err != nil {
			return err
		}
		*m = val
		return nil
	}
	val := createPartialMatch(match.Player1)
	*m = val
	return nil
}

func createMatch(player1 int, player2 int) (Match, error) {
	if player1 == player2 {
		return Match{}, errors.New("same player is not allowed")
	}
	if player1 < player2 {
		return Match{player1: player1, player2: player2, isPlayer2Set: true}, nil
	}
	return Match{player1: player2, player2: player1, isPlayer2Set: true}, nil
}

func createPartialMatch(player1 int) Match {
	return Match{player1: player1, isPlayer2Set: false}
}

func (m Match) getPlayers() []int {
	if m.isPlayer2Set {
		return []int{m.player1, m.player2}
	}
	return []int{m.player1}
}

func (m Match) String(players []Player) string {
	if m.isPlayer2Set {
		return fmt.Sprintf("%s vs %s", players[m.player1].Name, players[m.player2].Name)
	}
	return fmt.Sprintf("%s vs ...", players[m.player1].Name)
}
