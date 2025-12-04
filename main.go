package main

import (
	"fmt"
	"github.com/kregan77/dartbuddy/internal/model"
)

func main() {
	// The main function is intentionally left empty.
	g := model.New01Game(401)
	g.AddPlayer(model.NewPlayer("Alice", 32.0, model.TwentiesScoringPreference))
	g.AddPlayer(model.NewPlayer("Anthony", 50.0, model.TwentiesScoringPreference))
	g.Start()
	for _ = range 50 {
		r := g.PlayTurn()
		switch r.Type {
		case model.WinTurn:
			fmt.Printf("Player %s scored %d points for the win(3DA: %.2f)!\n\n",
				r.PlayerName, r.TotalScore, r.CurrentThreeDA)
			fmt.Println(g.GetGameSummary())
			return
		case model.BustTurn:
			fmt.Printf("Player %s busted\n\n", r.PlayerName)
		case model.NormalTurn:
			fmt.Printf("Player %s scored %d points(3DA: %.2f) remaining score: %d\n\n",
				r.PlayerName, r.TotalScore,
				r.CurrentThreeDA, r.RemainingScore)

		}
	}
}
