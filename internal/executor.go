package internal

import (
	"math"

	log "github.com/sirupsen/logrus"
)

type seasonResult struct {
	season Season
	score  float64
}

// ExecutePlanerSerial does n times planning and exports the best result
func ExecutePlanerSerial(settings SeasonSettings, directory string, executions int) error {
	bestSeason := seasonResult{season: Season{}, score: math.MaxFloat64}

	log.WithFields(log.Fields{
		"times": executions,
	}).Info("Start planning multiple times")
	for i := 0; i < executions; i++ {
		s, err := CreateSeasonFromSettings(settings)
		if err != nil {
			return err
		}
		optimizer := Optimizer{Season: &s}
		score := optimizer.Optimize()
		if score < bestSeason.score {
			bestSeason = seasonResult{season: s, score: score}
		}
	}
	log.WithFields(log.Fields{
		"score": bestSeason.score,
	}).Info("Finished Planning")

	log.Info("Start exporting")
	bestSeason.season.Export(directory)
	log.Info("Finished.")
	return nil
}

// ExecutePlanerParallel does n times planning using go routines and exports the best result
func ExecutePlanerParallel(settings SeasonSettings, directory string, executions int) error {
	bestSeason := seasonResult{season: Season{}, score: math.MaxFloat64}
	c := make(chan seasonResult)

	log.WithFields(log.Fields{
		"times": executions,
	}).Info("Start planning multiple times")
	for i := 0; i < executions; i++ {
		go executeOptimize(settings, c)
	}
	for i := 0; i < executions; i++ {
		result := <-c
		if result.score < bestSeason.score {
			bestSeason = result
		}
	}
	log.WithFields(log.Fields{
		"score": bestSeason.score,
	}).Info("Finished Planning")

	log.Info("Start exporting")
	bestSeason.season.Export(directory)
	log.Info("Finished.")
	return nil
}

func executeOptimize(settings SeasonSettings, c chan seasonResult) {
	s, err := CreateSeasonFromSettings(settings)
	if err != nil {
		log.Fatal("Could not create season from settings: ", err)
		return
	}
	optimizer := Optimizer{Season: &s}
	score := optimizer.Optimize()
	c <- seasonResult{s, score}
}
