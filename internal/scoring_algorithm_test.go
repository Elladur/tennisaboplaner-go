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

func getScheduleWithOnePlayerNotPlaying() Schedule {
	return Schedule{
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
	}
}

func getEvenSchedule() Schedule {
	return Schedule{
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 1, player2: 2, isPlayer2Set: true}},
		{{player1: 0, player2: 2, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 1, player2: 2, isPlayer2Set: true}},
		{{player1: 0, player2: 2, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
	}
}

func getBlockSchedule() Schedule {
	return Schedule{
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 2, isPlayer2Set: true}},
		{{player1: 0, player2: 2, isPlayer2Set: true}},
		{{player1: 1, player2: 2, isPlayer2Set: true}},
		{{player1: 1, player2: 2, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
	}
}

func getBalancedSchedule() Schedule {
	return Schedule{
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 2, isPlayer2Set: true}},
		{{player1: 1, player2: 2, isPlayer2Set: true}},
		{{player1: 0, player2: 2, isPlayer2Set: true}},
		{{player1: 1, player2: 2, isPlayer2Set: true}},
		{{player1: 0, player2: 2, isPlayer2Set: true}},
		{{player1: 1, player2: 2, isPlayer2Set: true}},
	}
}

func helperTestAssertion(t *testing.T, betterSchedule Schedule, worseSchedule Schedule, players []Player) {
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

func TestIncrementalScoreEqualsFull(t *testing.T) {
	players := getPlayers()
	schedule := getEvenSchedule()
	full := GetScore(schedule, players)
	stats := computeScheduleStats(schedule, players)
	fromStats := getScoreFromStats(stats, schedule, players, defaultScoreWeights())
	assert.InDelta(t, full, fromStats, 1e-9)
	assert.True(t, statsConsistent(stats, schedule, players))

	ops := []struct{ r, m, p1, p2 int }{
		{0, 0, 0, 2}, {1, 0, 0, 1}, {2, 0, 1, 2},
	}
	for _, op := range ops {
		if op.r >= len(schedule) || op.m >= len(schedule[op.r]) {
			continue
		}
		old := schedule[op.r][op.m]
		nm, err := createMatch(op.p1, op.p2)
		if err != nil {
			continue
		}
		schedule[op.r][op.m] = nm
		stats.removeMatch(old, op.r)
		stats.addMatch(nm, op.r)
		assert.InDelta(t, GetScore(schedule, players), getScoreFromStats(stats, schedule, players, defaultScoreWeights()), 1e-9)
		assert.True(t, statsConsistent(stats, schedule, players))
		schedule[op.r][op.m] = old
		stats.removeMatch(nm, op.r)
		stats.addMatch(old, op.r)
		assert.True(t, statsConsistent(stats, schedule, players))
	}

	if len(schedule) > 1 && len(schedule[0]) > 0 && len(schedule[1]) > 0 {
		m1, m2 := schedule[0][0], schedule[1][0]
		schedule[0][0], schedule[1][0] = m2, m1
		stats.removeMatch(m1, 0)
		stats.removeMatch(m2, 1)
		stats.addMatch(m2, 0)
		stats.addMatch(m1, 1)
		assert.InDelta(t, GetScore(schedule, players), getScoreFromStats(stats, schedule, players, defaultScoreWeights()), 1e-9)
		assert.True(t, statsConsistent(stats, schedule, players))
	}
}

func TestInsertSortedAndRemoveRoundIndex(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, insertSorted([]int{1, 3}, 2))
	assert.Equal(t, []int{0}, insertSorted(nil, 0))
	assert.Equal(t, []int{1, 2, 3, 3}, insertSorted([]int{1, 2, 3}, 3))
	// remove first occurrence only
	assert.Equal(t, []int{1, 3}, removeRoundIndex([]int{1, 2, 3}, 2))
	assert.Equal(t, []int{1, 2, 3}, removeRoundIndex([]int{1, 2, 3}, 9))
	assert.Nil(t, removeRoundIndex([]int{5}, 5))
}

func TestIncrementalStatsSwapPlayersOfRound(t *testing.T) {
	season := setupStaticTestSeason()
	if len(season.Schedule) == 0 {
		t.Skip("no schedule")
	}
	players := season.Players
	stats := computeScheduleStats(season.Schedule, players)
	assert.True(t, statsConsistent(stats, season.Schedule, players))

	roundIdx := 0
	roundPlayers := getPlayersOfRound(season.Schedule[roundIdx])
	if len(roundPlayers) < 2 {
		t.Skip("need at least two players in round")
	}
	p, q := roundPlayers[0], roundPlayers[2]
	oldRound, swapped := applySwapPlayersOfRound(&stats, &season, roundIdx, p, q)
	assert.True(t, swapped)
	assert.True(t, statsConsistent(stats, season.Schedule, players))

	// revert
	for _, m := range season.Schedule[roundIdx] {
		stats.removeMatch(m, roundIdx)
	}
	for _, m := range oldRound {
		stats.addMatch(m, roundIdx)
	}
	season.swapPlayersOfRound(roundIdx, q, p)
	assert.True(t, statsConsistent(stats, season.Schedule, players))
}

func TestConsecutivePlayDroughtPenalty(t *testing.T) {
	// player plays rounds 0,1,2 -> consecutive penalties
	schedule := Schedule{
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 0, player2: 1, isPlayer2Set: true}},
		{{player1: 2, player2: 3, isPlayer2Set: true}},
	}
	players := []Player{{Name: "A"}, {Name: "B"}, {Name: "C"}, {Name: "D"}}
	st := computeScheduleStats(schedule, players)
	p := getConsecutivePlayDroughtPenalty(st, schedule)
	assert.Greater(t, p, 0.0)
}
