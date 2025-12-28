package nvim

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

// Client handles JSON-RPC communication with Neovim
type Client struct {
	reader    *bufio.Reader
	writer    io.Writer
	requestID int
	mu        sync.Mutex

	// Pending requests waiting for responses
	pending map[int]chan *Response

	// Handler for incoming requests/notifications from Neovim
	handler Handler

	// Channel for incoming messages
	incoming chan interface{}

	// Done channel for shutdown
	done chan struct{}
}

// Handler processes incoming RPC messages from Neovim
type Handler interface {
	HandleChallengeComplete(result *ChallengeResult)
	HandleConfigUpdate(config *ConfigUpdate)
	HandlePause()
	HandleResume()
	HandleStartChallenge()
	HandleRestart()
}

// NewClient creates a new RPC client using stdin for reading and stderr for writing
// This allows Bubbletea to use stdout for terminal rendering while RPC uses stderr
func NewClient(handler Handler) *Client {
	return &Client{
		reader:   bufio.NewReader(os.Stdin),
		writer:   os.Stderr, // Use stderr so stdout remains free for terminal
		pending:  make(map[int]chan *Response),
		handler:  handler,
		incoming: make(chan interface{}, 100),
		done:     make(chan struct{}),
	}
}

// NewClientWithIO creates a client with custom IO (for testing)
func NewClientWithIO(r io.Reader, w io.Writer, handler Handler) *Client {
	return &Client{
		reader:   bufio.NewReader(r),
		writer:   w,
		pending:  make(map[int]chan *Response),
		handler:  handler,
		incoming: make(chan interface{}, 100),
		done:     make(chan struct{}),
	}
}

// Start begins listening for incoming messages
func (c *Client) Start() {
	go c.readLoop()
	go c.processLoop()
}

// Stop shuts down the client
func (c *Client) Stop() {
	close(c.done)
}

// readLoop continuously reads messages from the input
func (c *Client) readLoop() {
	for {
		select {
		case <-c.done:
			return
		default:
		}

		line, err := c.reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				// Log error but continue
			}
			continue
		}

		if len(line) == 0 {
			continue
		}

		// Try to parse as a generic message first
		var msg map[string]interface{}
		if err := json.Unmarshal(line, &msg); err != nil {
			continue
		}

		// Determine message type
		if _, hasMethod := msg["method"]; hasMethod {
			// It's a request or notification
			if _, hasID := msg["id"]; hasID {
				// Request (has ID)
				var req Request
				if err := json.Unmarshal(line, &req); err == nil {
					c.incoming <- &req
				}
			} else {
				// Notification (no ID)
				var notif Notification
				if err := json.Unmarshal(line, &notif); err == nil {
					c.incoming <- &notif
				}
			}
		} else if _, hasResult := msg["result"]; hasResult {
			// It's a response
			var resp Response
			if err := json.Unmarshal(line, &resp); err == nil {
				c.incoming <- &resp
			}
		} else if _, hasError := msg["error"]; hasError {
			// It's an error response
			var resp Response
			if err := json.Unmarshal(line, &resp); err == nil {
				c.incoming <- &resp
			}
		}
	}
}

// processLoop handles incoming messages
func (c *Client) processLoop() {
	for {
		select {
		case <-c.done:
			return
		case msg := <-c.incoming:
			c.handleMessage(msg)
		}
	}
}

func (c *Client) handleMessage(msg interface{}) {
	switch m := msg.(type) {
	case *Response:
		c.handleResponse(m)
	case *Request:
		c.handleRequest(m)
	case *Notification:
		c.handleNotification(m)
	}
}

func (c *Client) handleResponse(resp *Response) {
	c.mu.Lock()
	ch, ok := c.pending[resp.ID]
	if ok {
		delete(c.pending, resp.ID)
	}
	c.mu.Unlock()

	if ok {
		ch <- resp
	}
}

func (c *Client) handleRequest(req *Request) {
	// Handle incoming request and send response
	var result interface{}
	var rpcErr *RPCError

	switch req.Method {
	case MethodChallengeComplete:
		if params, ok := req.Params.(map[string]interface{}); ok {
			cr := parseChallengeResult(params)
			if c.handler != nil {
				c.handler.HandleChallengeComplete(cr)
			}
			result = map[string]bool{"ok": true}
		}
	case MethodConfigUpdate:
		if params, ok := req.Params.(map[string]interface{}); ok {
			cfg := parseConfigUpdate(params)
			if c.handler != nil {
				c.handler.HandleConfigUpdate(cfg)
			}
			result = map[string]bool{"ok": true}
		}
	case MethodPauseGame:
		if c.handler != nil {
			c.handler.HandlePause()
		}
		result = map[string]bool{"ok": true}
	case MethodResumeGame:
		if c.handler != nil {
			c.handler.HandleResume()
		}
		result = map[string]bool{"ok": true}
	case MethodStartChallenge:
		if c.handler != nil {
			c.handler.HandleStartChallenge()
		}
		result = map[string]bool{"ok": true}
	case MethodRestartGame:
		if c.handler != nil {
			c.handler.HandleRestart()
		}
		result = map[string]bool{"ok": true}
	default:
		rpcErr = NewError(ErrCodeMethodNotFound, fmt.Sprintf("method not found: %s", req.Method))
	}

	// Send response
	resp := NewResponse(req.ID, result, rpcErr)
	c.send(resp)
}

func (c *Client) handleNotification(notif *Notification) {
	// Handle notification (no response needed)
	switch notif.Method {
	case MethodChallengeComplete:
		if params, ok := notif.Params.(map[string]interface{}); ok {
			cr := parseChallengeResult(params)
			if c.handler != nil {
				c.handler.HandleChallengeComplete(cr)
			}
		}
	case MethodConfigUpdate:
		if params, ok := notif.Params.(map[string]interface{}); ok {
			cfg := parseConfigUpdate(params)
			if c.handler != nil {
				c.handler.HandleConfigUpdate(cfg)
			}
		}
	case MethodPauseGame:
		if c.handler != nil {
			c.handler.HandlePause()
		}
	case MethodResumeGame:
		if c.handler != nil {
			c.handler.HandleResume()
		}
	case MethodStartChallenge:
		if c.handler != nil {
			c.handler.HandleStartChallenge()
		}
	case MethodRestartGame:
		if c.handler != nil {
			c.handler.HandleRestart()
		}
	}
}

// send writes a message to the output
func (c *Client) send(msg interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(c.writer, "%s\n", data)
	return err
}

// Request sends a request and waits for a response
func (c *Client) Request(method string, params interface{}) (*Response, error) {
	c.mu.Lock()
	c.requestID++
	id := c.requestID
	ch := make(chan *Response, 1)
	c.pending[id] = ch
	c.mu.Unlock()

	req := NewRequest(id, method, params)
	if err := c.send(req); err != nil {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return nil, err
	}

	// Wait for response (with timeout handled by caller if needed)
	resp := <-ch
	return resp, nil
}

// Notify sends a notification (no response expected)
func (c *Client) Notify(method string, params interface{}) error {
	notif := NewNotification(method, params)
	return c.send(notif)
}

// RequestChallenge asks Neovim to present a challenge
func (c *Client) RequestChallenge(requestID, category string, difficulty int) error {
	return c.Notify(MethodRequestChallenge, &ChallengeRequest{
		RequestID:  requestID,
		Category:   category,
		Difficulty: difficulty,
	})
}

// SendGameState sends current game state to Neovim
func (c *Client) SendGameState(state string, wave, gold, health, enemies, towers int) error {
	return c.Notify(MethodGameStateUpdate, &GameStateUpdate{
		State:   state,
		Wave:    wave,
		Gold:    gold,
		Health:  health,
		Enemies: enemies,
		Towers:  towers,
	})
}

// SendGameReady notifies Neovim that the game is ready
func (c *Client) SendGameReady() error {
	return c.Notify(MethodGameReady, nil)
}

// SendGoldUpdate notifies Neovim of gold changes
func (c *Client) SendGoldUpdate(gold, earned int, source string, speedBonus float64) error {
	return c.Notify(MethodGoldUpdate, &GoldUpdate{
		Gold:       gold,
		Earned:     earned,
		Source:     source,
		SpeedBonus: speedBonus,
	})
}

// SendChallengeAvailable notifies Neovim that challenges are available
func (c *Client) SendChallengeAvailable(count, nextReward int, nextCategory string) error {
	return c.Notify(MethodChallengeAvailable, &ChallengeAvailable{
		Count:        count,
		NextReward:   nextReward,
		NextCategory: nextCategory,
	})
}

// SendGameOver notifies Neovim that the game is over
func (c *Client) SendGameOver(wave, gold, towers, health int) error {
	return c.Notify(MethodGameOver, &GameOverParams{
		Wave:   wave,
		Gold:   gold,
		Towers: towers,
		Health: health,
	})
}

// SendVictory notifies Neovim that the player has won
func (c *Client) SendVictory(wave, gold, towers, health int) error {
	return c.Notify(MethodVictory, &VictoryParams{
		Wave:   wave,
		Gold:   gold,
		Towers: towers,
		Health: health,
	})
}

// Helper functions to parse params

func parseChallengeResult(params map[string]interface{}) *ChallengeResult {
	cr := &ChallengeResult{}
	if v, ok := params["request_id"].(string); ok {
		cr.RequestID = v
	}
	if v, ok := params["success"].(bool); ok {
		cr.Success = v
	}
	if v, ok := params["skipped"].(bool); ok {
		cr.Skipped = v
	}
	if v, ok := params["keystroke_count"].(float64); ok {
		cr.KeystrokeCount = int(v)
	}
	if v, ok := params["time_ms"].(float64); ok {
		cr.TimeMs = int(v)
	}
	if v, ok := params["efficiency"].(float64); ok {
		cr.Efficiency = v
	}
	return cr
}

func parseConfigUpdate(params map[string]interface{}) *ConfigUpdate {
	cfg := &ConfigUpdate{}
	if v, ok := params["difficulty"].(string); ok {
		cfg.Difficulty = v
	}
	if v, ok := params["use_nerd_fonts"].(bool); ok {
		cfg.UseNerdFonts = v
	}
	if v, ok := params["starting_gold"].(float64); ok {
		cfg.StartingGold = int(v)
	}
	if v, ok := params["starting_health"].(float64); ok {
		cfg.StartingHealth = int(v)
	}
	return cfg
}
