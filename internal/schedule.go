package internal

func getMatchIndizesOfPlayer(schedule [][]Match, player uint8) [][2]int {
	var indizes [][2]int
	for i, round := range schedule {
		for j, match := range round {
			if match.player1 == player || (match.isPlayer2Set && match.player2 == player) {
				indizes = append(indizes, [2]int{i, j})
			}
		}
	}
	return indizes
}

func getMatchIndizesOfMatch(schedule [][]Match, match Match) [][2]int {
	var indizes [][2]int
	for i, round := range schedule {
		for j, m := range round {
			if m == match {
				indizes = append(indizes, [2]int{i, j})
			}
		}
	}
	return indizes
}

func getRoundIndizesOfMatch(schedule [][]Match, match Match) []int {
	indizes := make([]int, 0, len(schedule))
	for i, round := range schedule {
		for _, m := range round {
			if m == match {
				indizes = append(indizes, i)
			}
		}
	}
	return indizes
}

func getCountOfMatchInSchedule(schedule [][]Match, match Match) int {
	num := 0
	for _, round := range schedule {
		for _, m := range round {
			if m == match {
				num++
			}
		}
	}
	return num
}

func convertMatchIndizesToRoundIndizes(indizes [][2]int) []int {
	var roundIndizes []int
	for _, ind := range indizes {
		roundIndizes = append(roundIndizes, ind[0])
	}
	return roundIndizes
}
