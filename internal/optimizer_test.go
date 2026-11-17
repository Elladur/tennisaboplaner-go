package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyChangeMatchRejectedKeepsStats(t *testing.T) {
	season := setupStaticTestSeason()
	players := season.Players
	stats := computeScheduleStats(season.Schedule, players)

	// A replacement that duplicates a player in the round must be rejected.
	nm, _ := createMatch(0, 1)
	before := stats
	ok := applyChangeMatch(&stats, &season, 0, 0, nm)
	assert.False(t, ok)
	assert.True(t, statsConsistent(stats, season.Schedule, players))
	assert.True(t, statsConsistent(before, season.Schedule, players))
}

func TestApplySwapMatchesConsistency(t *testing.T) {
	season := setupStaticTestSeason()
	players := season.Players
	stats := computeScheduleStats(season.Schedule, players)

	ok := applySwapMatches(&stats, &season, 0, 0, 1, 0)
	if ok {
		assert.True(t, statsConsistent(stats, season.Schedule, players))
	} else {
		assert.True(t, statsConsistent(stats, season.Schedule, players))
	}
}

func TestApplySwapPlayersOfRoundConsistency(t *testing.T) {
	season := setupStaticTestSeason()
	players := season.Players
	stats := computeScheduleStats(season.Schedule, players)

	roundPlayers := getPlayersOfRound(season.Schedule[0])
	p, q := roundPlayers[0], roundPlayers[2]
	oldRound, ok := applySwapPlayersOfRound(&stats, &season, 0, p, q)
	assert.True(t, ok)
	assert.True(t, statsConsistent(stats, season.Schedule, players))

	// revert and check consistency again
	for _, m := range season.Schedule[0] {
		stats.removeMatch(m, 0)
	}
	for _, m := range oldRound {
		stats.addMatch(m, 0)
	}
	season.swapPlayersOfRound(0, q, p)
	assert.True(t, statsConsistent(stats, season.Schedule, players))
}
