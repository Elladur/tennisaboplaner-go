package internal

func getPlayersOfRound(round []Match) []int {
	players := make([]int, 0, 8)
	for _, match := range round {
		for _, p := range match.getPlayers() {
			if !isInSlice(p, players) {
				players = append(players, p)
			}
		}
	}
	return players
}
