package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getPlayers() []Player {
	return []Player{
		{Name: "Player 1", Weight: 1},
		{Name: "Player 2", Weight: 1},
		{Name: "Player 3", Weight: 2},
	}
}

func getScheduleWithOnePlayerNotPlaying() [][]Match {
	return [][]Match{
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
	}
}

func getEvenSchedule() [][]Match {
	return [][]Match{
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 1, player2: 2, isPlayer2Set: true}},
		{{player1: 0, player2: 2, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 1, player2: 2, isPlayer2Set: true}},
		{{player1: 0, player2: 2, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
	}
}

func getBlockSchedule() [][]Match {
	return [][]Match{
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 2, isPlayer2Set: true}},
		{{player1: 0, player2: 2, isPlayer2Set: true}},
		{{player1: 1, player2: 2, isPlayer2Set: true}},
		{{player1: 1, player2: 2, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
	}
}

func getBalancedSchedule() [][]Match {
	return [][]Match{
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 2, isPlayer2Set: true}},
		{{player1: 1, player2: 2, isPlayer2Set: true}},
		{{player1: 0, player2: 2, isPlayer2Set: true}},
		{{player1: 1, player2: 2, isPlayer2Set: true}},
		{{player1: 0, player2: 2, isPlayer2Set: true}},
		{{player1: 1, player2: 2, isPlayer2Set: true}},
	}
}

func helperTestAssertion(t *testing.T, betterSchedule [][]Match, worseSchedule [][]Match, players []Player) {
	valBetter, _ := getStdOfPlayerTimesPlaying(betterSchedule, players)
	valWorse, _ := getStdOfPlayerTimesPlaying(worseSchedule, players)
	assert.LessOrEqual(t, valBetter, valWorse)

	valBetter, _ = getStdOfPossibleMatches(betterSchedule, players)
	valWorse, _ = getStdOfPossibleMatches(worseSchedule, players)
	assert.LessOrEqual(t, valBetter, valWorse)

	valBetter, _ = getStdOfPauseBetweenPlaying(betterSchedule, players)
	valWorse, _ = getStdOfPauseBetweenPlaying(worseSchedule, players)
	assert.LessOrEqual(t, valBetter, valWorse)

	valBetter, _ = getStdOfPauseBetweenMatches(betterSchedule, players)
	valWorse, _ = getStdOfPauseBetweenMatches(worseSchedule, players)
	assert.LessOrEqual(t, valBetter, valWorse)
}

func TestEvenBetterThanNotPlaying(t *testing.T) {
	players := getPlayers()
	betterSchedule := getEvenSchedule()
	worseSchedule := getScheduleWithOnePlayerNotPlaying()

	helperTestAssertion(t, betterSchedule, worseSchedule, players)
}

func TestEvenBetterThanBlock(t *testing.T) {
	players := getPlayers()
	betterSchedule := getEvenSchedule()
	worseSchedule := getBlockSchedule()

	helperTestAssertion(t, betterSchedule, worseSchedule, players)
}

func TestBalancedBetterThanEven(t *testing.T) {
	players := getPlayers()
	betterSchedule := getBalancedSchedule()
	worseSchedule := getEvenSchedule()

	valBetter, _ := getStdOfPlayerTimesPlaying(betterSchedule, players)
	valWorse, _ := getStdOfPlayerTimesPlaying(worseSchedule, players)
	assert.LessOrEqual(t, valBetter, valWorse)

	valBetter, _ = getStdOfPossibleMatches(betterSchedule, players)
	valWorse, _ = getStdOfPossibleMatches(worseSchedule, players)
	assert.LessOrEqual(t, valBetter, valWorse)

	valBetter, _ = getStdOfPauseBetweenPlaying(betterSchedule, players)
	valWorse, _ = getStdOfPauseBetweenPlaying(worseSchedule, players)
	assert.LessOrEqual(t, valBetter, valWorse)

	valBetter, _ = getStdOfPauseBetweenMatches(betterSchedule, players)
	valWorse, _ = getStdOfPauseBetweenMatches(worseSchedule, players)
	assert.GreaterOrEqual(t, valBetter, valWorse)

	valBetter = GetScore(betterSchedule, players)
	valWorse = GetScore(worseSchedule, players)
	assert.LessOrEqual(t, valBetter, valWorse)
}
