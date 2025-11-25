package model

import (
	"github.com/google/uuid"
)

type Player struct {
	Profile *PlayerProfile
	Score   int
}

type Game struct {
	ID         uuid.UUID
	Players    []*Player
	StartScore int
	Turn       int
}

func NewGame(startingScore int) {
	return &Game{
		ID:         uuid.New(),
		StartScore: startingScore,
	}
}

func (g *Game) AddPlayer(profile *PlayerProfile) {
	player := &Player{
		Profile: profile,
		Score:   g.StartScore,
	}
	g.Players = append(g.Players, player)
}

func (g *Game) NextTurn() {
	g.Turn = (g.Turn + 1) % len(g.Players)
}

func (g *Game) CurrentPlayer() *Player {
	return g.Players[g.Turn]
}

func (g *Game) PlayTurn() {
	p := g.CurrentPlayer()
	p.ThrowDarts(p.Score)
}

func (p *Player) ThrowDarts(currentScore int) {
	p.Profile.ThrowDarts(currentScore)
}
