package internal

func getPlayersOfRound(round []Match) []uint8 {
	var players []uint8
	for _, match := range round {
		for _, p := range match.getPlayers() {
			if !isInSlice(p, players) {
				players = append(players, p)
			}
		}
	}
	return players
}
