package internal

func getMatchIndizesOfPlayer(schedule [][]Match, player uint8) [][]int {
	var indizes [][]int
	for i, round := range schedule {
		for j, match := range round {
			if match.player1 == player || (match.isPlayer2Set && match.player2 == player) {
				indizes = append(indizes, []int{i, j})
			}
		}
	}
	return indizes
}

func getMatchIndizesOfMatch(schedule [][]Match, match Match) [][]int {
	var indizes [][]int
	for i, round := range schedule {
		for j, m := range round {
			if m == match {
				indizes = append(indizes, []int{i, j})
			}
		}
	}
	return indizes
}