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
	betterStats := computeScheduleStats(betterSchedule, players)
	worseStats := computeScheduleStats(worseSchedule, players)

	valBetter := getStdOfPlayerTimesPlayingStats(betterStats, players)
	valWorse := getStdOfPlayerTimesPlayingStats(worseStats, players)
	assert.LessOrEqual(t, valBetter, valWorse)

	valBetter = getStdOfPossibleMatchesStats(betterStats, players)
	valWorse = getStdOfPossibleMatchesStats(worseStats, players)
	assert.LessOrEqual(t, valBetter, valWorse)

	valBetter = getStdOfPauseBetweenPlayingStats(betterStats, betterSchedule, players)
	valWorse = getStdOfPauseBetweenPlayingStats(worseStats, worseSchedule, players)
	assert.LessOrEqual(t, valBetter, valWorse)

	valBetter = getStdOfPauseBetweenMatchesStats(betterStats, betterSchedule, players)
	valWorse = getStdOfPauseBetweenMatchesStats(worseStats, worseSchedule, players)
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

	betterStats := computeScheduleStats(betterSchedule, players)
	worseStats := computeScheduleStats(worseSchedule, players)

	valBetter := getStdOfPlayerTimesPlayingStats(betterStats, players)
	valWorse := getStdOfPlayerTimesPlayingStats(worseStats, players)
	assert.LessOrEqual(t, valBetter, valWorse)

	valBetter = getStdOfPossibleMatchesStats(betterStats, players)
	valWorse = getStdOfPossibleMatchesStats(worseStats, players)
	assert.LessOrEqual(t, valBetter, valWorse)

	valBetter = getStdOfPauseBetweenPlayingStats(betterStats, betterSchedule, players)
	valWorse = getStdOfPauseBetweenPlayingStats(worseStats, worseSchedule, players)
	assert.LessOrEqual(t, valBetter, valWorse)

	valBetter = getStdOfPauseBetweenMatchesStats(betterStats, betterSchedule, players)
	valWorse = getStdOfPauseBetweenMatchesStats(worseStats, worseSchedule, players)
	assert.GreaterOrEqual(t, valBetter, valWorse)

	valBetter = GetScore(betterSchedule, players)
	valWorse = GetScore(worseSchedule, players)
	assert.LessOrEqual(t, valBetter, valWorse)
}

func TestComputeScheduleStats(t *testing.T) {
	players := getPlayers()
	schedule := getEvenSchedule()
	stats := computeScheduleStats(schedule, players)
	assert.Equal(t, len(players), len(stats.playerCounts))
	assert.Equal(t, len(players), len(stats.playerRounds))
	assert.Equal(t, len(players), len(stats.matchCounts))
	assert.Equal(t, len(players), len(stats.matchRounds))
	total := 0
	for _, c := range stats.playerCounts {
		total += c
	}
	assert.Greater(t, total, 0)
}
