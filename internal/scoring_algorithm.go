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

	// since weight should define how often someone is playing, the number of playing is same important than all other factors combined
	return 3*val1 + val2 + val3 + val4
}

func getStdOfPlayerTimesPlaying(schedule [][]Match, players []Player, numberRounds float64) float64 {
	var weightedTimesPlayer []float64
	for i, p := range players {
		val := float64(getMatchesCountOfPlayer(schedule, i)) / p.Weight
		weightedTimesPlayer = append(weightedTimesPlayer, val)
	}
	return getCoefficientOfVariation(weightedTimesPlayer)
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
	return getCoefficientOfVariation(weightedPossibleMatches)
}

func getStdOfPauseBetweenPlaying(schedule [][]Match, players []Player) float64 {
	pausesBetweenPlaying := []float64{}
	for i := range players {
		roundsPlaying := getRoundIndizesOfPlayer(schedule, i)
		val := calcCoefficientOfVariationOfPauses(schedule, roundsPlaying)
		pausesBetweenPlaying = append(pausesBetweenPlaying, val)
	}
	val, err := stats.Mean(pausesBetweenPlaying)
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
			val := calcCoefficientOfVariationOfPauses(schedule, roundsPlaying)
			pausesBetweenMatches = append(pausesBetweenMatches, val)
		}
	}
	val, err := stats.Mean(pausesBetweenMatches)
	if err != nil {
		return math.MaxFloat64
	}
	return val
}

func calcCoefficientOfVariationOfPauses(schedule [][]Match, roundsPlaying []int) float64 {
	if len(roundsPlaying) > 1 {
		pauses := getPausesBetweenRounds(schedule, roundsPlaying)
		return getCoefficientOfVariation(pauses)
	}
	return 1 // values are positive and have a max value, therefore Coefficient of Variation is most of the times smaller than 1
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

func getCoefficientOfVariation(input []float64) float64 {
	std, err := stats.StandardDeviation(input)
	mean, err2 := stats.Mean(input)
	if err != nil || err2 != nil || mean == 0 {
		return math.MaxFloat64
	}
	return std / mean
}
