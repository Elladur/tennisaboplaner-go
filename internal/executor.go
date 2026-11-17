// Package internal contains domain logic for the tennis abo planner.
package internal

import (
	"context"
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
func ExecutePlanerSerial(ctx context.Context, settings SeasonSettings, directory string, executions int) error {
	if executions <= 0 {
		return fmt.Errorf("executions must be positive")
	}
	bestSeason := seasonResult{season: Season{}, score: math.MaxFloat64}
	hasResult := false

	log.WithFields(log.Fields{
		"times": executions,
	}).Info("Start planning multiple times")
	for i := 0; i < executions; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		s, err := CreateSeasonFromSettings(settings)
		if err != nil {
			return err
		}
		optimizer := Optimizer{Season: &s}
		score := optimizer.Optimize(ctx)
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
func ExecutePlanerParallel(ctx context.Context, settings SeasonSettings, directory string, executions int) error {
	if executions <= 0 {
		return fmt.Errorf("executions must be positive")
	}
	bestSeason := seasonResult{season: Season{}, score: math.MaxFloat64}
	hasResult := false
	c := make(chan seasonResult, executions)

	log.WithFields(log.Fields{
		"times": executions,
	}).Info("Start planning multiple times")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for range executions {
		go executeOptimize(ctx, settings, c)
	}
	for n := 1; n <= executions; n++ {
		select {
		case result := <-c:
			if result.err != nil {
				cancel()
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
		case <-ctx.Done():
			cancel()
			return ctx.Err()
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

func executeOptimize(ctx context.Context, settings SeasonSettings, c chan seasonResult) {
	s, err := CreateSeasonFromSettings(settings)
	if err != nil {
		c <- seasonResult{err: err}
		return
	}
	optimizer := Optimizer{Season: &s}
	score := optimizer.Optimize(ctx)
	select {
	case c <- seasonResult{season: s, score: score}:
	case <-ctx.Done():
	}
}
