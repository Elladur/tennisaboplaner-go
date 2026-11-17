package internal

// Schedule is a season's schedule: a slice of rounds, each round a slice of matches.
type Schedule [][]Match

func (s Schedule) Rounds() int { return len(s) }

func (s Schedule) PlayersInRound(i int) []int {
	if i < 0 || i >= len(s) {
		return nil
	}
	return getPlayersOfRound(s[i])
}

func (s Schedule) RoundIndizesOfPlayer(player int) []int {
	indizes := make([]int, 0, len(s))
	for i, round := range s {
		for _, match := range round {
			if match.player1 == player || (match.isPlayer2Set && match.player2 == player) {
				indizes = append(indizes, i)
			}
		}
	}
	return indizes
}

func (s Schedule) MatchesCountOfPlayer(player int) int {
	num := 0
	for _, round := range s {
		for _, match := range round {
			if match.player1 == player || (match.isPlayer2Set && match.player2 == player) {
				num++
			}
		}
	}
	return num
}

func (s Schedule) RoundIndizesOfMatch(match Match) []int {
	indizes := make([]int, 0, len(s))
	for i, round := range s {
		for _, m := range round {
			if m == match {
				indizes = append(indizes, i)
			}
		}
	}
	return indizes
}

func (s Schedule) CountOfMatch(match Match) int {
	num := 0
	for _, round := range s {
		for _, m := range round {
			if m == match {
				num++
			}
		}
	}
	return num
}
