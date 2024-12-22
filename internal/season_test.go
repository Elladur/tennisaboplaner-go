package internal

import (
	"encoding/json"
	"fmt"
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

func TestGetPossiblePlayers(t *testing.T) {
	season := setupStaticTestSeason()
	date := season.dates[0]
	players := season.getPossiblePlayers(date)
	assert.Equal(t, []uint8{0, 1, 2, 3}, players)

	date2 := season.dates[1]
	players = season.getPossiblePlayers(date2)
	assert.Equal(t, 5, len(players))
}

func TestChangeMatch(t *testing.T) {
	season := setupStaticTestSeason()
	newMatch := Match{player1: 3, player2: 2, isPlayer2Set: true}
	success := season.changeMatch(1, 0, newMatch)
	assert.True(t, success)
	assert.Equal(t, newMatch, season.Schedule[1][0])

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
	season := setupStaticTestSeason()
	success := season.swapPlayersOfRound(0, 0, 3)
	assert.True(t, success)
	assert.Equal(t, uint8(2), season.Schedule[0][0].player1)
	assert.Equal(t, uint8(3), season.Schedule[0][0].player2)
	assert.Equal(t, uint8(0), season.Schedule[0][1].player1)
	assert.Equal(t, uint8(1), season.Schedule[0][1].player2)

	success = season.swapPlayersOfRound(0, 3, 2)
	assert.False(t, success)
	// check nothing changed
	assert.Equal(t, uint8(2), season.Schedule[0][0].player1)
	assert.Equal(t, uint8(3), season.Schedule[0][0].player2)
	assert.Equal(t, uint8(0), season.Schedule[0][1].player1)
	assert.Equal(t, uint8(1), season.Schedule[0][1].player2)

	// Test swapping players not in the round
	success = season.swapPlayersOfRound(0, 2, 4)
	assert.False(t, success)
	// check nothing changed
	assert.Equal(t, uint8(2), season.Schedule[0][0].player1)
	assert.Equal(t, uint8(3), season.Schedule[0][0].player2)
	assert.Equal(t, uint8(0), season.Schedule[0][1].player1)
	assert.Equal(t, uint8(1), season.Schedule[0][1].player2)
}

func TestSwitchMatches(t *testing.T) {
	season := setupStaticTestSeason()
	oldMatch1 := season.Schedule[1][0]
	oldMatch2 := season.Schedule[2][1]
	success := season.switchMatches(1, 0, 2, 1)
	assert.True(t, success)
	assert.Equal(t, season.Schedule[1][0], oldMatch2)
	assert.Equal(t, season.Schedule[2][1], oldMatch1)

	oldMatch1 = season.Schedule[0][0]
	oldMatch2 = season.Schedule[1][1]
	success = season.switchMatches(0, 0, 1, 1)
	assert.False(t, success)
	assert.Equal(t, season.Schedule[0][0], oldMatch1)
	assert.Equal(t, season.Schedule[1][1], oldMatch2)

	// Test switching matches in fixed rounds
	season.fixedRounds = append(season.fixedRounds, 1)
	success = season.switchMatches(1, 0, 2, 1)
	assert.False(t, success)
}

func TestCheckIfRoundIsValid(t *testing.T) {
	season := setupTestSeason()
	assert.True(t, season.checkIfRoundIsValid(0))

	invalidRound := []Match{{player1: 0, player2: 2, isPlayer2Set: true}, {player1: 0, player2: 2, isPlayer2Set: true}}
	season.Schedule[0] = invalidRound
	assert.False(t, season.checkIfRoundIsValid(0))

	season.Schedule[0] = []Match{season.Schedule[0][0]}
	assert.False(t, season.checkIfRoundIsValid(0))

	// player 4 cant play in first round
	invalidRound = []Match{{player1: 0, player2: 2, isPlayer2Set: true}, {player1: 3, player2: 4, isPlayer2Set: true}}
	season.Schedule[0] = invalidRound
	assert.False(t, season.checkIfRoundIsValid(0))
}

func TestCheckIfScheduleIsValid(t *testing.T) {
	season := setupTestSeason()
	assert.True(t, season.checkIfScheduleIsValid())

	invalidRound := []Match{{player1: 0, player2: 2, isPlayer2Set: true}, {player1: 0, player2: 2, isPlayer2Set: true}}
	season.Schedule[5] = invalidRound
	assert.False(t, season.checkIfScheduleIsValid())
}

func setupTestSeason() Season {
	players := []Player{{Name: "Player1", Weight: 1}, {Name: "Player2", Weight: 1}, {Name: "Player3", Weight: 1}, {Name: "Player4", Weight: 1}, {Name: "Player5", Weight: 1, CannotPlay: []time.Time{time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)}}}
	start := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	end := time.Date(2023, 2, 31, 12, 0, 0, 0, time.UTC)
	excludedDates := []time.Time{time.Date(2023, 1, 29, 0, 0, 0, 0, time.UTC)}

	season := createSeason(players, start, end, 2, "Test Season", 100.0, excludedDates)
	//data, err := json.Marshal(season)
	//if err != nil {
	//fmt.Println(err)
	//}
	//datastr := string(data)
	//fmt.Println(datastr)
	return season
}

func setupStaticTestSeason() Season {
	data := "{\"Players\":[{\"CannotPlay\":null,\"Name\":\"Player1\",\"Weight\":1},{\"CannotPlay\":null,\"Name\":\"Player2\",\"Weight\":1},{\"CannotPlay\":null,\"Name\":\"Player3\",\"Weight\":1},{\"CannotPlay\":null,\"Name\":\"Player4\",\"Weight\":1},{\"CannotPlay\":[\"2023-01-01\"],\"Name\":\"Player5\",\"Weight\":1}],\"Start\":\"2023-01-01T00:00:00Z\",\"End\":\"2023-03-03T00:00:00Z\",\"StartTime\":{\"Hour\":10,\"Minute\":0},\"EndTime\":{\"Hour\":12,\"Minute\":0},\"NumberOfCourts\":2,\"CalendarTitle\":\"Test Season\",\"OverallCosts\":100,\"ExcludedDates\":[\"2023-01-29T00:00:00Z\"],\"Schedule\":[[{\"IsPlayer2Set\":true,\"Player1\":0,\"Player2\":2},{\"IsPlayer2Set\":true,\"Player1\":1,\"Player2\":3}],[{\"IsPlayer2Set\":true,\"Player1\":3,\"Player2\":4},{\"IsPlayer2Set\":true,\"Player1\":0,\"Player2\":1}],[{\"IsPlayer2Set\":true,\"Player1\":0,\"Player2\":1},{\"IsPlayer2Set\":true,\"Player1\":2,\"Player2\":4}],[{\"IsPlayer2Set\":true,\"Player1\":1,\"Player2\":2},{\"IsPlayer2Set\":true,\"Player1\":0,\"Player2\":3}],[{\"IsPlayer2Set\":true,\"Player1\":1,\"Player2\":3},{\"IsPlayer2Set\":true,\"Player1\":0,\"Player2\":2}],[{\"IsPlayer2Set\":true,\"Player1\":2,\"Player2\":3},{\"IsPlayer2Set\":true,\"Player1\":1,\"Player2\":4}],[{\"IsPlayer2Set\":true,\"Player1\":0,\"Player2\":2},{\"IsPlayer2Set\":true,\"Player1\":1,\"Player2\":4}],[{\"IsPlayer2Set\":true,\"Player1\":0,\"Player2\":3},{\"IsPlayer2Set\":true,\"Player1\":1,\"Player2\":4}]]}"
	season := Season{}
	err := json.Unmarshal([]byte(data), &season)
	if err != nil {
		fmt.Println(err)
	}
	season.dates = generateDates(season.Start, season.End, season.ExcludedDates)
	return season
}
