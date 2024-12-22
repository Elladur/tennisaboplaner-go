package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateSeasonFromSettings(t *testing.T) {
	settings := SeasonSettings{
		Players:        []Player{{Name: "Player1"}, {Name: "Player2"}},
		Start:          "2023-01-01 10:00:00",
		End:            "2023-12-31 12:00:00",
		ExcludedDates:  []string{"2023-06-01"},
		NumberOfCourts: 2,
		OverallCost:    100.0,
		CalendarTitle:  "Test Season",
	}

	season, err := CreateSeasonFromSettings(settings)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(season.Players))
	assert.Equal(t, "Test Season", season.CalendarTitle)
	assert.Equal(t, 100.0, season.OverallCosts)
	assert.Equal(t, 2, season.NumberOfCourts)
	assert.Equal(t, 1, len(season.ExcludedDates))

	// Test with invalid start date
	settings.Start = "invalid-date"
	_, err = CreateSeasonFromSettings(settings)
	assert.Error(t, err)

	// Test with invalid end date
	settings.Start = "2023-01-01 10:00:00"
	settings.End = "invalid-date"
	_, err = CreateSeasonFromSettings(settings)
	assert.Error(t, err)

	// Test with invalid excluded date
	settings.End = "2023-12-31 10:00:00"
	settings.ExcludedDates = []string{"invalid-date"}
	_, err = CreateSeasonFromSettings(settings)
	assert.Error(t, err)
}

func TestCreateSeason(t *testing.T) {
	players := []Player{{Name: "Player1"}, {Name: "Player2"}}
	start := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	end := time.Date(2023, 12, 31, 12, 0, 0, 0, time.UTC)
	excludedDates := []time.Time{time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)}

	season := createSeason(players, start, end, 2, "Test Season", 100.0, excludedDates)
	assert.Equal(t, 2, len(season.Players))
	assert.Equal(t, "Test Season", season.CalendarTitle)
	assert.Equal(t, 100.0, season.OverallCosts)
	assert.Equal(t, 2, season.NumberOfCourts)
	assert.Equal(t, 1, len(season.ExcludedDates))
	assert.Equal(t, simpleTime{10, 0}, season.StartTime)
	assert.Equal(t, simpleTime{12, 0}, season.EndTime)
	assert.Equal(t, time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC), season.Start)
	assert.Equal(t, time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC), season.End)
	assert.NotEmpty(t, season.dates)
	assert.NotEmpty(t, season.Schedule)
}

func TestCreateSchedule(t *testing.T) {
	season := setupTestSeason()
	season.createSchedule()
	assert.NotEmpty(t, season.Schedule)
	assert.Equal(t, len(season.dates), len(season.Schedule))
}

func TestCreateRound(t *testing.T) {
	season := setupTestSeason()
	date := season.dates[0]
	matches, partial := season.createRound(date)
	assert.NotEmpty(t, matches)
	assert.False(t, partial)

	// Test with insufficient players
	season.Players = []Player{{Name: "Player1"}}
	matches, partial = season.createRound(date)
	assert.NotEmpty(t, matches)
	assert.True(t, partial)
}

func TestChangeMatch(t *testing.T) {
	season := setupTestSeason()
	season.createSchedule()
	newMatch := Match{player1: 0, player2: 1, isPlayer2Set: true}
	success := season.changeMatch(0, 0, newMatch)
	assert.True(t, success)
	assert.Equal(t, newMatch, season.Schedule[0][0])

	invalidMatch := season.Schedule[0][1]
	success = season.changeMatch(0, 0, invalidMatch)
	assert.False(t, success)
	assert.NotEqual(t, invalidMatch, season.Schedule[0][0])

	// Test changing match in a fixed round
	season.fixedRounds = append(season.fixedRounds, 0)
	success = season.changeMatch(0, 0, newMatch)
	assert.False(t, success)
}

func TestSwapPlayersOfRound(t *testing.T) {
	season := setupTestSeason()
	season.createSchedule()
	success := season.swapPlayersOfRound(0, 0, 1)
	assert.True(t, success)
	assert.Equal(t, uint8(1), season.Schedule[0][0].player1)
	assert.Equal(t, uint8(0), season.Schedule[0][0].player2)

	success = season.swapPlayersOfRound(0, 0, 2)
	assert.False(t, success)

	// Test swapping players not in the round
	success = season.swapPlayersOfRound(0, 2, 3)
	assert.False(t, success)
}

func TestSwitchMatches(t *testing.T) {
	season := setupTestSeason()
	season.createSchedule()
	success := season.switchMatches(0, 0, 1, 1)
	assert.True(t, success)
	assert.Equal(t, season.Schedule[0][0], season.Schedule[1][1])
	assert.Equal(t, season.Schedule[1][1], season.Schedule[0][0])

	success = season.switchMatches(0, 0, 1, 2)
	assert.False(t, success)

	// Test switching matches in fixed rounds
	season.fixedRounds = append(season.fixedRounds, 0)
	success = season.switchMatches(0, 0, 1, 1)
	assert.False(t, success)
}

func TestCheckIfRoundIsValid(t *testing.T) {
	season := setupTestSeason()
	season.createSchedule()
	assert.True(t, season.checkIfRoundIsValid(0))

	invalidMatch := Match{player1: 0, player2: 2, isPlayer2Set: true}
	season.Schedule[0][0] = invalidMatch
	assert.False(t, season.checkIfRoundIsValid(0))
}

func TestCheckIfScheduleIsValid(t *testing.T) {
	season := setupTestSeason()
	season.createSchedule()
	assert.True(t, season.checkIfScheduleIsValid())

	invalidMatch := Match{player1: 0, player2: 2, isPlayer2Set: true}
	season.Schedule[0][0] = invalidMatch
	assert.False(t, season.checkIfScheduleIsValid())
}

func setupTestSeason() Season {
	players := []Player{{Name: "Player1"}, {Name: "Player2"}, {Name: "Player3"}, {Name: "Player4"}}
	start := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	end := time.Date(2023, 12, 31, 12, 0, 0, 0, time.UTC)
	excludedDates := []time.Time{time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)}

	season := createSeason(players, start, end, 2, "Test Season", 100.0, excludedDates)
	return season
}
