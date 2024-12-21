package internal

func getPlayersOfRound(round []Match) []uint8 {
	var players []uint8
	for _, match := range round {
		players = append(players, match.GetPlayers()...)
	}
	return players
}