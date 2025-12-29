---@class RpcMessage
---@field jsonrpc string Always "2.0"
---@field method? string Method name for requests/notifications
---@field params? table Parameters
---@field id? number|string Request ID (absent for notifications)
---@field result? any Result for responses
---@field error? table Error object for error responses

local M = {}

-- RPC state
M._request_id = 0
M._pending_requests = {}
M._handlers = {}
M._socket = nil
M._buffer = ""
M._connected = false
M._socket_path = nil

--- Generate a unique request ID
---@return number
local function next_id()
  M._request_id = M._request_id + 1
  return M._request_id
end

--- Encode a message to JSON-RPC format
---@param msg table
---@return string
local function encode(msg)
  return vim.fn.json_encode(msg) .. "\n"
end

--- Decode a JSON-RPC message
---@param data string
---@return table|nil
local function decode(data)
  local ok, result = pcall(vim.fn.json_decode, data)
  if ok then
    return result
  end
  return nil
end

--- Register a handler for incoming RPC methods
---@param method string
---@param handler function
function M.on(method, handler)
  M._handlers[method] = handler
end

--- Send a request and wait for response
---@param method string
---@param params? table
---@param callback function Called with (error, result)
function M.request(method, params, callback)
  if not M._connected or not M._socket then
    callback({ message = "RPC not connected" }, nil)
    return
  end

  local id = next_id()
  local msg = {
    jsonrpc = "2.0",
    method = method,
    params = params or {},
    id = id,
  }

  M._pending_requests[id] = callback

  local data = encode(msg)
  M._socket:write(data)
end

--- Send a notification (no response expected)
---@param method string
---@param params? table
function M.notify(method, params)
  if not M._connected or not M._socket then
    return
  end

  local msg = {
    jsonrpc = "2.0",
    method = method,
    params = params or {},
  }

  local data = encode(msg)
  M._socket:write(data, function(err)
    if err then
      vim.schedule(function()
        vim.notify("RPC write error: " .. tostring(err), vim.log.levels.ERROR)
      end)
    end
  end)
end

--- Send a response to an incoming request
---@param id number|string
---@param result? any
---@param error? table
local function respond(id, result, error)
  if not M._connected or not M._socket then
    return
  end

  local msg = {
    jsonrpc = "2.0",
    id = id,
  }

  if error then
    msg.error = error
  else
    msg.result = result
  end

  local data = encode(msg)
  M._socket:write(data)
end

--- Handle incoming data
---@param data string
local function handle_data(data)
  -- Append to buffer
  M._buffer = M._buffer .. data

  -- Process complete messages (newline-delimited)
  while true do
    local newline = M._buffer:find("\n")
    if not newline then
      break
    end

    local line = M._buffer:sub(1, newline - 1)
    M._buffer = M._buffer:sub(newline + 1)

    if line ~= "" then
      local msg = decode(line)
      if msg then
        M._handle_message(msg)
      end
    end
  end
end

--- Handle a parsed message
---@param msg RpcMessage
function M._handle_message(msg)
  -- Response to our request
  if msg.id and (msg.result ~= nil or msg.error) then
    local callback = M._pending_requests[msg.id]
    if callback then
      M._pending_requests[msg.id] = nil
      callback(msg.error, msg.result)
    end
    return
  end

  -- Incoming request or notification
  if msg.method then
    local handler = M._handlers[msg.method]
    if handler then
      local ok, result = pcall(handler, msg.params)
      if msg.id then
        -- It's a request, send response
        if ok then
          respond(msg.id, result, nil)
        else
          respond(msg.id, nil, { code = -32000, message = tostring(result) })
        end
      end
    elseif msg.id then
      -- Unknown method, send error response
      respond(msg.id, nil, { code = -32601, message = "Method not found: " .. msg.method })
    end
  end
end

--- Connect to a Unix socket
---@param socket_path string Path to the Unix socket
---@param on_connect? function Called on successful connection
---@param on_error? function Called on connection error
function M.connect(socket_path, on_connect, on_error)
  if M._connected then
    M.disconnect()
  end

  M._socket_path = socket_path
  M._buffer = ""
  M._pending_requests = {}

  -- Create a new pipe (Unix socket client)
  M._socket = vim.loop.new_pipe(false)

  M._socket:connect(socket_path, function(err)
    vim.schedule(function()
      if err then
        M._connected = false
        if on_error then
          on_error(err)
        end
        return
      end

      M._connected = true

      -- Start reading from the socket
      M._socket:read_start(function(read_err, data)
        vim.schedule(function()
          if read_err then
            M.disconnect()
            return
          end

          if data then
            handle_data(data)
          else
            -- EOF - socket closed
            M.disconnect()
          end
        end)
      end)

      if on_connect then
        on_connect()
      end
    end)
  end)
end

--- Disconnect RPC
function M.disconnect()
  M._connected = false

  if M._socket then
    if not M._socket:is_closing() then
      M._socket:read_stop()
      M._socket:close()
    end
    M._socket = nil
  end

  M._buffer = ""
  M._socket_path = nil

  -- Cancel pending requests
  for _, callback in pairs(M._pending_requests) do
    callback({ message = "RPC disconnected" }, nil)
  end
  M._pending_requests = {}
end

--- Check if connected
---@return boolean
function M.is_connected()
  return M._connected and M._socket ~= nil
end

--- Get the socket path
---@return string|nil
function M.get_socket_path()
  return M._socket_path
end

-- Legacy compatibility functions (for existing code)

--- Connect to an RPC channel (job ID) - DEPRECATED
--- This is kept for backward compatibility but no longer used
---@param _job_id number (unused, deprecated)
---@deprecated Use M.connect(socket_path) instead
function M.connect_job(_job_id)
  -- No-op - socket connection replaces job-based RPC
  vim.notify("rpc.connect_job is deprecated, use socket connection", vim.log.levels.WARN)
end

--- Handle stdout data from job - DEPRECATED
---@param _data string[] (unused, deprecated)
---@deprecated Socket-based RPC handles data internally
function M.on_stdout(_data)
  -- No-op - socket connection handles data internally
end

-- Register default handlers for Keyforge-specific methods

--- Handler for gold_update notifications from the game
---@param params table GoldUpdate params
function M.handle_gold_update(params)
  local earned = params.earned or 0
  local source = params.source or "unknown"
  local speed_bonus = params.speed_bonus or 1.0

  -- Show notification to user
  local msg
  if source == "challenge" then
    msg = string.format("+%dg from challenge (%.1fx speed bonus)", earned, speed_bonus)
  elseif source == "mob" then
    msg = string.format("+%dg from defeating enemy", earned)
  elseif source == "wave_bonus" then
    msg = string.format("+%dg wave completion bonus!", earned)
  else
    msg = string.format("+%dg", earned)
  end

  vim.schedule(function()
    vim.notify(msg, vim.log.levels.INFO)
  end)
end

--- Handler for challenge_available notifications from the game
---@param params table ChallengeAvailable params
function M.handle_challenge_available(params)
  local count = params.count or 0
  local next_reward = params.next_reward or 0
  local next_category = params.next_category or "general"

  if count > 0 then
    vim.schedule(function()
      vim.notify(
        string.format("Challenges available: %d (next: ~%dg from %s)", count, next_reward, next_category),
        vim.log.levels.INFO
      )
    end)
  end
end

--- Handler for request_challenge requests from the game
---@param params table ChallengeRequest params
---@return table result
function M.handle_request_challenge(params)
  local challenge_buffer = require("keyforge.challenge_buffer")
  local keyforge = require("keyforge")

  vim.schedule(function()
    -- Store the challenge ID for the response
    keyforge._current_challenge_id = params.request_id
    keyforge._game_state = "challenge_waiting"

    -- Start the challenge using the new file-based buffer system
    -- Pass the full params which now includes challenge data from Go engine
    challenge_buffer.start_challenge(params.request_id, params)
  end)

  return { ok = true }
end

--- Handler for game_state_update notifications from the game
---@param params table GameStateUpdate params
function M.handle_game_state_update(params)
  local keyforge = require("keyforge")

  vim.schedule(function()
    local state = params.state or "playing"
    keyforge._game_state = state

    -- Handle game over or victory states
    if state == "game_over" or state == "victory" then
      local game_over = require("keyforge.game_over")
      game_over.show(params)
    end
  end)
end

--- Handler for game_over notification from the game
---@param params table GameOver params
function M.handle_game_over(params)
  local keyforge = require("keyforge")
  local game_over = require("keyforge.game_over")

  vim.schedule(function()
    keyforge._game_state = "game_over"
    game_over.show(vim.tbl_extend("force", params, { state = "game_over" }))
  end)
end

--- Handler for victory notification from the game
---@param params table Victory params
function M.handle_victory(params)
  local keyforge = require("keyforge")
  local game_over = require("keyforge.game_over")

  vim.schedule(function()
    keyforge._game_state = "victory"
    game_over.show(vim.tbl_extend("force", params, { state = "victory" }))
  end)
end

--- Handler for game_ready notification from the game
---@param _params table GameReady params (reserved for future use)
function M.handle_game_ready(_params)
  local keyforge = require("keyforge")

  vim.schedule(function()
    keyforge._game_state = "playing"
    vim.notify("Keyforge ready!", vim.log.levels.INFO)
  end)
end

--- Register the default handlers
function M.register_handlers()
  M.on("gold_update", M.handle_gold_update)
  M.on("challenge_available", M.handle_challenge_available)
  M.on("request_challenge", M.handle_request_challenge)
  M.on("game_state_update", M.handle_game_state_update)
  M.on("game_over", M.handle_game_over)
  M.on("victory", M.handle_victory)
  M.on("game_ready", M.handle_game_ready)
end

return M
