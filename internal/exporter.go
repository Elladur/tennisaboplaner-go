package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/xuri/excelize/v2"
)

const scheduleSheet = "Sheet1"
const scheduleSheetName = "Schedule"
const opponentSheet = "Partner by Player"
const numberOfMatchesSheet = "NumMatches"
const costSheet = "Cost"

// Export the season into various files (excel & ical)
func (s Season) Export(directory string) error {
	err := s.exportExcel(directory)
	if err != nil {
		return err
	}
	err = s.exportCalendarFiles(directory)
	if err != nil {
		return err
	}
	return nil
}

func (s Season) exportExcel(directoy string) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			// Handle error if needed
			fmt.Printf("Error closing excel file: %v", err)
		}
	}()

	// initial sheet for schedule
	if err := addInitialScheduleSheet(f, s); err != nil {
		return err
	}

	// new sheet for oppopnent tab
	if err := addOpponentSheet(f, s); err != nil {
		return err
	}

	// new sheet for matches tab
	if err := addNumberOfMatchesSheet(f, s); err != nil {
		return err
	}

	// new sheet for costs
	if err := addCostSheet(f, s); err != nil {
		return err
	}

	if err := f.SaveAs(filepath.Join(directoy, "schedule.xlsx")); err != nil {
		return err
	}
	return nil
}


func addInitialScheduleSheet(f *excelize.File, s Season) error {
	f.WorkBook.Sheets.Sheet[f.GetActiveSheetIndex()].Name = scheduleSheetName
	// add header
	err := f.SetCellValue(scheduleSheet, "A1", "Date")
	if err != nil {
		return err
	}
	for i := 0; i < s.NumberOfCourts; i++ {
		err = f.SetCellValue(scheduleSheet, getCell(0, i+1), fmt.Sprintf("Match %d", i+1))
		if err != nil {
			return err
		}
	}
	for i, round := range s.Schedule {
		row := i + 1
		err = f.SetCellValue(scheduleSheet, getCell(row, 0), s.dates[i])
		if err != nil {
			return err
		}
		for j, m := range round {
			err = f.SetCellValue(scheduleSheet, getCell(row, j+1), m.String(&s.Players))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func addOpponentSheet(f *excelize.File, s Season) error {
	_, err := f.NewSheet(opponentSheet)
	if err != nil {
		return err
	}
	if err := f.SetCellValue(opponentSheet, getCell(0, 0), "Date"); err != nil {
		return err
	}
	for i, p := range s.Players {
		if err := f.SetCellValue(opponentSheet, getCell(0, i+1), p.Name); err != nil {
			return err
		}
	}
	for i, round := range s.Schedule {
		row := i + 1
		if err := f.SetCellValue(opponentSheet, getCell(row, 0), s.dates[i]); err != nil {
			return err
		}
		for j := range s.Players {
			for _, m := range round {
				var opponent string
				switch {
				case m.player1 == j && m.isPlayer2Set:
					opponent = s.Players[m.player2].Name
				case m.player1 == j:
					opponent = "..."
				case m.isPlayer2Set && m.player2 == j:
					opponent = s.Players[m.player1].Name
				default:
					continue
				}
				if err := f.SetCellValue(opponentSheet, getCell(row, j+1), opponent); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func addNumberOfMatchesSheet(f *excelize.File, s Season) error {
	_, err := f.NewSheet(numberOfMatchesSheet)
	if err != nil {
		return err
	}
	if err := f.SetCellValue(numberOfMatchesSheet, getCell(0, 0), "Players"); err != nil {
		return err
	}
	for i, p := range s.Players {
		if err := f.SetCellValue(numberOfMatchesSheet, getCell(0, i+1), p.Name); err != nil {
			return err
		}
		for j := range s.Players {
			if i == j {
				continue
			}
			match, err := createMatch(i, j)
			if err != nil {
				return err
			}
			getnumMatches := float64(len(getRoundIndizesOfMatch(s.Schedule, match)))
			if err := f.SetCellValue(numberOfMatchesSheet, getCell(j+1, i+1), getnumMatches); err != nil {
				return err
			}
		}
	}
	for i, p := range s.Players {
		if err := f.SetCellValue(numberOfMatchesSheet, getCell(i+1, 0), p.Name); err != nil {
			return err
		}
	}
	return nil
}

func addCostSheet(f *excelize.File, s Season) error {
	costPerMatch := s.OverallCosts / float64(len(s.Schedule)*s.NumberOfCourts) / 2
	_, err := f.NewSheet(costSheet)
	if err != nil {
		return err
	}
	if err := f.SetCellValue(costSheet, getCell(0, 0), "Player"); err != nil {
		return err
	}
	if err := f.SetCellValue(costSheet, getCell(1, 0), "Match"); err != nil {
		return err
	}
	if err := f.SetCellValue(costSheet, getCell(2, 0), "Cost"); err != nil {
		return err
	}
	for i, p := range s.Players {
		timesPlaying := float64(len(getRoundIndizesOfPlayer(s.Schedule, i)))
		if err := f.SetCellValue(costSheet, getCell(0, i+1), p.Name); err != nil {
			return err
		}
		if err := f.SetCellValue(costSheet, getCell(1, i+1), timesPlaying); err != nil {
			return err
		}
		if err := f.SetCellValue(costSheet, getCell(2, i+1), costPerMatch*timesPlaying); err != nil {
			return err
		}
	}
	return nil
}

// calculates excel column of a zero-based-index
func getColName(index int) string {
	return string(rune('A' + index))
}

// calculates cell of zero-based indizes
func getCell(i, j int) string {
	col := getColName(j)
	return fmt.Sprintf("%s%d", col, i+1)
}

func (s Season) exportCalendarFiles(directory string) error {
	for i := range s.Players {
		if err := s.exportCalendarFileForPlayer(i, directory); err != nil {
			return err
		}
	}
	return nil
}

func (s Season) exportCalendarFileForPlayer(player int, directory string) error {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	for _, i := range getRoundIndizesOfPlayer(s.Schedule, player) {
		date := s.dates[i]
		event := cal.AddEvent(fmt.Sprintf("Tennisabo %s", date.Format(time.DateOnly)))
		event.SetCreatedTime(time.Now())
		event.SetDtStampTime(time.Now())
		event.SetModifiedAt(time.Now())
		startDate := time.Date(date.Year(), date.Month(), date.Day(), s.StartTime.Hour, s.StartTime.Minute, 0, 0, &s.location)
		endDate := time.Date(date.Year(), date.Month(), date.Day(), s.EndTime.Hour, s.EndTime.Minute, 0, 0, &s.location)
		event.SetStartAt(startDate)
		event.SetEndAt(endDate)
		event.SetSummary("Tennisabo")
		event.SetDescription(convertRoundToString(s.Schedule[i], &s.Players))
	}
	fileContent := cal.Serialize()
	f, err := os.Create(filepath.Join(directory, fmt.Sprintf("%s.ics", s.Players[player].Name)))
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			// Handle error if needed
			fmt.Printf("Error closing excel file: %v", err)
		}
	}()

	_, err = f.WriteString(fileContent)
	if err != nil {
		return err
	}
	return nil
}
