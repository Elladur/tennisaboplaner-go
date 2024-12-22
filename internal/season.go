package internal

import (
	"time"
)

type Season struct {
	Players        []Player
	Start          time.Time
	End            time.Time
	StartTime      simpleTime
	EndTime        simpleTime
	NumberOfCourts int
	CalendarTitle  string
	OverallCosts   float64
	ExcludedDates  []time.Time
	dates          []time.Time
	fixedRounds    []int
	Schedule       [][]Match
}

type simpleTime struct {
	Hour   uint8
	Minute uint8
}

type SeasonSettings struct {
	Players        []Player
	Start          string
	End            string
	ExcludedDates  []string
	NumberOfCourts int
	OverallCost    float64
	CalendarTitle  string
}

func CreateSeasonFromSettings(settings SeasonSettings) (Season, error) {
	start, err := time.Parse(time.DateTime, settings.Start)
	if err != nil {
		return Season{}, err
	}
	end, err := time.Parse(time.DateTime, settings.End)
	if err != nil {
		return Season{}, err
	}
	var excludedDates []time.Time
	for _, s := range settings.ExcludedDates {
		t, err := time.Parse(time.DateOnly, s)
		excludedDates = append(excludedDates, t)
		if err != nil {
			return Season{}, err
		}
	}
	return CreateSeason(settings.Players, start, end, settings.NumberOfCourts, settings.CalendarTitle, settings.OverallCost, excludedDates), nil
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
		Players:        players,
		Start:          start,
		End:            end,
		StartTime:      startTime,
		EndTime:        endTime,
		NumberOfCourts: numberOfCourts,
		CalendarTitle:  calendarTitle,
		OverallCosts:   overallCosts,
		ExcludedDates:  excludedDates,
		dates:          dates,
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

func (s *Season) changeMatch(roundIdx int, matchIdx int, newMatch Match) bool {
	if isInSlice(roundIdx, s.fixedRounds) {
		return false
	}
	old_match := s.Schedule[roundIdx][matchIdx]
	s.Schedule[roundIdx][matchIdx] = newMatch
	if s.checkIfRoundIsValid(roundIdx) {
		return true
	} else {
		s.Schedule[roundIdx][matchIdx] = old_match
		return false
	}
}

func (s Season) checkIfRoundIsValid(roundIdx int) bool {
	players := getPlayersOfRound(s.Schedule[roundIdx])
	if !isInSlice(roundIdx, s.fixedRounds) && len(players) != 2*s.NumberOfCourts {
		return false
	}
	date := s.dates[roundIdx]
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

func (s *Season) swapPlayersOfRound(roundIdx int, player1 uint8, player2 uint8) bool {
	players := getPlayersOfRound(s.Schedule[roundIdx])
	if !isInSlice(player1, players) || !isInSlice(player2, players) {
		return false
	}
	for i, m := range s.Schedule[roundIdx] {
		if m.isPlayer2Set {
			switch {
			case m.player1 == player1 && m.player2 == player2:
				return false
			case m.player1 == player2 && m.player2 == player1:
				return false
			case m.player1 == player1:
				match, _ := createMatch(player2, m.player2)
				s.Schedule[roundIdx][i] = match
			case m.player2 == player1:
				match, _ := createMatch(player2, m.player1)
				s.Schedule[roundIdx][i] = match
			case m.player1 == player2:
				match, _ := createMatch(player1, m.player2)
				s.Schedule[roundIdx][i] = match
			case m.player2 == player2:
				match, _ := createMatch(player1, m.player1)
				s.Schedule[roundIdx][i] = match
			}
		} else {
			switch {
			case m.player1 == player1:
				s.Schedule[roundIdx][i] = createPartialMatch(player2)
			case m.player1 == player2:
				s.Schedule[roundIdx][i] = createPartialMatch(player1)
			}
		}
	}
	return true
}

func (s *Season) switchMatches(round1 int, match1 int, round2 int, match2 int) bool {
	if isInSlice(round1, s.fixedRounds) || isInSlice(round2, s.fixedRounds) {
		return false
	}
	match1_old := s.Schedule[round1][match1]
	match2_old := s.Schedule[round2][match2]
	s.Schedule[round1][match1] = match2_old
	s.Schedule[round2][match2] = match1_old
	if s.checkIfRoundIsValid(round1) && s.checkIfRoundIsValid(round2) {
		return true
	} else {
		s.Schedule[round1][match1] = match1_old
		s.Schedule[round2][match2] = match2_old
		return false
	}
}
