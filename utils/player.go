package player

import "time"

type Player struct{
	Name string
	CannotPlay []time.Time
	Weight float64
}


func (p *Player) toDict() map[string]interface{} {
	return map[string]interface{}{
		"name": p.Name,
		"cannot_play": p.CannotPlay,
		"weight": p.Weight,
	}
}

func (p *Player) fromDict(d map[string]interface{}) {
	p.Name = d["name"].(string)
	p.CannotPlay = d["cannot_play"].([]time.Time)
	p.Weight = d["weight"].(float64)
}

func (p *Player) String() string {
	return (p.Name)
}