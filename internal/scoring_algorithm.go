package internal

import (
	"math"

	"github.com/montanaflynn/stats"
)

// Scorer caches scoring results and calculates the new score
type Scorer struct {
	isInitalized                   bool
	players                        []Player
	playersPlaying                 []float64
	matchesPlaying                 map[Match]float64
	playersPauses                  []float64
	matchesPauses                  map[Match]float64
	changedPlayerSinceCalculation  []int
	changedMatchesSinceCalculation []Match
}

func (s *Scorer) initialize(schedule [][]Match, players []Player) {
	s.players = players
	s.playersPlaying = make([]float64, len(players))
	s.playersPauses = make([]float64, len(players))
	s.matchesPlaying = make(map[Match]float64)
	s.matchesPauses = make(map[Match]float64)

	// initialize players
	for i := range s.players {
		s.calcPlayersPlaying(schedule, i)
		s.calcPausePlayer(schedule, i)
	}
	// initialize matches
	for i := range players {
		for j := i + 1; j < len(players); j++ {
			match, err := createMatch(i, j)
			if err != nil {
				continue
			}
			s.calcMatchesPlaying(schedule, match)
			s.calcPauseMatch(schedule, match)
		}
	}
	s.isInitalized = true
}

func (s *Scorer) updateCachedData(schedule [][]Match) {
	for _, p := range s.changedPlayerSinceCalculation {
		s.calcPlayersPlaying(schedule, p)
		s.calcPausePlayer(schedule, p)
	}
	s.changedPlayerSinceCalculation = s.changedPlayerSinceCalculation[:0]
	for _, m := range s.changedMatchesSinceCalculation {
		s.calcMatchesPlaying(schedule, m)
		s.calcPauseMatch(schedule, m)
	}
	s.changedMatchesSinceCalculation = s.changedMatchesSinceCalculation[:0]
}

// GetScore gives the value which was used to optimize the schedule
func (s *Scorer) GetScore(schedule [][]Match, players []Player) float64 {
	if !s.isInitalized {
		s.initialize(schedule, players)
	}
	numberRounds := float64(len(schedule))
	var score float64

	s.updateCachedData(schedule)

	val1, err := stats.StandardDeviation(s.playersPlaying)
	if err != nil {
		return math.MaxFloat64
	}
	score += val1 * numberRounds

	val2, err := stats.StandardDeviation(convertToSlice(s.matchesPlaying))
	if err != nil {
		return math.MaxFloat64
	}
	score += val2 * numberRounds

	val3, err := stats.Sum(s.playersPauses)
	if err != nil {
		return math.MaxFloat64
	}
	score += val3

	val4, err := stats.Sum(convertToSlice(s.matchesPauses))
	if err != nil {
		return math.MaxFloat64
	}
	score += val4

	return score
}

func (s *Scorer) calcPlayersPlaying(schedule [][]Match, playerIdx int) {
	s.playersPlaying[playerIdx] = float64(getMatchesCountOfPlayer(schedule, playerIdx)) / s.players[playerIdx].Weight
}

func (s *Scorer) calcMatchesPlaying(schedule [][]Match, match Match) {
	// only works for full matches and is only used there
	combinedWeight := s.players[match.player1].Weight + s.players[match.player2].Weight
	s.matchesPlaying[match] = float64(getCountOfMatchInSchedule(schedule, match)) / combinedWeight
}

func (s *Scorer) calcPausePlayer(schedule [][]Match, playerIdx int) {
	roundsPlaying := getRoundIndizesOfPlayer(schedule, playerIdx)
	s.playersPauses[playerIdx] = calcStdOfPauses(schedule, roundsPlaying)
}

func (s *Scorer) calcPauseMatch(schedule [][]Match, match Match) {
	roundsPlaying := getRoundIndizesOfMatch(schedule, match)
	s.matchesPauses[match] = calcStdOfPauses(schedule, roundsPlaying)
}

func calcStdOfPauses(schedule [][]Match, roundsPlaying []int) float64 {
	if len(roundsPlaying) > 1 {
		pauses := getPausesBetweenRounds(schedule, roundsPlaying)
		val, err := stats.StandardDeviation(pauses)
		if err != nil {
			return float64(len(schedule))
		}
		return val
	}
	return float64(len(schedule))
}

func getPausesBetweenRounds(schedule [][]Match, roundsPlaying []int) []float64 {
	pauses := make([]float64, 0, 8)
	if len(roundsPlaying) > 1 {
		for j := 0; j < len(roundsPlaying)-1; j++ {
			pauses = append(pauses, float64(roundsPlaying[j+1]-roundsPlaying[j]))
		}
		pauses = append(pauses, float64(len(schedule)-roundsPlaying[len(roundsPlaying)-1])) // pause to end
		pauses = append(pauses, float64(roundsPlaying[0]))                                  // pause at start
	}
	return pauses
}

func (s *Scorer) appendChangedPlayers(players ...int) {
	for _, p := range players {
		if !isInSlice(p, s.changedPlayerSinceCalculation) {
			s.changedPlayerSinceCalculation = append(s.changedPlayerSinceCalculation, p)
		}
	}
}

func (s *Scorer) appendChangedMatches(match ...Match) {
	s.changedMatchesSinceCalculation = append(s.changedMatchesSinceCalculation, match...)
	for _, m := range match {
		s.appendChangedPlayers(m.getPlayers()...)
	}
}
