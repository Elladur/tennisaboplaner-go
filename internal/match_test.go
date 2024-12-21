package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMatch(t *testing.T) {
	tests := []struct {
		player1   uint8
		player2   uint8
		expectErr bool
		expected  Match
	}{
		{1, 2, false, Match{1, 2, true}},
		{2, 1, false, Match{1, 2, true}},
		{1, 1, true, Match{}},
	}

	for _, tt := range tests {
		m, err := CreateMatch(tt.player1, tt.player2)
		if tt.expectErr {
			assert.Error(t, err, "expected an error but got none")
		} else {
			assert.NoError(t, err, "did not expect an error but got one")
			assert.Equal(t, tt.expected, m)
		}
	}
}

func TestCreatePartialMatch(t *testing.T) {
	player1 := uint8(1)
	match := CreatePartialMatch(player1)
	assert.Equal(t, player1, match.Player1, "Player1 should be set correctly")
	assert.False(t, match.IsPlayer2Set, "IsPlayer2Set should be false")
}

func TestGetPlayers(t *testing.T) {
	match := Match{Player1: 1, Player2: 2, IsPlayer2Set: true}
	players := match.GetPlayers()
	assert.Equal(t, []uint8{1, 2}, players, "Players should be [1, 2]")

	partialMatch := Match{Player1: 1, IsPlayer2Set: false}
	players = partialMatch.GetPlayers()
	assert.Equal(t, []uint8{1}, players, "Players should be [1]")
}

func TestCanMatchBeAdded(t *testing.T) {
	matches := []Match{
		{Player1: 1, Player2: 2, IsPlayer2Set: true},
		{Player1: 3, Player2: 4, IsPlayer2Set: true},
		{Player1: 7, IsPlayer2Set: false},
	}

	tests := []struct {
		match      Match
		canBeAdded bool
	}{
		{Match{Player1: 5, Player2: 6, IsPlayer2Set: true}, true},
		{Match{Player1: 1, Player2: 3, IsPlayer2Set: true}, false},
		{Match{Player1: 2, Player2: 4, IsPlayer2Set: true}, false},
		{Match{Player1: 5, IsPlayer2Set: false}, true},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.canBeAdded, CanMatchBeAdded(matches, tt.match), "CanMatchBeAdded(%v) should be %v", tt.match, tt.canBeAdded)
	}
}

func TestReplacePlayer(t *testing.T) {
	match := Match{Player1: 1, Player2: 2, IsPlayer2Set: true}

	err := match.ReplacePlayer(1, 3)
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, uint8(3), match.Player1, "Player1 should be 3")

	err = match.ReplacePlayer(2, 4)
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, uint8(4), match.Player2, "Player2 should be 4")

	err = match.ReplacePlayer(1, 4)
	assert.Equal(t, Match{3, 4, true}, match, "expected match to be {3, 4, true}")
	assert.Error(t, err, "expected an error")

	err = match.ReplacePlayer(5, 6)
	assert.Equal(t, Match{3, 4, true}, match, "expected match to be {3, 4, true}")
	assert.Error(t, err, "expected an error")
}
