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

func convertMatchIndizesToRoundIndizes(indizes [][2]int) []int {
	var roundIndizes []int
	for _, ind := range indizes {
		roundIndizes = append(roundIndizes, ind[0])
	}
	return roundIndizes
}
