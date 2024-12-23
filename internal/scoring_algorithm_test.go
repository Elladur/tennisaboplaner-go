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
	var betterScorer Scorer
	betterScorer.initialize(betterSchedule, players)
	var worseScorer Scorer
	worseScorer.initialize(worseSchedule, players)

	valBetter := betterScorer.GetScore(betterSchedule, players)
	valWorse := worseScorer.GetScore(worseSchedule, players)
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
	var betterScorer Scorer
	betterScorer.initialize(betterSchedule, players)
	var worseScorer Scorer
	worseScorer.initialize(worseSchedule, players)

	valBetter := betterScorer.GetScore(betterSchedule, players)
	valWorse := worseScorer.GetScore(worseSchedule, players)
	assert.LessOrEqual(t, valBetter, valWorse)
}
