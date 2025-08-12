package main

import (
	"testing"
)

func TestReadLeagues(t *testing.T) {
	config := Config{
		leaguesFile: "../../leagues.json",
	}

	leagues, err := readLeagues(config)
	t.Log(leagues)
	if err != nil {
		t.Errorf("failed to parse leagues: %v", err)
	}

}
