package nvim

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
)

// Debug logging to file (since stdout is used by bubbletea).
var debugFile *os.File

func init() {
	// Create debug log file in /tmp with restricted permissions
	var err error
	debugFile, err = os.OpenFile("/tmp/keyforge-debug.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		debugFile = nil
	}
}

func debugLog(format string, args ...interface{}) {
	if debugFile != nil {
		fmt.Fprintf(debugFile, format+"\n", args...)
		_ = debugFile.Sync() // Ignore sync errors for debug logging
	}
}

// SocketServer handles JSON-RPC communication over a Unix socket.
type SocketServer struct {
	socketPath string
	listener   net.Listener
	conn       net.Conn
	reader     *bufio.Reader
	mu         sync.Mutex

	// Pending requests waiting for responses
	pending map[int]chan *Response

	// Handler for incoming requests/notifications from Neovim
	handler Handler

	// Channel for incoming messages
	incoming chan interface{}

	// Done channel for shutdown
	done chan struct{}

	// Connected state
	connected bool
}

// NewSocketServer creates a new RPC server that listens on a Unix socket.
func NewSocketServer(socketPath string, handler Handler) *SocketServer {
	return &SocketServer{
		socketPath: socketPath,
		pending:    make(map[int]chan *Response),
		handler:    handler,
		incoming:   make(chan interface{}, 100),
		done:       make(chan struct{}),
	}
}

// Start begins listening on the socket and waits for a connection.
func (s *SocketServer) Start() error {
	var err error
	s.listener, err = net.Listen("unix", s.socketPath)
	if err != nil {
		return fmt.Errorf("failed to listen on socket: %w", err)
	}

	// Accept connections in the background
	go s.acceptLoop()
	go s.processLoop()

	return nil
}

// Stop shuts down the server.
func (s *SocketServer) Stop() {
	close(s.done)

	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			debugLog("socket: error closing connection: %v", err)
		}
	}
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			debugLog("socket: error closing listener: %v", err)
		}
	}
}

// IsConnected returns true if a client is connected.
func (s *SocketServer) IsConnected() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.connected
}

// acceptLoop waits for incoming connections.
func (s *SocketServer) acceptLoop() {
	debugLog("acceptLoop started, waiting for connections on %s", s.socketPath)
	for {
		select {
		case <-s.done:
			return
		default:
		}

		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.done:
				return
			default:
				debugLog("Accept error: %v", err)
				continue
			}
		}

		debugLog("Client connected!")
		s.mu.Lock()
		s.conn = conn
		s.reader = bufio.NewReader(conn)
		s.connected = true
		s.mu.Unlock()

		// Start reading from this connection
		go s.readLoop()

		// Send game ready notification
		if err := s.SendGameReady(); err != nil {
			debugLog("Failed to send game_ready: %v", err)
		} else {
			debugLog("Sent game_ready notification")
		}
	}
}

// readLoop continuously reads messages from the connection.
func (s *SocketServer) readLoop() {
	debugLog("readLoop started")
	for {
		select {
		case <-s.done:
			debugLog("readLoop: done signal received")
			return
		default:
		}

		if s.reader == nil {
			debugLog("readLoop: reader is nil")
			return
		}

		line, err := s.reader.ReadBytes('\n')
		if err != nil {
			debugLog("readLoop: read error: %v", err)
			s.mu.Lock()
			s.connected = false
			s.mu.Unlock()
			return
		}

		if len(line) == 0 {
			continue
		}

		debugLog("readLoop: received line: %s", string(line))

		// Try to parse as a generic message first
		var msg map[string]interface{}
		if err := json.Unmarshal(line, &msg); err != nil {
			debugLog("readLoop: JSON parse error: %v", err)
			continue
		}

		// Determine message type
		if _, hasMethod := msg["method"]; hasMethod {
			// It's a request or notification
			if _, hasID := msg["id"]; hasID {
				// Request (has ID)
				var req Request
				if err := json.Unmarshal(line, &req); err == nil {
					debugLog("readLoop: parsed as Request, method=%s", req.Method)
					s.incoming <- &req
				}
			} else {
				// Notification (no ID)
				var notif Notification
				if err := json.Unmarshal(line, &notif); err == nil {
					debugLog("readLoop: parsed as Notification, method=%s", notif.Method)
					s.incoming <- &notif
				}
			}
		} else if _, hasResult := msg["result"]; hasResult {
			// It's a response
			var resp Response
			if err := json.Unmarshal(line, &resp); err == nil {
				s.incoming <- &resp
			}
		} else if _, hasError := msg["error"]; hasError {
			// It's an error response
			var resp Response
			if err := json.Unmarshal(line, &resp); err == nil {
				s.incoming <- &resp
			}
		}
	}
}

// processLoop handles incoming messages.
func (s *SocketServer) processLoop() {
	for {
		select {
		case <-s.done:
			return
		case msg := <-s.incoming:
			s.handleMessage(msg)
		}
	}
}

func (s *SocketServer) handleMessage(msg interface{}) {
	switch m := msg.(type) {
	case *Response:
		s.handleResponse(m)
	case *Request:
		s.handleRequest(m)
	case *Notification:
		s.handleNotification(m)
	}
}

func (s *SocketServer) handleResponse(resp *Response) {
	s.mu.Lock()
	ch, ok := s.pending[resp.ID]
	if ok {
		delete(s.pending, resp.ID)
	}
	s.mu.Unlock()

	if ok {
		ch <- resp
	}
}

func (s *SocketServer) handleRequest(req *Request) {
	var result interface{}
	var rpcErr *RPCError

	switch req.Method {
	case MethodChallengeComplete:
		if params, ok := req.Params.(map[string]interface{}); ok {
			cr := parseChallengeResult(params)
			if s.handler != nil {
				s.handler.HandleChallengeComplete(cr)
			}
			result = map[string]bool{"ok": true}
		}
	case MethodConfigUpdate:
		if params, ok := req.Params.(map[string]interface{}); ok {
			cfg := parseConfigUpdate(params)
			if s.handler != nil {
				s.handler.HandleConfigUpdate(cfg)
			}
			result = map[string]bool{"ok": true}
		}
	case MethodPauseGame:
		if s.handler != nil {
			s.handler.HandlePause()
		}
		result = map[string]bool{"ok": true}
	case MethodResumeGame:
		if s.handler != nil {
			s.handler.HandleResume()
		}
		result = map[string]bool{"ok": true}
	case MethodStartChallenge:
		if s.handler != nil {
			s.handler.HandleStartChallenge()
		}
		result = map[string]bool{"ok": true}
	case MethodRestartGame:
		if s.handler != nil {
			s.handler.HandleRestart()
		}
		result = map[string]bool{"ok": true}
	case MethodGoToLevelSelect:
		if s.handler != nil {
			s.handler.HandleGoToLevelSelect()
		}
		result = map[string]bool{"ok": true}
	default:
		rpcErr = NewError(ErrCodeMethodNotFound, "method not found: "+req.Method)
	}

	// Send response
	resp := NewResponse(req.ID, result, rpcErr)
	if err := s.send(resp); err != nil {
		debugLog("socket: failed to send response: %v", err)
	}
}

func (s *SocketServer) handleNotification(notif *Notification) {
	// Debug: log received notifications
	debugLog("handleNotification: method=%s", notif.Method)

	switch notif.Method {
	case MethodChallengeComplete:
		debugLog("Processing challenge_complete")
		if params, ok := notif.Params.(map[string]interface{}); ok {
			cr := parseChallengeResult(params)
			debugLog("Parsed result: request_id=%s success=%v", cr.RequestID, cr.Success)
			if s.handler != nil {
				s.handler.HandleChallengeComplete(cr)
				debugLog("Called handler.HandleChallengeComplete")
			}
		}
	case MethodConfigUpdate:
		if params, ok := notif.Params.(map[string]interface{}); ok {
			cfg := parseConfigUpdate(params)
			if s.handler != nil {
				s.handler.HandleConfigUpdate(cfg)
			}
		}
	case MethodPauseGame:
		if s.handler != nil {
			s.handler.HandlePause()
		}
	case MethodResumeGame:
		if s.handler != nil {
			s.handler.HandleResume()
		}
	case MethodStartChallenge:
		if s.handler != nil {
			s.handler.HandleStartChallenge()
		}
	case MethodRestartGame:
		if s.handler != nil {
			s.handler.HandleRestart()
		}
	case MethodGoToLevelSelect:
		if s.handler != nil {
			s.handler.HandleGoToLevelSelect()
		}
	}
}

// send writes a message to the connection.
func (s *SocketServer) send(msg interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn == nil || !s.connected {
		return errors.New("not connected")
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(s.conn, "%s\n", data)
	return err
}

// Notify sends a notification (no response expected).
func (s *SocketServer) Notify(method string, params interface{}) error {
	notif := NewNotification(method, params)
	return s.send(notif)
}

// RequestChallenge asks Neovim to present a challenge.
func (s *SocketServer) RequestChallenge(requestID string, challenge *ChallengeData) error {
	req := &ChallengeRequest{
		RequestID: requestID,
	}
	if challenge != nil {
		req.ChallengeID = challenge.ID
		req.ChallengeName = challenge.Name
		req.Category = challenge.Category
		req.Difficulty = challenge.Difficulty
		req.Description = challenge.Description
		req.InitialBuffer = challenge.InitialBuffer
		req.ExpectedBuffer = challenge.ExpectedBuffer
		req.ValidationType = challenge.ValidationType
		req.ExpectedCursor = challenge.ExpectedCursor
		req.ExpectedContent = challenge.ExpectedContent
		req.FunctionName = challenge.FunctionName
		req.CursorStart = challenge.CursorStart
		req.ParKeystrokes = challenge.ParKeystrokes
		req.GoldBase = challenge.GoldBase
		req.Filetype = challenge.Filetype
		req.HintAction = challenge.HintAction
		req.HintFallback = challenge.HintFallback
		req.Mode = challenge.Mode
		// Include feedback from previous challenge
		req.PrevSuccess = challenge.PrevSuccess
		req.PrevStreak = challenge.PrevStreak
		req.PrevGold = challenge.PrevGold
	}
	return s.Notify(MethodRequestChallenge, req)
}

// SendGameState sends current game state to Neovim.
func (s *SocketServer) SendGameState(state string, wave, gold, health, enemies, towers int) error {
	return s.Notify(MethodGameStateUpdate, &GameStateUpdate{
		State:   state,
		Wave:    wave,
		Gold:    gold,
		Health:  health,
		Enemies: enemies,
		Towers:  towers,
	})
}

// SendGameReady notifies Neovim that the game is ready.
func (s *SocketServer) SendGameReady() error {
	return s.Notify(MethodGameReady, nil)
}

// SendGoldUpdate notifies Neovim of gold changes.
func (s *SocketServer) SendGoldUpdate(gold, earned int, source string, speedBonus float64) error {
	return s.Notify(MethodGoldUpdate, &GoldUpdate{
		Gold:       gold,
		Earned:     earned,
		Source:     source,
		SpeedBonus: speedBonus,
	})
}

// SendChallengeAvailable notifies Neovim that challenges are available.
func (s *SocketServer) SendChallengeAvailable(count, nextReward int, nextCategory string) error {
	return s.Notify(MethodChallengeAvailable, &ChallengeAvailable{
		Count:        count,
		NextReward:   nextReward,
		NextCategory: nextCategory,
	})
}

// SendGameOver notifies Neovim that the game is over.
func (s *SocketServer) SendGameOver(wave, gold, towers, health int) error {
	return s.Notify(MethodGameOver, &GameOverParams{
		Wave:   wave,
		Gold:   gold,
		Towers: towers,
		Health: health,
	})
}

// SendVictory notifies Neovim that the player has won.
func (s *SocketServer) SendVictory(wave, gold, towers, health int) error {
	return s.Notify(MethodVictory, &VictoryParams{
		Wave:   wave,
		Gold:   gold,
		Towers: towers,
		Health: health,
	})
}
