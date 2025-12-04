package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/kregan77/dartbuddy/internal/model"
)

// Server holds the HTTP server and game state
type Server struct {
	games map[uuid.UUID]*GameState
	mu    sync.RWMutex
}

// GameState wraps a Game01 with additional metadata
type GameState struct {
	Game       *model.Game01
	IsRealGame bool // true if any real players are in the game
	OutChart   *model.OutChart
}

// NewServer creates a new API server
func NewServer() *Server {
	return &Server{
		games: make(map[uuid.UUID]*GameState),
	}
}

// CreateGameRequest represents a request to create a new game
type CreateGameRequest struct {
	StartingScore int  `json:"starting_score"`
	UseOutChart   bool `json:"use_out_chart"` // if true, players use the out chart
}

// CreateGameResponse represents the response from creating a game
type CreateGameResponse struct {
	GameID        string `json:"game_id"`
	StartingScore int    `json:"starting_score"`
}

// AddPlayerRequest represents a request to add a player
type AddPlayerRequest struct {
	Name              string  `json:"name"`
	IsSimulated       bool    `json:"is_simulated"`
	ThreeDA           float64 `json:"three_da"`           // Only for simulated players
	ScoringPreference string  `json:"scoring_preference"` // "twenties" or "nineteens"
}

// AddPlayerResponse represents the response from adding a player
type AddPlayerResponse struct {
	PlayerID string `json:"player_id"`
	Name     string `json:"name"`
}

// SubmitScoreRequest represents a real player submitting their score
type SubmitScoreRequest struct {
	Scores []int `json:"scores"` // Array of 1-3 scores for the turn
}

// GameStateResponse represents the current state of the game
type GameStateResponse struct {
	GameID         string          `json:"game_id"`
	Turn           int             `json:"turn"`
	CurrentPlayer  PlayerState     `json:"current_player"`
	Players        []PlayerState   `json:"players"`
	LastTurnResult *TurnResultData `json:"last_turn_result,omitempty"`
	GameOver       bool            `json:"game_over"`
	Winner         string          `json:"winner,omitempty"`
}

// PlayerState represents the state of a player
type PlayerState struct {
	PlayerID     string  `json:"player_id"`
	Name         string  `json:"name"`
	CurrentScore int     `json:"current_score"`
	IsSimulated  bool    `json:"is_simulated"`
	ThreeDA      float64 `json:"three_da"`
	Turns        int     `json:"turns"`
	TotalPoints  int     `json:"total_points"`
	AverageScore float64 `json:"average_score"`
}

// TurnResultData represents the result of a turn
type TurnResultData struct {
	PlayerName     string  `json:"player_name"`
	TotalScore     int     `json:"total_score"`
	RemainingScore int     `json:"remaining_score"`
	ThreeDA        float64 `json:"three_da"`
}

// CreateGame handles POST /games
func (s *Server) CreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	// Default to 501
	if req.StartingScore == 0 {
		req.StartingScore = 501
	}

	game := model.New01Game(req.StartingScore)
	gameState := &GameState{
		Game:       game,
		IsRealGame: false,
	}

	if req.UseOutChart {
		gameState.OutChart = model.NewOutChart()
	}

	s.mu.Lock()
	s.games[game.ID] = gameState
	s.mu.Unlock()

	resp := CreateGameResponse{
		GameID:        game.ID.String(),
		StartingScore: req.StartingScore,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// AddPlayer handles POST /games/{id}/players
func (s *Server) AddPlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract game ID from URL path
	gameIDStr := r.URL.Path[len("/games/"):]
	if idx := len(gameIDStr) - len("/players"); idx > 0 {
		gameIDStr = gameIDStr[:idx]
	}

	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	var req AddPlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	gameState, exists := s.games[gameID]
	s.mu.Unlock()

	if !exists {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	// Parse scoring preference
	var pref model.ScoringPreference
	if req.ScoringPreference == "nineteens" {
		pref = model.NinteensScoringPreference
	} else {
		pref = model.TwentiesScoringPreference
	}

	// For simulated players, use the provided 3DA
	// For real players, use a default 3DA (won't be used for targeting)
	threeDA := req.ThreeDA
	if !req.IsSimulated {
		gameState.IsRealGame = true
		if threeDA == 0 {
			threeDA = 60.0 // Default for real players (just for display)
		}
	} else if threeDA == 0 {
		threeDA = 60.0 // Default for simulated players
	}

	profile := model.NewPlayer(req.Name, threeDA, pref)
	gameState.Game.AddPlayer(profile)

	resp := AddPlayerResponse{
		PlayerID: profile.ID.String(),
		Name:     profile.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// PlaySimulatedTurn handles POST /games/{id}/turns/simulate
func (s *Server) PlaySimulatedTurn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract game ID from URL path
	gameIDStr := r.URL.Path[len("/games/"):]
	if idx := len(gameIDStr) - len("/turns/simulate"); idx > 0 {
		gameIDStr = gameIDStr[:idx]
	}

	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	gameState, exists := s.games[gameID]
	s.mu.Unlock()

	if !exists {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	if err := gameState.Game.Start(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to start game: %v", err), http.StatusBadRequest)
		return
	}

	result, won := gameState.Game.PlayTurn()

	resp := s.buildGameStateResponse(gameState, result, won)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// SubmitScore handles POST /games/{id}/turns/submit
func (s *Server) SubmitScore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract game ID from URL path
	gameIDStr := r.URL.Path[len("/games/"):]
	if idx := len(gameIDStr) - len("/turns/submit"); idx > 0 {
		gameIDStr = gameIDStr[:idx]
	}

	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	var req SubmitScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	gameState, exists := s.games[gameID]
	s.mu.Unlock()

	if !exists {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	if err := gameState.Game.Start(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to start game: %v", err), http.StatusBadRequest)
		return
	}

	// Process the submitted scores
	player := gameState.Game.GetCurrentPlayer()
	totalScore := 0
	won := false

	for _, score := range req.Scores {
		player.CurrentScore -= score
		player.TotalPoints += score
		player.Throws++
		totalScore += score

		if player.CurrentScore <= 0 {
			won = true
			break
		}
	}

	result := &model.TurnResult{
		PlayerName:     player.Profile.GetName(),
		TotalScore:     totalScore,
		RemainingScore: player.CurrentScore,
		ThreeDA:        float64(totalScore) / float64(len(req.Scores)),
	}

	if !won {
		gameState.Game.Turn++
		gameState.Game.NextPlayer()
	}

	resp := s.buildGameStateResponse(gameState, result, won)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetGameState handles GET /games/{id}
func (s *Server) GetGameState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract game ID from URL path
	gameIDStr := r.URL.Path[len("/games/"):]

	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	gameState, exists := s.games[gameID]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	resp := s.buildGameStateResponse(gameState, nil, false)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// buildGameStateResponse constructs a GameStateResponse from the current game state
func (s *Server) buildGameStateResponse(gameState *GameState, lastResult *model.TurnResult, gameOver bool) GameStateResponse {
	players := make([]PlayerState, len(gameState.Game.Players))
	for i, p := range gameState.Game.Players {
		avgScore := 0.0
		if p.Throws > 0 {
			avgScore = float64(p.TotalPoints) / float64(p.Throws) * 3.0
		}

		players[i] = PlayerState{
			PlayerID:     p.Profile.ID.String(),
			Name:         p.Profile.GetName(),
			CurrentScore: p.CurrentScore,
			IsSimulated:  true, // TODO: track this properly
			ThreeDA:      p.Profile.GetThreeDA(),
			Turns:        p.Turns,
			TotalPoints:  p.TotalPoints,
			AverageScore: avgScore,
		}
	}

	resp := GameStateResponse{
		GameID:   gameState.Game.ID.String(),
		Turn:     gameState.Game.Turn,
		Players:  players,
		GameOver: gameOver,
	}

	if len(players) > 0 {
		resp.CurrentPlayer = players[gameState.Game.CurrentPlayer]
	}

	if lastResult != nil {
		resp.LastTurnResult = &TurnResultData{
			PlayerName:     lastResult.PlayerName,
			TotalScore:     lastResult.TotalScore,
			RemainingScore: lastResult.RemainingScore,
			ThreeDA:        lastResult.ThreeDA,
		}

		if gameOver {
			resp.Winner = lastResult.PlayerName
		}
	}

	return resp
}

// RegisterRoutes registers all API routes on the given mux
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/games", s.CreateGame)
	mux.HandleFunc("/games/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Route to appropriate handler based on path
		if path[len(path)-8:] == "/players" {
			s.AddPlayer(w, r)
		} else if len(path) > 16 && path[len(path)-16:] == "/turns/simulate" {
			s.PlaySimulatedTurn(w, r)
		} else if len(path) > 13 && path[len(path)-13:] == "/turns/submit" {
			s.SubmitScore(w, r)
		} else {
			// Just game ID - get state
			s.GetGameState(w, r)
		}
	})
}
