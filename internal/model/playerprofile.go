package model

import (
	"fmt"

	"github.com/google/uuid"
)

type PlayerType int

const (
	RealPlayer PlayerType = iota
	SimulatedPlayer
)

type ScoringPreference int

const (
	TwentiesScoringPreference ScoringPreference = iota
	NinteensScoringPreference
)

type PlayerProfile struct {
	ID                uuid.UUID
	Name              string
	PlayerType        PlayerType
	ThreeDA           float64 // x01 Three Dart Average - currently only applicable to simulated players.
	ScoringPreference ScoringPreference
}

func NewPlayer(name string, threeDA float64, scoringPreference ScoringPreference) *PlayerProfile {
	return &PlayerProfile{
		ID:                uuid.New(),
		Name:              name,
		ThreeDA:           threeDA,
		ScoringPreference: scoringPreference,
	}
}

func (p *PlayerProfile) GetName() string {
	return p.Name
}

func (p *PlayerProfile) GetType() PlayerType {
	return p.PlayerType
}

func (p *PlayerProfile) GetThreeDA() float64 {
	return p.ThreeDA
}

func (p *PlayerProfile) GetScoringPreference() ScoringPreference {
	return p.ScoringPreference
}
