package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/Elladur/tennisaboplaner-go/internal"
)

func main() {
	content, err := os.ReadFile("test/integration/input/performance.json")
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

	f, _ := os.Create("cpu.prof")
	defer func() {
		if err := f.Close(); err != nil {
			// Handle error if needed
			fmt.Printf("Error closing excel file: %v", err)
		}
	}()
	err = pprof.StartCPUProfile(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pprof.StopCPUProfile()

	err = internal.ExecutePlanerParallel(settings, "output", 100)
	if err != nil {
		fmt.Println(err)
		return
	}
}
