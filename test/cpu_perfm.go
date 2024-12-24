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
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	internal.ExecutePlanerParallel(settings, "output", 100)
}
