package internal

func GetMatchIndizesOfPlayer(schedule [][]Match, player uint8) [][]int {
	var indizes [][]int
	for i, round := range schedule {
		for j, match := range round {
			if match.Player1 == player || (match.IsPlayer2Set && match.Player2 == player) {
				indizes = append(indizes, []int{i, j})
			}
		}
	}
	return indizes
}

func GetMatchIndizesOfMatch(schedule [][]Match, match Match) [][]int {
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