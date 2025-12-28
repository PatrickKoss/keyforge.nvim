package nvim

// JSON-RPC 2.0 message types for Neovim communication

// Request represents a JSON-RPC request
type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      int         `json:"id,omitempty"`
}

// Response represents a JSON-RPC response
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      int         `json:"id"`
}

// Notification represents a JSON-RPC notification (no ID, no response expected)
type Notification struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// RPCError represents a JSON-RPC error
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Standard JSON-RPC error codes
const (
	ErrCodeParseError     = -32700
	ErrCodeInvalidRequest = -32600
	ErrCodeMethodNotFound = -32601
	ErrCodeInvalidParams  = -32602
	ErrCodeInternalError  = -32603
)

// Game -> Neovim message types

// ChallengeRequest asks Neovim to present a challenge to the user
type ChallengeRequest struct {
	RequestID  string `json:"request_id"`
	Category   string `json:"category"`
	Difficulty int    `json:"difficulty"`
}

// GameStateUpdate notifies Neovim of game state changes
type GameStateUpdate struct {
	State    string `json:"state"` // playing, paused, game_over, victory
	Wave     int    `json:"wave"`
	Gold     int    `json:"gold"`
	Health   int    `json:"health"`
	Enemies  int    `json:"enemies"`
	Towers   int    `json:"towers"`
}

// Neovim -> Game message types

// ChallengeResult contains the result of a completed challenge
type ChallengeResult struct {
	RequestID      string  `json:"request_id"`
	Success        bool    `json:"success"`
	Skipped        bool    `json:"skipped,omitempty"`
	KeystrokeCount int     `json:"keystroke_count,omitempty"`
	TimeMs         int     `json:"time_ms,omitempty"`
	Efficiency     float64 `json:"efficiency,omitempty"`
}

// ConfigUpdate contains configuration from Neovim
type ConfigUpdate struct {
	Difficulty     string `json:"difficulty,omitempty"`
	UseNerdFonts   bool   `json:"use_nerd_fonts,omitempty"`
	StartingGold   int    `json:"starting_gold,omitempty"`
	StartingHealth int    `json:"starting_health,omitempty"`
}

// Method names
const (
	// Game -> Neovim
	MethodRequestChallenge = "request_challenge"
	MethodGameStateUpdate  = "game_state_update"
	MethodGameReady        = "game_ready"

	// Neovim -> Game
	MethodChallengeComplete = "challenge_complete"
	MethodConfigUpdate      = "config_update"
	MethodPauseGame         = "pause_game"
	MethodResumeGame        = "resume_game"
)

// NewRequest creates a new JSON-RPC request
func NewRequest(id int, method string, params interface{}) *Request {
	return &Request{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      id,
	}
}

// NewNotification creates a new JSON-RPC notification
func NewNotification(method string, params interface{}) *Notification {
	return &Notification{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}
}

// NewResponse creates a new JSON-RPC response
func NewResponse(id int, result interface{}, err *RPCError) *Response {
	return &Response{
		JSONRPC: "2.0",
		Result:  result,
		Error:   err,
		ID:      id,
	}
}

// NewError creates a new RPC error
func NewError(code int, message string) *RPCError {
	return &RPCError{
		Code:    code,
		Message: message,
	}
}
