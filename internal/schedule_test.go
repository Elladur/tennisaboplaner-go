package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMatchIndizesOfPlayer(t *testing.T) {
	player := uint8(1)
	schedule := [][]Match{
		{
			{player1: 1, isPlayer2Set: true, player2: 2},
			{player1: 3, isPlayer2Set: true, player2: 4},
		},
		{
			{player1: 5, isPlayer2Set: true, player2: 1},
			{player1: 2, isPlayer2Set: false},
		},
		{
			{player1: 5, isPlayer2Set: true, player2: 3},
			{player1: 1, isPlayer2Set: false},
		},
	}

	expected := [][]int{{0, 0}, {1, 0}, {2, 1}}
	result := getMatchIndizesOfPlayer(schedule, player)
	assert.Equal(t, expected, result)
}

func TestGetMatchIndizesOfMatch(t *testing.T) {
	match := Match{player1: 1, isPlayer2Set: true, player2: 2}
	schedule := [][]Match{
		{
			{player1: 1, isPlayer2Set: true, player2: 2},
			{player1: 3, isPlayer2Set: true, player2: 4},
		},
		{
			{player1: 5, isPlayer2Set: true, player2: 1},
			{player1: 1, isPlayer2Set: true, player2: 2},
		},
	}

	expected := [][]int{{0, 0}, {1, 1}}
	result := getMatchIndizesOfMatch(schedule, match)
	assert.Equal(t, expected, result)
}

func TestGetMatchIndizesOfPartialMatch(t *testing.T) {
	match := Match{player1: 1, isPlayer2Set: false}
	schedule := [][]Match{
		{
			{player1: 1, isPlayer2Set: false},
			{player1: 3, isPlayer2Set: true, player2: 4},
		},
		{
			{player1: 5, isPlayer2Set: true, player2: 1},
			{player1: 1, isPlayer2Set: false},
		},
	}

	expected := [][]int{{0, 0}, {1, 1}}
	result := getMatchIndizesOfMatch(schedule, match)
	assert.Equal(t, expected, result)
}
