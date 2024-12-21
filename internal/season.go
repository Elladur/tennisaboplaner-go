package internal

import "time"

type Season struct {
	Players []Player
	Start time.Time
	End time.Time
	StartTime simpleTime
	EndTime simpleTime
	NumberOfCourts int
	CalendarTitle string
	OverallCosts float64
	ExcludedDates []time.Time
	dates []time.Time
	fixedRounds []int
	Schedule [][]Match
}

type simpleTime struct {
	Hour uint8
	Minute uint8
}

func CreateSeason(players []Player, start time.Time, end time.Time, numberOfCourts int, calendarTitle string, overallCosts float64, excludedDates []time.Time) Season {
	startTime := simpleTime{uint8(start.Hour()), uint8(start.Minute())}
	endTime := simpleTime{uint8(end.Hour()), uint8(end.Minute())}
	start = start.Truncate(24 * time.Hour)
	end = end.Truncate(24 * time.Hour)
	var dates []time.Time
	for d := start; d.Before(end); d = d.AddDate(0, 0, 7) {
		if !isInSlice(d, excludedDates) {
			dates = append(dates, d)
		}
	}

	return Season{
		Players: players,
		Start: start,
		End: end,
		StartTime: startTime,
		EndTime: endTime,
		NumberOfCourts: numberOfCourts,
		CalendarTitle: calendarTitle,
		OverallCosts: overallCosts,
		ExcludedDates: excludedDates,
		dates: dates,
	}
}