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
	player1 := uint8(1)
	match := createPartialMatch(player1)
	assert.Equal(t, player1, match.player1, "Player1 should be set correctly")
	assert.False(t, match.isPlayer2Set, "IsPlayer2Set should be false")
}

func TestGetPlayers(t *testing.T) {
	match := Match{player1: 1, player2: 2, isPlayer2Set: true}
	players := match.getPlayers()
	assert.Equal(t, []uint8{1, 2}, players, "Players should be [1, 2]")

	partialMatch := Match{player1: 1, isPlayer2Set: false}
	players = partialMatch.getPlayers()
	assert.Equal(t, []uint8{1}, players, "Players should be [1]")
}

func TestCanMatchBeAdded(t *testing.T) {
	matches := []Match{
		{player1: 1, player2: 2, isPlayer2Set: true},
		{player1: 3, player2: 4, isPlayer2Set: true},
		{player1: 7, isPlayer2Set: false},
	}

	tests := []struct {
		match      Match
		canBeAdded bool
	}{
		{Match{player1: 5, player2: 6, isPlayer2Set: true}, true},
		{Match{player1: 1, player2: 3, isPlayer2Set: true}, false},
		{Match{player1: 2, player2: 4, isPlayer2Set: true}, false},
		{Match{player1: 5, isPlayer2Set: false}, true},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.canBeAdded, canMatchBeAdded(matches, tt.match), "CanMatchBeAdded(%v) should be %v", tt.match, tt.canBeAdded)
	}
}

func TestReplacePlayer(t *testing.T) {
	match := Match{player1: 1, player2: 2, isPlayer2Set: true}

	err := match.replacePlayer(1, 3)
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, uint8(3), match.player1, "Player1 should be 3")

	err = match.replacePlayer(2, 4)
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, uint8(4), match.player2, "Player2 should be 4")

	err = match.replacePlayer(1, 4)
	assert.Equal(t, Match{3, 4, true}, match, "expected match to be {3, 4, true}")
	assert.Error(t, err, "expected an error")

	err = match.replacePlayer(5, 6)
	assert.Equal(t, Match{3, 4, true}, match, "expected match to be {3, 4, true}")
	assert.Error(t, err, "expected an error")
}
