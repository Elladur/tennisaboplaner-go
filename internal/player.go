package internal

import (
	"encoding/json"
	"time"
)

type Player struct {
	Name       string
	CannotPlay []time.Time
	Weight     float64
}

func (p Player) ToDict() map[string]interface{} {
	return map[string]interface{}{
		"name":        p.Name,
		"cannot_play": p.CannotPlay,
		"weight":      p.Weight,
	}
}

func (p *Player) FromDict(d map[string]interface{}) {
	p.Name = d["name"].(string)
	p.CannotPlay = d["cannot_play"].([]time.Time)
	p.Weight = d["weight"].(float64)
}

func (p Player) String() string {
	return p.Name
}

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
