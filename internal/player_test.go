package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPlayerString(t *testing.T) {
	expected := "John Doe"
	player := Player{
		Name: expected,
	}

	assert.Equal(t, expected, string(player.String()))
}

func TestPlayerMarshalJSON(t *testing.T) {
	player := Player{
		Name:       "John Doe",
		CannotPlay: []time.Time{time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
		Weight:     75.5,
	}

	data, err := player.MarshalJSON()
	assert.NoError(t, err)

	expected := `{"Name":"John Doe","CannotPlay":["2020-01-01","2021-01-01"],"Weight":75.5}`
	assert.JSONEq(t, expected, string(data))
}

func TestPlayerUnmarshalJSON(t *testing.T) {
	data := []byte(`{"Name":"Jane Doe","CannotPlay":["2020-01-01","2021-01-01"],"Weight":65.5}`)

	var player Player
	err := player.UnmarshalJSON(data)
	assert.NoError(t, err)

	expected := Player{
		Name:       "Jane Doe",
		CannotPlay: []time.Time{time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
		Weight:     65.5,
	}

	assert.Equal(t, expected.Name, player.Name)
	assert.Equal(t, expected.CannotPlay, player.CannotPlay)
	assert.Equal(t, expected.Weight, player.Weight)
}
