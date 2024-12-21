package player

import (
	"reflect"
	"testing"
	"time"
)

func TestPlayerToDict(t *testing.T) {
	player := Player{
		Name:       "John Doe",
		CannotPlay: []time.Time{time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
		Weight:     1,
	}

	dict := player.toDict()

	if dict["name"] != player.Name {
		t.Errorf("expected %v, got %v", player.Name, dict["name"])
	}

	if !reflect.DeepEqual(dict["cannot_play"], player.CannotPlay) {
		t.Errorf("expected %v, got %v", player.CannotPlay, dict["cannot_play"])
	}

	if dict["weight"] != player.Weight {
		t.Errorf("expected %v, got %v", player.Weight, dict["weight"])
	}
}

func TestPlayerFromDict(t *testing.T) {
	now := time.Now()
	dict := map[string]interface{}{
		"name":        "Jane Doe",
		"cannot_play": []time.Time{now},
		"weight":      65.5,
	}

	player := Player{}
	player.fromDict(dict)

	if player.Name != dict["name"] {
		t.Errorf("expected %v, got %v", dict["name"], player.Name)
	}

	if !reflect.DeepEqual(player.CannotPlay, dict["cannot_play"]) {
		t.Errorf("expected %v, got %v", dict["cannot_play"], player.CannotPlay)
	}

	if player.Weight != dict["weight"] {
		t.Errorf("expected %v, got %v", dict["weight"], player.Weight)
	}
}

func TestPlayerString(t *testing.T) {
	player := Player{
		Name: "John Doe",
	}

	expected := "John Doe"
	if player.String() != expected {
		t.Errorf("expected %v, got %v", expected, player.String())
	}
}