package main

import (
	"fmt"
	"github.com/dank/go-csgsi"
	"github.com/hugolgst/rich-go/client"
	"time"
)

type MatchDetails struct {
	tScore    int
	ctScore   int
	mapName   string
	timestamp client.Timestamps
}

type Connection struct {
	state    *csgsi.State
	activity client.Activity
}

var lastMatch MatchDetails

func initializeRPC() {
	err := client.Login("937726683442712657")
	if err != nil {
		panic(err)
	}
}

func setState(state *csgsi.State) {
	c := Connection{
		state:    state,
		activity: client.Activity{},
	}

	if state.Map != nil {
		c.setGameState()
	} else {
		// Ended a game
		if lastMatch.mapName != "Menu" {
			now := time.Now()
			lastMatch = MatchDetails{
				tScore:  0,
				ctScore: 0,
				mapName: "Menu",
				timestamp: client.Timestamps{
					Start: &now,
				},
			}
		}

		err := client.SetActivity(client.Activity{
			Details:    "On Menu",
			LargeImage: "csgo",
			LargeText:  "Counter-Strike: Global Offensive",
			Timestamps: &lastMatch.timestamp,
		})

		if err != nil {
			panic(err)
		}
	}
}

func (c *Connection) setGameState() {
	c.setMapIcon()
	c.checkIfIsSameGame()
	c.setScoreboard()

	err := client.SetActivity(c.activity)

	if err != nil {
		panic(err)
	}
}

func (c *Connection) checkIfIsSameGame() {
	if lastMatch.mapName != c.state.Map.Name || lastMatch.tScore > c.state.Map.Team_t.Score || lastMatch.ctScore > c.state.Map.Team_ct.Score {
		now := time.Now()
		lastMatch.timestamp = client.Timestamps{
			Start: &now,
		}
	}

	c.activity.Timestamps = &lastMatch.timestamp
	lastMatch.mapName = c.state.Map.Name
	lastMatch.tScore = c.state.Map.Team_t.Score
	lastMatch.ctScore = c.state.Map.Team_ct.Score
}

func (c *Connection) setMapIcon() {
	mapsWithIcon := []string{"cs_agency", "cs_assault", "cs_insertion2", "cs_italy", "cs_militia", "cs_office", "de_ancient", "de_basalt", "de_cache", "de_canals", "de_cbble", "de_dust2", "de_inferno", "de_mirage", "de_nuke", "de_overpass", "de_train", "de_vertigo"}
	currentMap := "csgo"
	for _, mapName := range mapsWithIcon {
		if c.state.Map.Name == mapName {
			currentMap = c.state.Map.Name
			break
		}
	}

	c.activity.Details = c.state.Map.Name
	// Default CSGO icon if map has no icon
	c.activity.LargeImage = currentMap
	c.activity.LargeText = c.state.Map.Name
}

func (c *Connection) setScoreboard() {
	switch c.state.Map.Phase {
	case "live":
		c.activity.State += "Playing "
	case "warmup":
		c.activity.State += "Warming up "
	case "intermission":
		c.activity.State += "Pause "
	case "gameover":
		c.activity.State += "Ending "
	}

	if c.state.Player.Team == "CT" {
		c.activity.State += fmt.Sprintf("[%d : %d]", c.state.Map.Team_ct.Score, c.state.Map.Team_t.Score)
		c.activity.SmallImage = "ct"
		c.activity.SmallText = "Counter-Terrorist"
	} else if c.state.Player.Team == "T" {
		c.activity.State += fmt.Sprintf("[%d : %d]", c.state.Map.Team_t.Score, c.state.Map.Team_ct.Score)
		c.activity.SmallImage = "t"
		c.activity.SmallText = "Terrorist"
	} else {
		c.activity.State += fmt.Sprintf("[%d : %d]", c.state.Map.Team_ct.Score, c.state.Map.Team_t.Score)
		c.activity.SmallImage = "spectator"
		c.activity.SmallText = "Spectator"
	}
}