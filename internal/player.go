package internal

import (
	"encoding/json"
	"time"
)

// Player represents a player with a name, a list of dates when the player cannot play and a weight to determine if he wants to play full abo or less
type Player struct {
	Name       string
	CannotPlay []time.Time
	Weight     float64
}

func (p Player) String() string {
	return p.Name
}

// MarshalJSON marshals a Player to JSON
func (p Player) MarshalJSON() ([]byte, error) {
	var cannotPlay []string
	for _, t := range p.CannotPlay {
		cannotPlay = append(cannotPlay, t.Format(time.DateOnly))
	}
	return json.Marshal(map[string]interface{}{
		"Name":       p.Name,
		"CannotPlay": cannotPlay,
		"Weight":     p.Weight,
	})
}

// UnmarshalJSON unmarshals a Player from JSON
func (p *Player) UnmarshalJSON(data []byte) error {
	var player struct {
		Name       string
		CannotPlay []string
		Weight     float64
	}
	err := json.Unmarshal(data, &player)
	if err != nil {
		return err
	}
	p.Name = player.Name
	for _, s := range player.CannotPlay {
		t, err := time.Parse(time.DateOnly, s)
		if err != nil {
			return err
		}
		p.CannotPlay = append(p.CannotPlay, t)
	}
	p.Weight = player.Weight
	return nil
}
