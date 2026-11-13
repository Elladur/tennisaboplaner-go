package internal

import (
	"sort"
	"time"
)

// Season manages all the necessary information to create a schedule for an abo
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
	location       time.Location
}

type simpleTime struct {
	Hour   int
	Minute int
}

// SeasonSettings are necessary to create a Season
type SeasonSettings struct {
	Players        []Player
	Start          string
	End            string
	Location       string
	ExcludedDates  []string
	NumberOfCourts int
	OverallCost    float64
	CalendarTitle  string
}

// CreateSeasonFromSettings creates a Season from the given settings
func CreateSeasonFromSettings(settings SeasonSettings) (Season, error) {
	location, err := time.LoadLocation(settings.Location)
	if err != nil {
		return Season{}, err
	}

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
		if err != nil {
			return Season{}, err
		}
		excludedDates = append(excludedDates, t)
	}
	season := createSeason(settings.Players, start, end, settings.NumberOfCourts, settings.CalendarTitle, settings.OverallCost, excludedDates, *location)
	return season, nil
}

func createSeason(players []Player, start time.Time, end time.Time, numberOfCourts int, calendarTitle string, overallCosts float64, excludedDates []time.Time, location time.Location) Season {
	startTime := simpleTime{start.Hour(), start.Minute()}
	endTime := simpleTime{end.Hour(), end.Minute()}
	start = start.Truncate(24 * time.Hour)
	end = end.Truncate(24 * time.Hour)
	dates := generateDates(start, end.Add(time.Duration(endTime.Hour)*time.Hour+time.Duration(endTime.Minute)*time.Minute), excludedDates)

	season := Season{
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
		location:       location,
	}
	season.CreateSchedule()
	return season
}

func generateDates(start time.Time, end time.Time, excludedDates []time.Time) []time.Time {
	var dates []time.Time
	for d := start; d.Before(end); d = d.AddDate(0, 0, 7) {
		if !isInSlice(d, excludedDates) {
			dates = append(dates, d)
		}
	}
	return dates
}

// CreateSchedule createas new schedule for this season.
// This is done with some randomness, such that optimization has a different starting point.
func (s *Season) CreateSchedule() {
	s.Schedule = [][]Match{}
	for i := range s.dates {
		r, partial := s.createRound(i)
		s.Schedule = append(s.Schedule, r)
		if partial {
			s.fixedRounds = append(s.fixedRounds, i)
		}
	}
}

func (s *Season) createRound(index int) ([]Match, bool) {
	possible := s.getPossiblePlayers(index)

	// cumulative play counts up to previous rounds
	counts := make([]int, len(s.Players))
	for r := 0; r < index && r < len(s.Schedule); r++ {
		for _, m := range s.Schedule[r] {
			if m.isPlayer2Set {
				counts[m.player1]++
				counts[m.player2]++
			} else {
				counts[m.player1]++
			}
		}
	}

	// shuffle for diversity, then stable-sort by play count (greedy balance)
	shuffle(possible)
	type cand struct {
		p     int
		count int
		order int
	}
	cands := make([]cand, 0, len(possible))
	for i, p := range possible {
		cands = append(cands, cand{p: p, count: counts[p], order: i})
	}
	sort.SliceStable(cands, func(a, b int) bool {
		if cands[a].count != cands[b].count {
			return cands[a].count < cands[b].count
		}
		return cands[a].order < cands[b].order
	})
	playerIdx := make([]int, 0, len(cands))
	for _, c := range cands {
		playerIdx = append(playerIdx, c.p)
	}

	var matches []Match
	for len(playerIdx) > 0 && len(matches) < s.NumberOfCourts {
		switch {
		case len(playerIdx) >= 2:
			var p, q int
			playerIdx, p = pop(playerIdx)
			playerIdx, q = pop(playerIdx)
			match, _ := createMatch(p, q)
			matches = append(matches, match)
		case len(playerIdx) == 1:
			var p int
			playerIdx, p = pop(playerIdx)
			match := createPartialMatch(p)
			matches = append(matches, match)
		}
	}

	return matches, len(getPlayersOfRound(matches)) < s.NumberOfCourts*2
}

func (s *Season) getPossiblePlayers(index int) []int {
	date := s.dates[index]
	possiblePlayers := make([]int, 0, len(s.Players))
	for i, p := range s.Players {
		if !isInSlice(date, p.CannotPlay) {
			possiblePlayers = append(possiblePlayers, i)
		}
	}
	return possiblePlayers
}

func (s *Season) changeMatch(roundIdx int, matchIdx int, newMatch Match) bool {
	oldMatch := s.Schedule[roundIdx][matchIdx]
	s.Schedule[roundIdx][matchIdx] = newMatch
	if s.checkIfRoundIsValid(roundIdx) {
		return true
	}
	s.Schedule[roundIdx][matchIdx] = oldMatch
	return false
}

func (s *Season) checkIfRoundIsValid(roundIdx int) bool {
	round := s.Schedule[roundIdx]
	players := getPlayersOfRound(round)
	// distinct players: number of unique players must equal total slots
	slots := 0
	for _, m := range round {
		slots += len(m.getPlayers())
	}
	if len(players) != slots {
		return false
	}
	// possible already excludes players who cannot play on this date.
	possible := s.getPossiblePlayers(roundIdx)
	for _, idx := range players {
		if !isInSlice(idx, possible) {
			return false
		}
	}
	if isInSlice(roundIdx, s.fixedRounds) {
		// fixed/partial round: filled as much as possible — match count equals
		// ceil(availablePlayers/2) capped by courts.
		maxMatches := (len(possible) + 1) / 2
		if maxMatches > s.NumberOfCourts {
			maxMatches = s.NumberOfCourts
		}
		if len(round) != maxMatches {
			return false
		}
	} else if len(players) != 2*s.NumberOfCourts {
		return false
	}
	return true
}

func (s *Season) checkIfScheduleIsValid() bool {
	for i := range s.Schedule {
		if !s.checkIfRoundIsValid(i) {
			return false
		}
	}
	return true
}

func (s *Season) swapPlayersOfRound(roundIdx int, player1 int, player2 int) bool {
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
			switch m.player1 {
			case player1:
				s.Schedule[roundIdx][i] = createPartialMatch(player2)
			case player2:
				s.Schedule[roundIdx][i] = createPartialMatch(player1)
			}
		}
	}
	return true
}

func (s *Season) swapMatches(round1 int, match1 int, round2 int, match2 int) bool {
	oldMatch1 := s.Schedule[round1][match1]
	oldMatch2 := s.Schedule[round2][match2]
	s.Schedule[round1][match1] = oldMatch2
	s.Schedule[round2][match2] = oldMatch1
	if s.checkIfRoundIsValid(round1) && s.checkIfRoundIsValid(round2) {
		return true
	}
	s.Schedule[round1][match1] = oldMatch1
	s.Schedule[round2][match2] = oldMatch2
	return false
}

func (s *Season) replaceRound(index int, round []Match) ([]Match, bool) {
	if isInSlice(index, s.fixedRounds) {
		return nil, false
	}
	oldRound := s.Schedule[index]
	s.Schedule[index] = round
	if s.checkIfRoundIsValid(index) {
		return oldRound, true
	}
	s.Schedule[index] = oldRound
	return nil, false
}
