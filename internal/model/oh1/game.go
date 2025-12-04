package oh1

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kregan77/dartbuddy/internal/model"
)

// Player tracks the 01 player state.
type Player struct {
	model.PlayerProfile
	spread       float64
	CurrentScore int
	Turns        int
	TotalPoints  int
	Throws       int
}

func (p *Player) CurrentThreeDA() float64 {
	if p.Throws == 0 {
		return 0.0
	}
	return float64(p.TotalPoints) / float64(p.Throws/3)
}

type Game struct {
	ID            uuid.UUID
	Simulator     *model.Simulator
	Players       []*Player
	CurrentPlayer int
	StartScore    int
	Turn          int
	Outs          *OutChart
}

func New01Game(startingScore int) *Game {
	return &Game{
		ID:         uuid.New(),
		StartScore: startingScore,
		Simulator:  model.NewSimulator(),
		Outs:       NewOutChart(),
	}
}

func (g *Game) Start() error {
	if len(g.Players) == 0 {
		return errors.New("cannot start game with no players")
	}
	// TODO: initialization?
	return nil
}

// TODO: for now the game will go in order of player added.
func (g *Game) AddPlayer(profile *model.PlayerProfile) {
	player := &Player{
		PlayerProfile: *profile,
		CurrentScore:  g.StartScore,
	}
	g.Players = append(g.Players, player)
}

func (g *Game) GetCurrentPlayer() *Player {
	if len(g.Players) == 0 {
		return nil
	}
	return g.Players[g.CurrentPlayer]
}

func (g *Game) NextPlayer() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
}

type TurnResultType int

const (
	ScoringTurn TurnResultType = iota
	BustTurn
	WinTurn
)

type TurnResult struct {
	Type           TurnResultType
	PlayerName     string
	Results        []*model.DartResult
	TotalScore     int
	RemainingScore int
	CurrentThreeDA float64
}

func (g *Game) PlayTurn() *TurnResult {
	p := g.GetCurrentPlayer()
	if p.GetType() == model.RealPlayer {
		fmt.Printf("awaiting real player score submission for %s...\n", p.GetName())
		return nil
	}

	fmt.Printf("%s turn.  Current Score: %d; Leg 3DA: %.2f\n",
		p.GetName(), p.CurrentScore, p.CurrentThreeDA())
	totalScore := 0
	currentScore := p.CurrentScore
	results := make([]*model.DartResult, 0, 3)
	for dart := range 3 {
		result := g.ThrowDart(dart, currentScore, p)

		results = append(results, result)
		currentScore -= result.Score
		if currentScore < 0 ||
			(currentScore == 0 && result.GetMultiplier() != model.Double) ||
			currentScore == 1 {
			fmt.Printf("	BUST!  Score resets to %d\n", p.CurrentScore)
			return &TurnResult{
				Type:           BustTurn,
				PlayerName:     p.GetName(),
				Results:        results,
				TotalScore:     totalScore,
				CurrentThreeDA: p.CurrentThreeDA(),
			}
		}
		// either win or continue
		p.TotalPoints += result.Score
		p.Throws++
		totalScore += result.Score

		if currentScore == 0 {
			fmt.Println("	WIN!!")
			return &TurnResult{
				Type:           WinTurn,
				PlayerName:     p.GetName(),
				Results:        results,
				TotalScore:     totalScore,
				CurrentThreeDA: p.CurrentThreeDA(),
				RemainingScore: currentScore,
			}
		}
	}

	p.CurrentScore = currentScore
	g.Turn++
	g.NextPlayer()
	return &TurnResult{
		Type:           ScoringTurn,
		PlayerName:     p.GetName(),
		TotalScore:     totalScore,
		Results:        results,
		CurrentThreeDA: float64(p.TotalPoints) / float64(p.Throws/3),
		RemainingScore: p.CurrentScore,
	}
}

func (g *Game) ThrowDart(dart int, currentScore int, p *Player) *model.DartResult {
	target := g.Outs.GetNextTarget(currentScore, p.GetScoringPreference())
	result := g.Simulator.ThrowDart(target, p.GetSpread())
	fmt.Printf("	Dart %d(target: %s): %s\n", dart+1,
		target.String(),
		result.String())
	return result
}

func (g *Game) GetGameSummary() string {
	summary := "Game Summary:\n"
	for _, p := range g.Players {
		summary += fmt.Sprintf("%s:\n\tFinal Score: %d, Total Points: %d, Throws: %d, 3DA: %.2f\n",
			p.GetName(), p.CurrentScore, p.TotalPoints, p.Throws,
			float64(p.TotalPoints)/float64(p.Throws/3))
	}
	return summary
}
