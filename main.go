package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/Elladur/tennisaboplaner-go/internal"
)

func main() {

	content, err := os.ReadFile("settings.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	var settings internal.SeasonSettings
	err = json.Unmarshal(content, &settings)
	if err != nil {
		fmt.Println(err)
		return
	}
	season, err := internal.CreateSeasonFromSettings(settings)
	if err != nil {
		fmt.Println(err)
		return
	}
	optimizer := internal.Optimizer{Season: &season}

	f, _ := os.Create("cpu.prof")
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	optimizer.Optimize()

	fmt.Printf("Optimized Schedule and new Score is %.2f\n", internal.GetScore(season.Schedule, season.Players))
}
