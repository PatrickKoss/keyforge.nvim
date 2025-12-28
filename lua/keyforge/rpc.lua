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
M._job_id = nil
M._buffer = ""

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
  if not M._job_id then
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
  vim.fn.chansend(M._job_id, data)
end

--- Send a notification (no response expected)
---@param method string
---@param params? table
function M.notify(method, params)
  if not M._job_id then
    return
  end

  local msg = {
    jsonrpc = "2.0",
    method = method,
    params = params or {},
  }

  local data = encode(msg)
  vim.fn.chansend(M._job_id, data)
end

--- Send a response to an incoming request
---@param id number|string
---@param result? any
---@param error? table
local function respond(id, result, error)
  if not M._job_id then
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
  vim.fn.chansend(M._job_id, data)
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

--- Connect to an RPC channel (job ID)
---@param job_id number
function M.connect(job_id)
  M._job_id = job_id
  M._buffer = ""
  M._pending_requests = {}
end

--- Disconnect RPC
function M.disconnect()
  M._job_id = nil
  M._buffer = ""
  -- Cancel pending requests
  for id, callback in pairs(M._pending_requests) do
    callback({ message = "RPC disconnected" }, nil)
  end
  M._pending_requests = {}
end

--- Check if connected
---@return boolean
function M.is_connected()
  return M._job_id ~= nil
end

--- Handle stdout data from job
---@param data string[]
function M.on_stdout(data)
  if data then
    for _, line in ipairs(data) do
      handle_data(line .. "\n")
    end
  end
end

return M
