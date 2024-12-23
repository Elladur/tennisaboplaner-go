package internal

import (
	"math"

	"github.com/montanaflynn/stats"
)

// GetScore gives the value which was used to optimize the schedule
func GetScore(schedule [][]Match, players []Player) float64 {
	numberRounds := float64(len(schedule))

	val1 := getStdOfPlayerTimesPlaying(schedule, players, numberRounds)
	val2 := getStdOfPossibleMatches(schedule, players, numberRounds)
	val3 := getStdOfPauseBetweenPlaying(schedule, players)
	val4 := getStdOfPauseBetweenMatches(schedule, players)

	return val1 + val2 + val3 + val4
}

func getStdOfPlayerTimesPlaying(schedule [][]Match, players []Player, numberRounds float64) float64 {
	var weightedTimesPlayer []float64
	for i, p := range players {
		val := float64(getMatchesCountOfPlayer(schedule, i)) / p.Weight
		weightedTimesPlayer = append(weightedTimesPlayer, val)
	}
	val, err := stats.StandardDeviation(weightedTimesPlayer)
	if err != nil {
		return math.MaxFloat64
	}
	return val * numberRounds
}

func getStdOfPossibleMatches(schedule [][]Match, players []Player, numberRounds float64) float64 {
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
	val, err := stats.StandardDeviation(weightedPossibleMatches)
	if err != nil {
		return math.MaxFloat64
	}
	return val * numberRounds
}

func getStdOfPauseBetweenPlaying(schedule [][]Match, players []Player) float64 {
	pausesBetweenPlaying := []float64{}
	for i := range players {
		roundsPlaying := getRoundIndizesOfPlayer(schedule, i)
		val := calcStdOfPauses(schedule, roundsPlaying)
		pausesBetweenPlaying = append(pausesBetweenPlaying, val)
	}
	val, err := stats.Sum(pausesBetweenPlaying)
	if err != nil {
		return math.MaxFloat64
	}
	return val
}

func getStdOfPauseBetweenMatches(schedule [][]Match, players []Player) float64 {
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
	val, err := stats.Sum(pausesBetweenMatches)
	if err != nil {
		return math.MaxFloat64
	}
	return val
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
