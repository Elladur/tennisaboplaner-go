package internal

func getRoundIndizesOfPlayer(schedule [][]Match, player int) []int {
	indizes := make([]int, 0, len(schedule))
	for i, round := range schedule {
		for _, match := range round {
			if match.player1 == player || (match.isPlayer2Set && match.player2 == player) {
				indizes = append(indizes, i)
			}
		}
	}
	return indizes
}

func getMatchesCountOfPlayer(schedule [][]Match, player int) int {
	num := 0
	for _, round := range schedule {
		for _, match := range round {
			if match.player1 == player || (match.isPlayer2Set && match.player2 == player) {
				num++
			}
		}
	}
	return num
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
