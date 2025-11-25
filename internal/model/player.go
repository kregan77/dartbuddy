package model

import (
	"github.com/google/uuid"
)

type ScoringPreference int

const (
	TwentiesScoringPreference ScoringPreference = iota
	NinteensScoringPreference
)

type PlayerProfile struct {
	ID                uuid.UUID
	Name              string
	threeDA           float64
	scoringPreference ScoringPreference
}

func NewPlayer(name string, threeDA float64) *PlayerProfile {
	return &PlayerProfile{
		Name:    name,
		threeDA: threeDA,
	}
}

func (p *PlayerProfile) ThrowDarts(currentScore int) {
	nextTarget := p.CalculateNextTarget(currentScore)

}

func (p *PlayerProfile) CalculateNextTarget(currentScore int) int {
	// TODO: make the player not dumb
	switch p.scoringPreference {
	case TwentiesScoringPreference:
		return Twenty
	case NinteensScoringPreference:
		return Nineteen
	default:
		return Twenty
	}
}
