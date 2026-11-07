package internal

import (
	"fmt"
	"math"

	log "github.com/sirupsen/logrus"
)

type seasonResult struct {
	season Season
	score  float64
	err    error
}

// ExecutePlanerSerial does n times planning and exports the best result
func ExecutePlanerSerial(settings SeasonSettings, directory string, executions int) error {
	if executions <= 0 {
		return fmt.Errorf("executions must be positive")
	}
	bestSeason := seasonResult{season: Season{}, score: math.MaxFloat64}
	hasResult := false

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
			hasResult = true
		}
	}
	if !hasResult {
		return fmt.Errorf("no season produced")
	}
	log.WithFields(log.Fields{
		"score": bestSeason.score,
	}).Info("Finished Planning")

	log.Info("Start exporting")
	err := bestSeason.season.Export(directory)
	if err != nil {
		return err
	}
	log.Info("Finished.")
	return nil
}

// ExecutePlanerParallel does n times planning using go routines and exports the best result
func ExecutePlanerParallel(settings SeasonSettings, directory string, executions int) error {
	if executions <= 0 {
		return fmt.Errorf("executions must be positive")
	}
	bestSeason := seasonResult{season: Season{}, score: math.MaxFloat64}
	hasResult := false
	c := make(chan seasonResult)

	log.WithFields(log.Fields{
		"times": executions,
	}).Info("Start planning multiple times")

	for range executions {
		go executeOptimize(settings, c)
	}
	for n := 1; n <= executions; n++ {
		result := <-c
		if result.err != nil {
			return result.err
		}
		if result.score < bestSeason.score {
			bestSeason = result
			hasResult = true
		}
		step := executions / 10
		if step < 1 {
			step = 1
		}
		if n == 1 || n == executions || n%step == 0 {
			log.WithFields(log.Fields{
				"completed": n,
				"total":     executions,
				"percent":   n * 100 / executions,
			}).Info("Planning progress")
		}
	}
	if !hasResult {
		return fmt.Errorf("no season produced")
	}
	log.WithFields(log.Fields{
		"score": bestSeason.score,
	}).Info("Finished Planning")

	log.Info("Start exporting")
	err := bestSeason.season.Export(directory)
	if err != nil {
		return err
	}
	log.Info("Finished.")
	return nil
}

func executeOptimize(settings SeasonSettings, c chan seasonResult) {
	s, err := CreateSeasonFromSettings(settings)
	if err != nil {
		c <- seasonResult{err: err}
		return
	}
	optimizer := Optimizer{Season: &s}
	score := optimizer.Optimize()
	c <- seasonResult{season: s, score: score}
}
