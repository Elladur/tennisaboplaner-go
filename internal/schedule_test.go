package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMatchIndizesOfPlayer(t *testing.T) {
	player := uint8(1)
	schedule := [][]Match{
		{
			{Player1: 1, IsPlayer2Set: true, Player2: 2},
			{Player1: 3, IsPlayer2Set: true, Player2: 4},
		},
		{
			{Player1: 5, IsPlayer2Set: true, Player2: 1},
			{Player1: 2, IsPlayer2Set: false},
		},
		{
			{Player1: 5, IsPlayer2Set: true, Player2: 3},
			{Player1: 1, IsPlayer2Set: false},
		},
	}

	expected := [][]int{{0, 0}, {1, 0}, {2, 1}}
	result := GetMatchIndizesOfPlayer(schedule, player)
	assert.Equal(t, expected, result)
}

func TestGetMatchIndizesOfMatch(t *testing.T) {
	match := Match{Player1: 1, IsPlayer2Set: true, Player2: 2}
	schedule := [][]Match{
		{
			{Player1: 1, IsPlayer2Set: true, Player2: 2},
			{Player1: 3, IsPlayer2Set: true, Player2: 4},
		},
		{
			{Player1: 5, IsPlayer2Set: true, Player2: 1},
			{Player1: 1, IsPlayer2Set: true, Player2: 2},
		},
	}

	expected := [][]int{{0, 0}, {1, 1}}
	result := GetMatchIndizesOfMatch(schedule, match)
	assert.Equal(t, expected, result)
}

func TestGetMatchIndizesOfPartialMatch(t *testing.T) {
	match := Match{Player1: 1, IsPlayer2Set: false}
	schedule := [][]Match{
		{
			{Player1: 1, IsPlayer2Set: false},
			{Player1: 3, IsPlayer2Set: true, Player2: 4},
		},
		{
			{Player1: 5, IsPlayer2Set: true, Player2: 1},
			{Player1: 1, IsPlayer2Set: false},
		},
	}

	expected := [][]int{{0, 0}, {1, 1}}
	result := GetMatchIndizesOfMatch(schedule, match)
	assert.Equal(t, expected, result)
}