package internal

import log "github.com/sirupsen/logrus"

// Optimizer is responsible for optimizing the schedule
type Optimizer struct {
	Season *Season
}

// Optimize is the main function of the optimizer
// It will optimize the season in a way that the score is minimal for the schedule
func (o *Optimizer) Optimize() {
	swaps := 1
	for swaps > 0 {
		swaps = 0
		log.Debug("Start a new round of optimization")

		log.Debug("Optimizing by swapping players")
		swaps += o.optimizeBySwappingPlayers()

		log.Debug("Optimizing by swapping matches")
		swaps += o.optimizeBySwappingMatches()

		log.WithFields(log.Fields{
			"score": GetScore(o.Season.Schedule, o.Season.Players),
			"swaps": swaps,
		}).Info("Finished a round of optimization")
	}

	log.Info("Optimization finished")
}

func (o *Optimizer) optimizeBySwappingPlayers() int {
	swaps := 0
	currentScore := GetScore(o.Season.Schedule, o.Season.Players)
	var newScore float64

	log.Debug("Start replacing matches")
	for i, round := range o.Season.Schedule {
		if isInSlice(i, o.Season.fixedRounds) {
			continue
		}
		for j, currentMatch := range round {
			for p := range o.Season.Players {
				for q := p + 1; q < len(o.Season.Players); q++ {
					possibleMatch, err := createMatch(p, q)
					if err != nil {
						continue
					}
					if possibleMatch == currentMatch {
						continue
					}
					changed := o.Season.changeMatch(i, j, possibleMatch)
					if !changed {
						continue
					}
					newScore = GetScore(o.Season.Schedule, o.Season.Players)
					if newScore < currentScore {
						swaps++
						currentScore = newScore
						currentMatch = possibleMatch
						log.WithFields(log.Fields{
							"score":    currentScore,
							"newMatch": currentMatch.String(&o.Season.Players),
							"round":    i,
							"match":    j,
						}).Debug("Replaced a match")
					} else {
						o.Season.changeMatch(i, j, currentMatch)
					}
				}
			}
		}
	}
	log.Debug("finished replacing matches")

	currentScore = GetScore(o.Season.Schedule, o.Season.Players)

	log.Debug("Start switching scheduled player of existing round")
	for i, round := range o.Season.Schedule {
		if isInSlice(i, o.Season.fixedRounds) {
			continue
		}
		playerOfRound := getPlayersOfRound(round)
		for j, p := range playerOfRound {
			for k := j + 1; k < len(playerOfRound); k++ {
				q := playerOfRound[k]
				swapped := o.Season.swapPlayersOfRound(i, p, q)
				if !swapped {
					continue
				}
				newScore = GetScore(o.Season.Schedule, o.Season.Players)
				if newScore < currentScore {
					swaps++
					currentScore = newScore
					log.WithFields(log.Fields{
						"score":   currentScore,
						"round":   i,
						"player1": p,
						"player2": q,
					}).Debug("Swapped players")
					break
				} else {
					o.Season.swapPlayersOfRound(i, q, p)
				}
			}
		}
	}
	log.Debug("Finished switching scheduled player of existing round")

	return swaps
}

func (o *Optimizer) optimizeBySwappingMatches() int {
	swaps := 0
	currentScore := GetScore(o.Season.Schedule, o.Season.Players)
	indexCombinations := getIndexCombinations(o.Season.Schedule)

	log.Debug("Finished swapping existing matches")
	for _, combination := range indexCombinations {
		roundIdx1 := combination[0]
		matchIdx1 := combination[1]
		roundIdx2 := combination[2]
		matchIdx2 := combination[3]
		if isInSlice(roundIdx1, o.Season.fixedRounds) || isInSlice(roundIdx2, o.Season.fixedRounds) || o.Season.Schedule[roundIdx1][matchIdx1] == o.Season.Schedule[roundIdx2][matchIdx2] {
			continue
		}
		swapped := o.Season.swapMatches(roundIdx1, matchIdx1, roundIdx2, matchIdx2)
		if !swapped {
			continue
		}
		newScore := GetScore(o.Season.Schedule, o.Season.Players)
		if newScore < currentScore {
			swaps++
			currentScore = newScore
			log.WithFields(log.Fields{
				"score":  currentScore,
				"round1": roundIdx1,
				"match1": matchIdx1,
				"round2": roundIdx2,
				"match2": matchIdx2,
			}).Debug("swapped matches")
		} else {
			o.Season.swapMatches(roundIdx1, matchIdx1, roundIdx2, matchIdx2)
		}
	}
	log.Debug("Finished swapping existing matches")

	return swaps
}

func getIndexCombinations(schedule [][]Match) [][4]int {
	var indizes [][2]int
	var indexCombinations [][4]int
	for i, round := range schedule {
		for j := range round {
			indizes = append(indizes, [2]int{i, j})
		}
	}
	for i, index1 := range indizes {
		for j := i + 1; j < len(indizes); j++ {
			index2 := indizes[j]
			indexCombinations = append(indexCombinations, [4]int{index1[0], index1[1], index2[0], index2[1]})
		}
	}
	return shuffle(indexCombinations)
}
