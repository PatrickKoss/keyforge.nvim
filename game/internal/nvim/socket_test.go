package nvim

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// MockHandler implements Handler for testing.
type MockHandler struct {
	ChallengeResults []*ChallengeResult
	ConfigUpdates    []*ConfigUpdate
	PauseCalls       int
	ResumeCalls      int
	StartChallenges  int
	RestartCalls     int
	LevelSelectCalls int
}

func (h *MockHandler) HandleChallengeComplete(result *ChallengeResult) {
	h.ChallengeResults = append(h.ChallengeResults, result)
}

func (h *MockHandler) HandleConfigUpdate(config *ConfigUpdate) {
	h.ConfigUpdates = append(h.ConfigUpdates, config)
}

func (h *MockHandler) HandlePause() {
	h.PauseCalls++
}

func (h *MockHandler) HandleResume() {
	h.ResumeCalls++
}

func (h *MockHandler) HandleStartChallenge() {
	h.StartChallenges++
}

func (h *MockHandler) HandleRestart() {
	h.RestartCalls++
}

func (h *MockHandler) HandleGoToLevelSelect() {
	h.LevelSelectCalls++
}

func TestSocketServerStartStop(t *testing.T) {
	tmpDir := os.TempDir()
	socketPath := filepath.Join(tmpDir, "test_socket.sock")
	defer os.Remove(socketPath)

	handler := &MockHandler{}
	server := NewSocketServer(socketPath, handler)

	err := server.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Give it a moment
	time.Sleep(50 * time.Millisecond)

	server.Stop()
}

func TestSocketServerReceivesChallengeComplete(t *testing.T) {
	tmpDir := os.TempDir()
	socketPath := filepath.Join(tmpDir, "test_challenge.sock")
	defer os.Remove(socketPath)

	handler := &MockHandler{}
	server := NewSocketServer(socketPath, handler)

	err := server.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer server.Stop()

	// Connect as a client
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Wait for connection to establish
	time.Sleep(100 * time.Millisecond)

	// Send a challenge_complete notification
	notif := Notification{
		JSONRPC: "2.0",
		Method:  "challenge_complete",
		Params: map[string]interface{}{
			"request_id":      "test_challenge_1",
			"success":         true,
			"keystroke_count": 5.0,
			"time_ms":         1000.0,
			"efficiency":      1.5,
			"gold_earned":     75.0,
		},
	}

	data, _ := json.Marshal(notif)
	conn.Write(append(data, '\n'))

	// Wait for processing
	time.Sleep(200 * time.Millisecond)

	// Verify handler was called
	if len(handler.ChallengeResults) != 1 {
		t.Fatalf("Expected 1 challenge result, got %d", len(handler.ChallengeResults))
	}

	result := handler.ChallengeResults[0]
	if result.RequestID != "test_challenge_1" {
		t.Errorf("Expected request_id 'test_challenge_1', got '%s'", result.RequestID)
	}
	if !result.Success {
		t.Error("Expected success=true")
	}
	if result.GoldEarned != 75 {
		t.Errorf("Expected gold_earned=75, got %d", result.GoldEarned)
	}
}

func TestSocketServerSendsRequestChallenge(t *testing.T) {
	tmpDir := os.TempDir()
	socketPath := filepath.Join(tmpDir, "test_request.sock")
	defer os.Remove(socketPath)

	handler := &MockHandler{}
	server := NewSocketServer(socketPath, handler)

	err := server.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer server.Stop()

	// Connect as a client
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Wait for connection
	time.Sleep(100 * time.Millisecond)

	// Clear any initial messages (like game_ready)
	conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	buf := make([]byte, 4096)
	conn.Read(buf) // Discard game_ready notification

	// Now send a request_challenge from the server
	challenge := &ChallengeData{
		ID:       "test_challenge",
		Name:     "Test Challenge",
		Category: "movement",
	}
	err = server.RequestChallenge("req_123", challenge)
	if err != nil {
		t.Fatalf("Failed to request challenge: %v", err)
	}

	// Read the notification sent to us
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}

	var received map[string]interface{}
	err = json.Unmarshal(buf[:n-1], &received) // -1 to remove newline
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v (data: %s)", err, string(buf[:n]))
	}

	if received["method"] != "request_challenge" {
		t.Errorf("Expected method 'request_challenge', got '%v'", received["method"])
	}

	params := received["params"].(map[string]interface{})
	if params["request_id"] != "req_123" {
		t.Errorf("Expected request_id 'req_123', got '%v'", params["request_id"])
	}
}

// TestEndToEndChallengeFlow tests the full flow:
// 1. Server sends request_challenge to client
// 2. Client sends challenge_complete back
// 3. Handler receives the result.
func TestEndToEndChallengeFlow(t *testing.T) {
	tmpDir := os.TempDir()
	socketPath := filepath.Join(tmpDir, "test_e2e.sock")
	defer os.Remove(socketPath)

	handler := &MockHandler{}
	server := NewSocketServer(socketPath, handler)

	err := server.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer server.Stop()

	// Connect as a client (simulating Neovim)
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	time.Sleep(100 * time.Millisecond)

	// Step 1: Server requests a challenge
	challenge := &ChallengeData{
		ID:       "e2e_test_challenge",
		Name:     "E2E Test Challenge",
		Category: "lsp",
	}
	err = server.RequestChallenge("e2e_challenge_1", challenge)
	if err != nil {
		t.Fatalf("Failed to request challenge: %v", err)
	}

	// Step 2: Client receives it and sends back result
	buf := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	conn.Read(buf) // game_ready
	conn.Read(buf) // request_challenge

	// Send challenge_complete
	response := Notification{
		JSONRPC: "2.0",
		Method:  "challenge_complete",
		Params: map[string]interface{}{
			"request_id":  "e2e_challenge_1",
			"success":     true,
			"gold_earned": 50.0,
		},
	}
	data, _ := json.Marshal(response)
	conn.Write(append(data, '\n'))

	// Wait for processing
	time.Sleep(200 * time.Millisecond)

	// Step 3: Verify handler received it
	if len(handler.ChallengeResults) != 1 {
		t.Fatalf("Expected 1 challenge result, got %d", len(handler.ChallengeResults))
	}

	result := handler.ChallengeResults[0]
	if result.RequestID != "e2e_challenge_1" {
		t.Errorf("Expected request_id 'e2e_challenge_1', got '%s'", result.RequestID)
	}
	if result.GoldEarned != 50 {
		t.Errorf("Expected gold_earned=50, got %d", result.GoldEarned)
	}
}

// =============================================================================
// Game Over / Victory RPC Tests
// =============================================================================

// TestSocketServerHandlesRestartNotification tests that restart_game RPC is handled.
func TestSocketServerHandlesRestartNotification(t *testing.T) {
	tmpDir := os.TempDir()
	socketPath := filepath.Join(tmpDir, "test_restart.sock")
	defer os.Remove(socketPath)

	handler := &MockHandler{}
	server := NewSocketServer(socketPath, handler)

	err := server.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer server.Stop()

	// Connect as client
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	time.Sleep(100 * time.Millisecond)

	// Send restart_game notification
	notif := Notification{
		JSONRPC: "2.0",
		Method:  "restart_game",
		Params:  map[string]interface{}{},
	}
	data, _ := json.Marshal(notif)
	conn.Write(append(data, '\n'))

	time.Sleep(200 * time.Millisecond)

	if handler.RestartCalls != 1 {
		t.Errorf("Expected 1 restart call, got %d", handler.RestartCalls)
	}
}

// TestSocketServerHandlesGoToLevelSelectNotification tests the new RPC method.
func TestSocketServerHandlesGoToLevelSelectNotification(t *testing.T) {
	tmpDir := os.TempDir()
	socketPath := filepath.Join(tmpDir, "test_levelselect.sock")
	defer os.Remove(socketPath)

	handler := &MockHandler{}
	server := NewSocketServer(socketPath, handler)

	err := server.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer server.Stop()

	// Connect as client
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	time.Sleep(100 * time.Millisecond)

	// Send go_to_level_select notification
	notif := Notification{
		JSONRPC: "2.0",
		Method:  "go_to_level_select",
		Params:  map[string]interface{}{},
	}
	data, _ := json.Marshal(notif)
	conn.Write(append(data, '\n'))

	time.Sleep(200 * time.Millisecond)

	if handler.LevelSelectCalls != 1 {
		t.Errorf("Expected 1 level select call, got %d", handler.LevelSelectCalls)
	}
}

// TestSocketServerHandlesRestartRequest tests restart_game as JSON-RPC request.
func TestSocketServerHandlesRestartRequest(t *testing.T) {
	tmpDir := os.TempDir()
	socketPath := filepath.Join(tmpDir, "test_restart_req.sock")
	defer os.Remove(socketPath)

	handler := &MockHandler{}
	server := NewSocketServer(socketPath, handler)

	err := server.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer server.Stop()

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	time.Sleep(100 * time.Millisecond)

	// Consume game_ready
	buf := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	conn.Read(buf)

	// Send restart_game as request (with ID)
	req := Request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "restart_game",
		Params:  map[string]interface{}{},
	}
	data, _ := json.Marshal(req)
	conn.Write(append(data, '\n'))

	// Read response
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	var resp Response
	err = json.Unmarshal(buf[:n-1], &resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Error != nil {
		t.Errorf("Expected no error, got: %v", resp.Error)
	}

	if handler.RestartCalls != 1 {
		t.Errorf("Expected 1 restart call, got %d", handler.RestartCalls)
	}
}

// TestSocketServerHandlesGoToLevelSelectRequest tests go_to_level_select as request.
func TestSocketServerHandlesGoToLevelSelectRequest(t *testing.T) {
	tmpDir := os.TempDir()
	socketPath := filepath.Join(tmpDir, "test_levelselect_req.sock")
	defer os.Remove(socketPath)

	handler := &MockHandler{}
	server := NewSocketServer(socketPath, handler)

	err := server.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer server.Stop()

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	time.Sleep(100 * time.Millisecond)

	// Consume game_ready
	buf := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	conn.Read(buf)

	// Send go_to_level_select as request
	req := Request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "go_to_level_select",
		Params:  map[string]interface{}{},
	}
	data, _ := json.Marshal(req)
	conn.Write(append(data, '\n'))

	// Read response
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	var resp Response
	err = json.Unmarshal(buf[:n-1], &resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Error != nil {
		t.Errorf("Expected no error, got: %v", resp.Error)
	}

	if handler.LevelSelectCalls != 1 {
		t.Errorf("Expected 1 level select call, got %d", handler.LevelSelectCalls)
	}
}

// TestEndToEndRestartFlow tests full restart flow via RPC.
func TestEndToEndRestartFlow(t *testing.T) {
	tmpDir := os.TempDir()
	socketPath := filepath.Join(tmpDir, "test_e2e_restart.sock")
	defer os.Remove(socketPath)

	handler := &MockHandler{}
	server := NewSocketServer(socketPath, handler)

	err := server.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer server.Stop()

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	time.Sleep(100 * time.Millisecond)

	// Simulate game over → user presses 'r' → Lua sends restart_game
	notif := Notification{
		JSONRPC: "2.0",
		Method:  "restart_game",
		Params:  map[string]interface{}{},
	}
	data, _ := json.Marshal(notif)
	conn.Write(append(data, '\n'))

	time.Sleep(200 * time.Millisecond)

	// Handler should have received the restart command
	if handler.RestartCalls != 1 {
		t.Errorf("Expected HandleRestart to be called once, got %d", handler.RestartCalls)
	}

	// Level select should not have been called
	if handler.LevelSelectCalls != 0 {
		t.Errorf("Expected no level select calls, got %d", handler.LevelSelectCalls)
	}
}

// TestEndToEndLevelSelectFlow tests full level select flow via RPC.
func TestEndToEndLevelSelectFlow(t *testing.T) {
	tmpDir := os.TempDir()
	socketPath := filepath.Join(tmpDir, "test_e2e_levelselect.sock")
	defer os.Remove(socketPath)

	handler := &MockHandler{}
	server := NewSocketServer(socketPath, handler)

	err := server.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer server.Stop()

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	time.Sleep(100 * time.Millisecond)

	// Simulate game over → user presses 'l' → Lua sends go_to_level_select
	notif := Notification{
		JSONRPC: "2.0",
		Method:  "go_to_level_select",
		Params:  map[string]interface{}{},
	}
	data, _ := json.Marshal(notif)
	conn.Write(append(data, '\n'))

	time.Sleep(200 * time.Millisecond)

	// Handler should have received the level select command
	if handler.LevelSelectCalls != 1 {
		t.Errorf("Expected HandleGoToLevelSelect to be called once, got %d", handler.LevelSelectCalls)
	}

	// Restart should not have been called
	if handler.RestartCalls != 0 {
		t.Errorf("Expected no restart calls, got %d", handler.RestartCalls)
	}
}
