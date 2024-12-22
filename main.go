package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Elladur/tennisaboplaner-go/internal"
)

func main() {
	fmt.Println("Hello, Go!")
	players := []internal.Player{
		{Name: "Alice"},
		{Name: "Bob"},
	}
	jsonPlayers, _ := json.Marshal(players)
	fmt.Println(players)
	fmt.Println(string(jsonPlayers))

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
	fmt.Println(settings)
	season, err := internal.CreateSeasonFromSettings(settings)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(season)
}
