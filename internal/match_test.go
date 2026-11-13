package internal

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMatch(t *testing.T) {
	tests := []struct {
		player1   int
		player2   int
		expectErr bool
		expected  Match
	}{
		{1, 2, false, Match{1, 2, true}},
		{2, 1, false, Match{1, 2, true}},
		{1, 1, true, Match{}},
	}

	for _, tt := range tests {
		m, err := createMatch(tt.player1, tt.player2)
		if tt.expectErr {
			assert.Error(t, err, "expected an error but got none")
		} else {
			assert.NoError(t, err, "did not expect an error but got one")
			assert.Equal(t, tt.expected, m)
		}
	}
}

func TestCreatePartialMatch(t *testing.T) {
	player1 := 1
	match := createPartialMatch(player1)
	assert.Equal(t, player1, match.player1, "Player1 should be set correctly")
	assert.False(t, match.isPlayer2Set, "IsPlayer2Set should be false")
}

func TestGetPlayers(t *testing.T) {
	match := Match{player1: 1, player2: 2, isPlayer2Set: true}
	players := match.getPlayers()
	assert.Equal(t, []int{1, 2}, players, "Players should be [1, 2]")

	partialMatch := Match{player1: 1, isPlayer2Set: false}
	players = partialMatch.getPlayers()
	assert.Equal(t, []int{1}, players, "Players should be [1]")
}

func TestMatchJSONRoundTrip(t *testing.T) {
	m := createPartialMatch(7)
	data, err := json.Marshal(m)
	assert.NoError(t, err)
	var got Match
	err = json.Unmarshal(data, &got)
	assert.NoError(t, err)
	assert.Equal(t, m, got)
	assert.False(t, got.isPlayer2Set)
}
