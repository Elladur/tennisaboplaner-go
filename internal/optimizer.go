package internal

import log "github.com/sirupsen/logrus"

// Optimizer is responsible for optimizing the schedule
type Optimizer struct {
	Season *Season
}

// applyChangeMatch mutates schedule and stats together; returns false if rejected.
func applyChangeMatch(stats *scheduleStats, s *Season, roundIdx, matchIdx int, newMatch Match) bool {
	old := s.Schedule[roundIdx][matchIdx]
	if !s.changeMatch(roundIdx, matchIdx, newMatch) {
		return false
	}
	stats.removeMatch(old, roundIdx)
	stats.addMatch(newMatch, roundIdx)
	return true
}

// applySwapMatches mutates schedule and stats together; returns false if rejected.
func applySwapMatches(stats *scheduleStats, s *Season, r1, m1, r2, m2 int) bool {
	old1 := s.Schedule[r1][m1]
	old2 := s.Schedule[r2][m2]
	if !s.swapMatches(r1, m1, r2, m2) {
		return false
	}
	stats.removeMatch(old1, r1)
	stats.removeMatch(old2, r2)
	stats.addMatch(old2, r1)
	stats.addMatch(old1, r2)
	return true
}

// applySwapPlayersOfRound swaps players within a round and keeps stats in sync.
// Returns the pre-swap round snapshot and ok.
func applySwapPlayersOfRound(stats *scheduleStats, s *Season, roundIdx, p, q int) ([]Match, bool) {
	oldRound := make([]Match, len(s.Schedule[roundIdx]))
	copy(oldRound, s.Schedule[roundIdx])
	if !s.swapPlayersOfRound(roundIdx, p, q) {
		return nil, false
	}
	for _, m := range oldRound {
		stats.removeMatch(m, roundIdx)
	}
	for _, m := range s.Schedule[roundIdx] {
		stats.addMatch(m, roundIdx)
	}
	return oldRound, true
}

// Optimize is the main function of the optimizer
// It will optimize the season in a way that the score is minimal for the schedule
func (o *Optimizer) Optimize() float64 {
	swaps := 1
	var score float64
	for swaps > 0 {
		swaps = 0
		log.Debug("Start a new round of optimization")

		log.Debug("Optimizing by swapping players")
		swaps += o.optimizeBySwappingPlayers()

		log.Debug("Optimizing by swapping matches")
		swaps += o.optimizeBySwappingMatches()

		score = GetScore(o.Season.Schedule, o.Season.Players)
		log.WithFields(log.Fields{
			"score": score,
			"swaps": swaps,
		}).Debug("Finished a round of optimization")
	}

	log.Debug("Optimization finished")
	return score
}

func (o *Optimizer) optimizeBySwappingPlayers() int {
	swaps := 0
	weights := defaultScoreWeights()
	stats := computeScheduleStats(o.Season.Schedule, o.Season.Players)
	currentScore := getScoreFromStats(stats, o.Season.Schedule, o.Season.Players, weights)
	var newScore float64

	log.Debug("Start replacing matches")
	for i, round := range o.Season.Schedule {
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
					changed := applyChangeMatch(&stats, o.Season, i, j, possibleMatch)
					if !changed {
						continue
					}
					newScore = getScoreFromStats(stats, o.Season.Schedule, o.Season.Players, weights)
					if newScore < currentScore {
						swaps++
						currentScore = newScore
						currentMatch = possibleMatch
						log.WithFields(log.Fields{
							"score":    currentScore,
							"newMatch": currentMatch.String(o.Season.Players),
							"round":    i,
							"match":    j,
						}).Debug("Replaced a match")
					} else {
						applyChangeMatch(&stats, o.Season, i, j, currentMatch)
					}
				}
			}
		}
	}
	log.Debug("finished replacing matches")

	// swapPlayersOfRound is a structural change; recompute for correctness.
	stats = computeScheduleStats(o.Season.Schedule, o.Season.Players)
	currentScore = getScoreFromStats(stats, o.Season.Schedule, o.Season.Players, weights)

	log.Debug("Start switching scheduled player of existing round")
	for i, round := range o.Season.Schedule {
		playerOfRound := getPlayersOfRound(round)
		for j, p := range playerOfRound {
			for k := j + 1; k < len(playerOfRound); k++ {
				q := playerOfRound[k]
				oldRound, swapped := applySwapPlayersOfRound(&stats, o.Season, i, p, q)
				if !swapped {
					continue
				}
				oldScore := currentScore
				newScore = getScoreFromStats(stats, o.Season.Schedule, o.Season.Players, weights)
				if newScore < currentScore {
					swaps++
					currentScore = newScore
					log.WithFields(log.Fields{
						"score":    currentScore,
						"oldScore": oldScore,
						"round":    i,
						"player1":  p,
						"player2":  q,
					}).Debug("Swapped players")
					break
				} else {
					// restore stats to reflect the original schedule
					for _, m := range o.Season.Schedule[i] {
						stats.removeMatch(m, i)
					}
					for _, m := range oldRound {
						stats.addMatch(m, i)
					}
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
	weights := defaultScoreWeights()
	stats := computeScheduleStats(o.Season.Schedule, o.Season.Players)
	currentScore := getScoreFromStats(stats, o.Season.Schedule, o.Season.Players, weights)
	indexCombinations := getIndexCombinations(o.Season.Schedule)

	log.Debug("Start swapping existing matches")
	for _, combination := range indexCombinations {
		roundIdx1 := combination[0]
		matchIdx1 := combination[1]
		roundIdx2 := combination[2]
		matchIdx2 := combination[3]
		if o.Season.Schedule[roundIdx1][matchIdx1] == o.Season.Schedule[roundIdx2][matchIdx2] {
			continue
		}
		swapped := applySwapMatches(&stats, o.Season, roundIdx1, matchIdx1, roundIdx2, matchIdx2)
		if !swapped {
			continue
		}
		newScore := getScoreFromStats(stats, o.Season.Schedule, o.Season.Players, weights)
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
			applySwapMatches(&stats, o.Season, roundIdx1, matchIdx1, roundIdx2, matchIdx2)
		}
	}
	log.Debug("Finished swapping existing matches")

	return swaps
}

func getIndexCombinations(schedule Schedule) [][4]int {
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
