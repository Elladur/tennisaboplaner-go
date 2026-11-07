package internal

import (
	"math"

	"github.com/montanaflynn/stats"
)

// scheduleStats holds aggregates computed in a single pass over the schedule,
// so the four metric functions do not each re-scan the whole schedule.
type scheduleStats struct {
	playerCounts []int
	playerRounds [][]int
	matchCounts  [][]int
	matchRounds  [][][]int
}

func computeScheduleStats(schedule [][]Match, players []Player) scheduleStats {
	n := len(players)
	s := scheduleStats{
		playerCounts: make([]int, n),
		playerRounds: make([][]int, n),
		matchCounts:  make([][]int, n),
		matchRounds:  make([][][]int, n),
	}
	for i := range s.matchCounts {
		s.matchCounts[i] = make([]int, n)
		s.matchRounds[i] = make([][]int, n)
	}
	for r, round := range schedule {
		for _, m := range round {
			if m.isPlayer2Set {
				p1, p2 := m.player1, m.player2
				s.playerCounts[p1]++
				s.playerCounts[p2]++
				s.playerRounds[p1] = append(s.playerRounds[p1], r)
				s.playerRounds[p2] = append(s.playerRounds[p2], r)
				s.matchCounts[p1][p2]++
				s.matchRounds[p1][p2] = append(s.matchRounds[p1][p2], r)
			} else {
				p := m.player1
				s.playerCounts[p]++
				s.playerRounds[p] = append(s.playerRounds[p], r)
			}
		}
	}
	return s
}

// scoreWeights controls the relative importance of the four scoring components.
type scoreWeights struct {
	PlayTime   float64 `json:"PlayTime"`
	Matchup    float64 `json:"Matchup"`
	PausePlay  float64 `json:"PausePlay"`
	PauseMatch float64 `json:"PauseMatch"`
}

// defaultScoreWeights returns the default weights (play time dominant, matchup
// slightly above pause regularity).
func defaultScoreWeights() scoreWeights {
	return scoreWeights{PlayTime: 3, Matchup: 1.5, PausePlay: 1, PauseMatch: 1}
}

// normalizeWeights replaces zero (unset) values with defaults.
func normalizeWeights(w scoreWeights) scoreWeights {
	d := defaultScoreWeights()
	if w.PlayTime <= 0 {
		w.PlayTime = d.PlayTime
	}
	if w.Matchup <= 0 {
		w.Matchup = d.Matchup
	}
	if w.PausePlay <= 0 {
		w.PausePlay = d.PausePlay
	}
	if w.PauseMatch <= 0 {
		w.PauseMatch = d.PauseMatch
	}
	return w
}
func GetScore(schedule [][]Match, players []Player) float64 {
	return getScoreWithWeights(schedule, players, defaultScoreWeights())
}

// getScoreWithWeights is like GetScore but uses the provided weights.
func getScoreWithWeights(schedule [][]Match, players []Player, weights scoreWeights) float64 {
	weights = normalizeWeights(weights)
	statsData := computeScheduleStats(schedule, players)

	val1 := getStdOfPlayerTimesPlayingStats(statsData, players)
	val2 := getStdOfPossibleMatchesStats(statsData, players)
	val3 := getStdOfPauseBetweenPlayingStats(statsData, schedule, players)
	val4 := getStdOfPauseBetweenMatchesStats(statsData, schedule, players)

	return weights.PlayTime*val1 + weights.Matchup*val2 + weights.PausePlay*val3 + weights.PauseMatch*val4
}

func getStdOfPlayerTimesPlayingStats(s scheduleStats, players []Player) float64 {
	weightedTimes := make([]float64, len(players))
	for i, p := range players {
		weightedTimes[i] = float64(s.playerCounts[i]) / p.Weight
	}
	return getCoefficientOfVariation(weightedTimes)
}

func getStdOfPossibleMatchesStats(s scheduleStats, players []Player) float64 {
	var weightedMatches []float64
	for i := range players {
		for j := i + 1; j < len(players); j++ {
			combinedWeight := players[i].Weight + players[j].Weight
			match, err := createMatch(i, j)
			if err != nil {
				continue
			}
			val := float64(s.matchCounts[match.player1][match.player2]) / combinedWeight
			weightedMatches = append(weightedMatches, val)
		}
	}
	return getCoefficientOfVariation(weightedMatches)
}

func getStdOfPauseBetweenPlayingStats(s scheduleStats, schedule [][]Match, players []Player) float64 {
	pausesBetweenPlaying := make([]float64, 0, len(players))
	for i := range players {
		val := calcCoefficientOfVariationOfPauses(schedule, s.playerRounds[i])
		pausesBetweenPlaying = append(pausesBetweenPlaying, val)
	}
	val, err := stats.Mean(pausesBetweenPlaying)
	if err != nil {
		return math.MaxFloat64
	}
	return val
}

func getStdOfPauseBetweenMatchesStats(s scheduleStats, schedule [][]Match, players []Player) float64 {
	pausesBetweenMatches := make([]float64, 0, len(players)*(len(players)-1)/2)
	for i := range players {
		for j := i + 1; j < len(players); j++ {
			match, err := createMatch(i, j)
			if err != nil {
				continue
			}
			val := calcCoefficientOfVariationOfPauses(schedule, s.matchRounds[match.player1][match.player2])
			pausesBetweenMatches = append(pausesBetweenMatches, val)
		}
	}
	val, err := stats.Mean(pausesBetweenMatches)
	if err != nil {
		return math.MaxFloat64
	}
	return val
}

func getPausesBetweenRounds(schedule [][]Match, roundsPlaying []int) []float64 {
	pauses := make([]float64, 0, 16)
	if len(roundsPlaying) > 1 {
		for j := 0; j < len(roundsPlaying)-1; j++ {
			pauses = append(pauses, float64(roundsPlaying[j+1]-roundsPlaying[j]))
		}
		pauses = append(pauses, float64(len(schedule)-roundsPlaying[len(roundsPlaying)-1])) // pause to end
		pauses = append(pauses, float64(roundsPlaying[0]))                                  // pause at start
	}
	return pauses
}

func calcCoefficientOfVariationOfPauses(schedule [][]Match, roundsPlaying []int) float64 {
	if len(roundsPlaying) > 1 {
		pauses := getPausesBetweenRounds(schedule, roundsPlaying)
		return getCoefficientOfVariation(pauses)
	}
	return 1 // values are positive and have a max value, therefore Coefficient of Variation is most of the times smaller than 1
}

func getCoefficientOfVariation(input []float64) float64 {
	std, err := stats.StandardDeviation(input)
	mean, err2 := stats.Mean(input)
	if err != nil || err2 != nil || mean == 0 {
		return math.MaxFloat64
	}
	return std / mean
}
