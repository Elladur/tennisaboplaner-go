package internal

import (
	"math"

	"github.com/montanaflynn/stats"
)

// GetScore gives the value which was used to optimize the schedule
func GetScore(schedule [][]Match, players []Player) float64 {
	numberRounds := float64(len(schedule))
	var score float64

	val1, err := getStdOfPlayerTimesPlaying(schedule, players)
	if err != nil {
		return math.MaxFloat64
	}
	score += val1 * numberRounds

	val2, err := getStdOfPossibleMatches(schedule, players)
	if err != nil {
		return math.MaxFloat64
	}
	score += val2 * numberRounds

	val3, err := getStdOfPauseBetweenPlaying(schedule, players)
	if err != nil {
		return math.MaxFloat64
	}
	score += val3

	val4, err := getStdOfPauseBetweenMatches(schedule, players)
	if err != nil {
		return math.MaxFloat64
	}
	score += val4

	return score
}

func getStdOfPlayerTimesPlaying(schedule [][]Match, players []Player) (float64, error) {
	var weightedTimesPlayer []float64
	for i, p := range players {
		val := float64(getMatchesCountOfPlayer(schedule, i)) / p.Weight
		weightedTimesPlayer = append(weightedTimesPlayer, val)
	}
	return stats.StandardDeviation(weightedTimesPlayer)
}

func getStdOfPossibleMatches(schedule [][]Match, players []Player) (float64, error) {
	var weightedPossibleMatches []float64
	for i := range players {
		for j := i + 1; j < len(players); j++ {
			combinedWeight := players[i].Weight + players[j].Weight
			match, err := createMatch(i, j)
			if err != nil {
				continue
			}
			val := float64(getCountOfMatchInSchedule(schedule, match)) / combinedWeight
			weightedPossibleMatches = append(weightedPossibleMatches, val)
		}
	}
	return stats.StandardDeviation(weightedPossibleMatches)
}

func getStdOfPauseBetweenPlaying(schedule [][]Match, players []Player) (float64, error) {
	pausesBetweenPlaying := []float64{}
	for i := range players {
		roundsPlaying := getRoundIndizesOfPlayer(schedule, i)
		val := calcStdOfPauses(schedule, roundsPlaying)
		pausesBetweenPlaying = append(pausesBetweenPlaying, val)
	}
	return stats.Sum(pausesBetweenPlaying)
}

func getStdOfPauseBetweenMatches(schedule [][]Match, players []Player) (float64, error) {
	pausesBetweenMatches := []float64{}
	for i := range players {
		for j := i + 1; j < len(players); j++ {
			match, err := createMatch(i, j)
			if err != nil {
				continue
			}
			roundsPlaying := getRoundIndizesOfMatch(schedule, match)
			val := calcStdOfPauses(schedule, roundsPlaying)
			pausesBetweenMatches = append(pausesBetweenMatches, val)
		}
	}
	return stats.Sum(pausesBetweenMatches)
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
