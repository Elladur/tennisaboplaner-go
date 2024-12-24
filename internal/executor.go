package internal

import (
	"math"

	log "github.com/sirupsen/logrus"
)

// ExecutePlanerSerial does n times planning and exports the best result
func ExecutePlanerSerial(settings SeasonSettings, directory string, executions int) error {
	var bestSeason, s Season
	var err error
	bestSeasonScore := math.MaxFloat64

	log.WithFields(log.Fields{
		"times": executions,
	}).Info("Start planning multiple times")
	for i := 0; i < executions; i++ {
		s, err = CreateSeasonFromSettings(settings)
		if err != nil {
			return err
		}
		optimizer := Optimizer{Season: &s}
		score := optimizer.Optimize()
		if score < bestSeasonScore {
			bestSeasonScore = score
			bestSeason = s
		}
	}
	log.WithFields(log.Fields{
		"score": bestSeasonScore,
	}).Info("Finished Planning")

	log.Info("Start exporting")
	bestSeason.Export(directory)
	log.Info("Finished.")
	return nil
}
