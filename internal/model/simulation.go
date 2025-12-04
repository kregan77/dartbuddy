package model

import (
	"fmt"
	"math"
	"math/rand"
)

// Dartboard dimensions (in mm, standard dartboard)
const (
	DoubleBullRadius  = 6.35  // Inner bull (double bull)
	SingleBullRadius  = 15.9  // Outer bull (single bull)
	TripleInnerRadius = 99.0  // Inner edge of triple ring
	TripleOuterRadius = 107.0 // Outer edge of triple ring
	DoubleInnerRadius = 162.0 // Inner edge of double ring
	DoubleOuterRadius = 170.0 // Outer edge of double ring
	BoardRadius       = 170.0 // Outer edge of board
)

// Multiplier represents the multiplier for a dart hit
type Multiplier int

func (m Multiplier) String() string {
	switch m {
	case Miss:
		return "Miss"
	case Single:
		return "Single"
	case Double:
		return "Double"
	case Triple:
		return "Triple"
	default:
		return "Unknown"
	}
}

const (
	Miss   Multiplier = 0
	Single Multiplier = 1
	Double Multiplier = 2
	Triple Multiplier = 3
)

// DartTarget represents a single dart target with multiplier and number
type DartTarget struct {
	Multiplier Multiplier
	Number     int
}

func (t *DartTarget) GetMultiplier() Multiplier {
	return t.Multiplier
}

func (t *DartTarget) GetNumber() int {
	return t.Number
}

func (t *DartTarget) String() string {
	return fmt.Sprintf("%s %d", t.GetMultiplier().String(), t.GetNumber())
}

// DartResult represents the outcome of a single dart throw
type DartResult struct {
	DartTarget
	Score int
}

func (d *DartResult) String() string {
	return fmt.Sprintf("%s %d(%d)", d.GetMultiplier().String(), d.GetNumber(), d.Score)
}

// Simulator handles dart throw simulation with cached dartboard geometry
type Simulator struct {
	angleMap map[int]float64 // maps number to center angle in radians
	rng      *rand.Rand
}

// NewSimulator creates and initializes a new dart simulator
func NewSimulator() *Simulator {
	sim := &Simulator{
		angleMap: make(map[int]float64),
		rng:      rand.New(rand.NewSource(rand.Int63())),
	}
	sim.initializeDartboard()
	return sim
}

// initializeDartboard pre-computes and caches angles for all numbers
func (s *Simulator) initializeDartboard() {
	// Each segment is 18 degrees (2π/20 radians)
	segmentAngle := 2.0 * math.Pi / 20.0

	for i, number := range Board {
		// Calculate center angle for this segment
		// 0 radians is at top where 20 is located
		s.angleMap[number] = float64(i) * segmentAngle
	}

	// Bullseye doesn't have a specific angle
	s.angleMap[Bullseye] = 0.0
}

// getAngleForNumber returns the center angle (in radians) for a given number
func (s *Simulator) getAngleForNumber(number int) float64 {
	if angle, ok := s.angleMap[number]; ok {
		return angle
	}
	return 0.0 // Default to top if number not found
}

// getNumberFromAngle returns the number at a given angle (in radians)
func (s *Simulator) getNumberFromAngle(angle float64) int {
	// Normalize angle to [0, 2π)
	for angle < 0 {
		angle += 2.0 * math.Pi
	}
	for angle >= 2.0*math.Pi {
		angle -= 2.0 * math.Pi
	}

	// Each segment is 18 degrees (2π/20 radians)
	segmentAngle := 2.0 * math.Pi / 20.0

	// Find which segment this angle falls into
	segmentIndex := int(math.Round(angle/segmentAngle)) % 20

	return Board[segmentIndex]
}

// getTargetRadius returns the target radius for a given multiplier
func (s *Simulator) getTargetRadius(multiplier Multiplier) float64 {
	switch multiplier {
	case Triple:
		return (TripleInnerRadius + TripleOuterRadius) / 2.0
	case Double:
		return (DoubleInnerRadius + DoubleOuterRadius) / 2.0
	case Single:
		// Aim at the larger single area (between triple and double)
		return (TripleOuterRadius + DoubleInnerRadius) / 2.0
	default:
		return (TripleInnerRadius + TripleOuterRadius) / 2.0
	}
}

// ThrowDart simulates a single dart throw using 2D Gaussian distribution
func (s *Simulator) ThrowDart(target DartTarget, spread float64) *DartResult {
	// Special case: aiming at bullseye
	if target.Number == Bullseye {
		// Aim at center (radius 0)
		hitX := s.rng.NormFloat64() * spread
		hitY := s.rng.NormFloat64() * spread
		hitRadius := math.Sqrt(hitX*hitX + hitY*hitY)
		hitAngle := math.Atan2(hitX, hitY)
		return s.determineHit(hitRadius, hitAngle)
	}

	// Get target coordinates
	targetAngle := s.getAngleForNumber(target.Number)
	targetRadius := s.getTargetRadius(target.Multiplier)

	// Convert to Cartesian coordinates
	targetX := targetRadius * math.Sin(targetAngle)
	targetY := targetRadius * math.Cos(targetAngle)

	// Sample from 2D Gaussian distribution
	hitX := targetX + s.rng.NormFloat64()*spread
	hitY := targetY + s.rng.NormFloat64()*spread

	// Convert back to polar coordinates
	hitRadius := math.Sqrt(hitX*hitX + hitY*hitY)
	hitAngle := math.Atan2(hitX, hitY)

	// Determine what was hit
	return s.determineHit(hitRadius, hitAngle)
}

// determineHit determines what segment was hit based on radius and angle
func (s *Simulator) determineHit(radius, angle float64) *DartResult {
	var number int
	var multiplier Multiplier

	// Check radial distance first
	if radius <= DoubleBullRadius {
		// Double bull (50 points)
		return &DartResult{
			DartTarget: DartTarget{
				Number:     Bullseye,
				Multiplier: Double,
			},
			Score: 50,
		}
	} else if radius <= SingleBullRadius {
		// Single bull (25 points)
		return &DartResult{
			DartTarget: DartTarget{
				Number:     Bullseye,
				Multiplier: Single,
			},
			Score: 25,
		}
	} else if radius > BoardRadius {
		// Miss
		return &DartResult{
			DartTarget: DartTarget{
				Number:     0,
				Multiplier: Miss,
			},
			Score: 0,
		}
	}

	// Determine number from angle
	number = s.getNumberFromAngle(angle)

	// Determine multiplier from radius
	if radius >= DoubleInnerRadius && radius <= DoubleOuterRadius {
		multiplier = Double
	} else if radius >= TripleInnerRadius && radius <= TripleOuterRadius {
		multiplier = Triple
	} else {
		multiplier = Single
	}

	score := number * int(multiplier)

	return &DartResult{
		DartTarget: DartTarget{
			Number:     number,
			Multiplier: multiplier,
		},
		Score: score,
	}
}

// calculateSpread converts a 3DA to a standard deviation (spread) in mm
func (s *Simulator) calculateSpread(threeDA float64) float64 {
	// Professional players (90+ 3DA) should have tight grouping (~10-15mm)
	// Intermediate players (60 3DA) should have moderate spread (~20-25mm)
	// Beginners (30 3DA) should have wide spread (~35-40mm)

	// Using an inverse relationship: spread = k / (threeDA + offset)
	// This formula produces realistic spreads relative to the 8mm triple ring width
	const k = 1800.0
	const offset = 20.0

	spread := k / (threeDA + offset)

	// Clamp to reasonable values (5mm to 50mm)
	if spread < 5.0 {
		spread = 5.0
	}
	if spread > 50.0 {
		spread = 50.0
	}

	fmt.Printf("Calculated spread for 3DA %.2f is %.2f mm\n", threeDA, spread)
	return spread
}
