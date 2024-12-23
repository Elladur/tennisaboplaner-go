package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundIndizesOfPlayer(t *testing.T) {
	player := 1
	schedule := [][]Match{
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
	result := getRoundIndizesOfPlayer(schedule, player)
	assert.Equal(t, expected, result)

	result = getRoundIndizesOfPlayer(schedule, 10)
	assert.Empty(t, result)
}

func TestGetMatchesCountOfPlayer(t *testing.T) {
	player := 1
	schedule := [][]Match{
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

	result := getMatchesCountOfPlayer(schedule, player)
	assert.Equal(t, 2, result)

	result = getMatchesCountOfPlayer(schedule, 10)
	assert.Equal(t, 0, result)
}

func TestGetRoundIndizesOfMatch(t *testing.T) {
	match := Match{player1: 1, isPlayer2Set: true, player2: 2}
	schedule := [][]Match{
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
	result := getRoundIndizesOfMatch(schedule, match)
	assert.Equal(t, expected, result)

	result = getRoundIndizesOfMatch(schedule, Match{player1: 10})
	assert.Empty(t, result)
}

func TestGetCountMatchInSchedule(t *testing.T) {
	match := Match{player1: 1, isPlayer2Set: true, player2: 2}
	schedule := [][]Match{
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

	result := getCountOfMatchInSchedule(schedule, match)
	assert.Equal(t, 2, result)

	result = getCountOfMatchInSchedule(schedule, Match{player1: 10})
	assert.Empty(t, 0, result)
}
