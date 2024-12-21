package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPlayerToDict(t *testing.T) {
	player := Player{
		Name:       "John Doe",
		CannotPlay: []time.Time{time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
		Weight:     1,
	}

	dict := player.toDict()

	assert.Equal(t, player.Name, dict["name"])
	assert.Equal(t, player.CannotPlay, dict["cannot_play"])
	assert.Equal(t, player.Weight, dict["weight"])
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

	assert.Equal(t, dict["name"], player.Name)
	assert.Equal(t, dict["cannot_play"], player.CannotPlay)
	assert.Equal(t, dict["weight"], player.Weight)
}

func TestPlayerString(t *testing.T) {
	expected := "John Doe"
	player := Player{
		Name: expected,
	}

	assert.Equal(t, expected, string(player.String()))
}
