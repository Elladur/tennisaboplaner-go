package internal

import (
	"time"
)

type Season struct {
	Players []Player
	Start time.Time
	End time.Time
	StartTime simpleTime
	EndTime simpleTime
	NumberOfCourts int
	CalendarTitle string
	OverallCosts float64
	ExcludedDates []time.Time
	dates []time.Time
	fixedRounds []int
	Schedule [][]Match
}

type simpleTime struct {
	Hour uint8
	Minute uint8
}

func CreateSeason(players []Player, start time.Time, end time.Time, numberOfCourts int, calendarTitle string, overallCosts float64, excludedDates []time.Time) Season {
	startTime := simpleTime{uint8(start.Hour()), uint8(start.Minute())}
	endTime := simpleTime{uint8(end.Hour()), uint8(end.Minute())}
	start = start.Truncate(24 * time.Hour)
	end = end.Truncate(24 * time.Hour)
	var dates []time.Time
	for d := start; d.Before(end); d = d.AddDate(0, 0, 7) {
		if !isInSlice(d, excludedDates) {
			dates = append(dates, d)
		}
	}

	return Season{
		Players: players,
		Start: start,
		End: end,
		StartTime: startTime,
		EndTime: endTime,
		NumberOfCourts: numberOfCourts,
		CalendarTitle: calendarTitle,
		OverallCosts: overallCosts,
		ExcludedDates: excludedDates,
		dates: dates,
	}
}

func (s *Season) createSchedule() {
	for i, d := range s.dates {
		if isInSlice(d, s.ExcludedDates) {
			continue
		}
		r, partial := s.createRound(d)
		s.Schedule = append(s.Schedule, r)
		if partial {
			s.fixedRounds = append(s.fixedRounds, i)
		}
	}
}

func (s *Season) createRound(date time.Time) ([]Match, bool) {
	var matches []Match
	playerIdx := shuffle(s.getPossiblePlayers(date))

	for len(playerIdx) > 0 || len(matches) <= s.NumberOfCourts {
		switch {
		case len(playerIdx) >= 2:
			playerIdx, p := pop(playerIdx)
			playerIdx, q := pop(playerIdx)
			match, _ := createMatch(p, q)
			matches = append(matches, match)
		case len(playerIdx) == 1:
			_, p := pop(playerIdx)
			match := createPartialMatch(p)
			matches = append(matches, match)
		}
	}

	return matches, len(matches) < s.NumberOfCourts
}

func (s Season) getPossiblePlayers(date time.Time) []uint8 {
	var players []uint8
	for i, p := range s.Players {
		if !isInSlice(date, p.CannotPlay) {
			players = append(players, uint8(i))
		}
	}
	return players
}

func (s *Season) changeMatch(round_index int, match_index int, new_match Match) bool {
	if isInSlice(round_index, s.fixedRounds) {
		return false
	}
	old_match := s.Schedule[round_index][match_index]
	s.Schedule[round_index][match_index] = new_match
	if s.checkIfRoundIsValid(round_index) {
		return true
	} else {
		s.Schedule[round_index][match_index] = old_match
		return false
	}
}

func (s Season) checkIfRoundIsValid(round_index int) bool {
	players := getPlayersOfRound(s.Schedule[round_index])
	if !isInSlice(round_index, s.fixedRounds) && len(players) != 2 * s.NumberOfCourts  {
		return false
	}
	date := s.dates[round_index]
	for _, idx := range players {
		if isInSlice(date, s.Players[idx].CannotPlay) {
			return false
		}
	}
	return true
}

func (s Season) checkIfScheduleIsValid() bool {
	for i := range s.Schedule {
		if !s.checkIfRoundIsValid(i) {
			return false
		}
	}
	return true
}
