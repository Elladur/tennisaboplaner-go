package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundIndizesOfPlayer(t *testing.T) {
	player := 1
	schedule := Schedule{
		{
			{player1: 1, isPlayer2Set: true, player2: 2},
			{player1: 3, isPlayer2Set: true, player2: 4},
		},
		{
			{player1: 5, isPlayer2Set: true, player2: 3},
			{player1: 2, isPlayer2Set: false},
		},
		{
			{player1: 5, isPlayer2Set: true, player2: 3},
			{player1: 1, isPlayer2Set: false},
		},
	}

	expected := []int{0, 2}
	result := schedule.RoundIndizesOfPlayer(player)
	assert.Equal(t, expected, result)

	result = schedule.RoundIndizesOfPlayer(10)
	assert.Empty(t, result)
}

func TestGetMatchesCountOfPlayer(t *testing.T) {
	player := 1
	schedule := Schedule{
		{
			{player1: 1, isPlayer2Set: true, player2: 2},
			{player1: 3, isPlayer2Set: true, player2: 4},
		},
		{
			{player1: 5, isPlayer2Set: true, player2: 3},
			{player1: 2, isPlayer2Set: false},
		},
		{
			{player1: 5, isPlayer2Set: true, player2: 3},
			{player1: 1, isPlayer2Set: false},
		},
	}

	result := schedule.MatchesCountOfPlayer(player)
	assert.Equal(t, 2, result)

	result = schedule.MatchesCountOfPlayer(10)
	assert.Equal(t, 0, result)
}

func TestGetRoundIndizesOfMatch(t *testing.T) {
	match := Match{player1: 1, isPlayer2Set: true, player2: 2}
	schedule := Schedule{
		{
			{player1: 1, isPlayer2Set: true, player2: 2},
			{player1: 3, isPlayer2Set: true, player2: 4},
		},
		{
			{player1: 5, isPlayer2Set: true, player2: 1},
			{player1: 1, isPlayer2Set: true, player2: 3},
		},
		{
			{player1: 5, isPlayer2Set: true, player2: 1},
			{player1: 1, isPlayer2Set: true, player2: 2},
		},
	}

	expected := []int{0, 2}
	result := schedule.RoundIndizesOfMatch(match)
	assert.Equal(t, expected, result)

	result = schedule.RoundIndizesOfMatch(Match{player1: 10})
	assert.Empty(t, result)
}

func TestGetCountMatchInSchedule(t *testing.T) {
	match := Match{player1: 1, isPlayer2Set: true, player2: 2}
	schedule := Schedule{
		{
			{player1: 1, isPlayer2Set: true, player2: 2},
			{player1: 3, isPlayer2Set: true, player2: 4},
		},
		{
			{player1: 5, isPlayer2Set: true, player2: 1},
			{player1: 1, isPlayer2Set: true, player2: 3},
		},
		{
			{player1: 5, isPlayer2Set: true, player2: 1},
			{player1: 1, isPlayer2Set: true, player2: 2},
		},
	}

	result := schedule.CountOfMatch(match)
	assert.Equal(t, 2, result)

	result = schedule.CountOfMatch(Match{player1: 10})
	assert.Empty(t, 0, result)
}
