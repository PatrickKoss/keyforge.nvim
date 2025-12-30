--- Challenge queue system for Keyforge
--- Manages user-triggered challenges with gold rewards
local M = {}

local challenges = require("keyforge.challenges")
local keymap_hints = require("keyforge.keymap_hints")
local failure_feedback = require("keyforge.failure_feedback")

-- Queue state
M._state = {
  available = {},     -- challenges ready to be started
  current = nil,      -- active challenge (if any)
  completed = {},     -- session history
  skipped = {},       -- skipped challenges this session
  total_gold = 0,     -- total gold earned from challenges this session
}

-- Challenge buffer state
M._challenge_buf = nil
M._challenge_win = nil
M._initial_content = nil

-- Callbacks for RPC communication
M._on_challenge_start = nil
M._on_challenge_complete = nil

--- Calculate speed bonus multiplier
---@param time_ms number Time taken in milliseconds
---@param par_time_ms number Par time in milliseconds
---@return number multiplier Speed bonus (1.0 to 2.0)
function M.calculate_speed_bonus(time_ms, par_time_ms)
  if time_ms <= 0 or par_time_ms <= 0 then
    return 1.0
  end

  -- If completed slower than par, no bonus
  if time_ms >= par_time_ms then
    return 1.0
  end

  -- Calculate speed ratio
  local speed_ratio = par_time_ms / time_ms

  -- Apply formula: 1.0 + (speedRatio - 1.0) * 0.5
  local bonus = 1.0 + (speed_ratio - 1.0) * 0.5

  -- Cap at 2.0
  return math.min(2.0, bonus)
end

--- Calculate total gold for a challenge completion
---@param challenge table Challenge data
---@param efficiency number Efficiency score (0-1)
---@param speed_bonus number Speed bonus multiplier
---@return number gold Total gold earned
function M.calculate_gold(challenge, efficiency, speed_bonus)
  local base_gold = challenge.gold_base or 50
  local difficulty = challenge.difficulty or 1

  -- Difficulty multiplier: 1.0 + (difficulty * 0.25)
  local difficulty_mult = 1.0 + difficulty * 0.25

  -- Efficiency multiplier: 0.5 + (efficiency * 0.5)
  local efficiency_mult = 0.5 + efficiency * 0.5

  local gold = math.floor(base_gold * difficulty_mult * efficiency_mult * speed_bonus)
  return math.max(1, gold)
end

--- Load available challenges based on detected plugins
function M.load_available_challenges()
  M._state.available = {}

  -- Get detected plugins
  local plugins = keymap_hints.get_detected_plugins()

  -- Load all sample challenges for now
  -- In the future, this will load from YAML and filter by required_plugin
  for _, challenge in ipairs(challenges.sample_challenges) do
    -- Check if challenge requires a plugin
    local required_plugin = challenge.required_plugin
    if required_plugin then
      if plugins[required_plugin] then
        table.insert(M._state.available, vim.deepcopy(challenge))
      end
    else
      -- No plugin requirement, always available
      table.insert(M._state.available, vim.deepcopy(challenge))
    end
  end

  -- Shuffle the available challenges
  for i = #M._state.available, 2, -1 do
    local j = math.random(i)
    M._state.available[i], M._state.available[j] = M._state.available[j], M._state.available[i]
  end
end

--- Get the next available challenge
---@return table|nil challenge
function M.get_next_challenge()
  if #M._state.available == 0 then
    -- Reload challenges if empty
    M.load_available_challenges()
  end

  if #M._state.available == 0 then
    return nil
  end

  return M._state.available[1]
end

--- Enrich a challenge with keymap hints
---@param challenge table
---@return table enriched_challenge
function M.get_challenge_with_hints(challenge)
  local enriched = vim.deepcopy(challenge)

  -- Get hints for the challenge category
  local hints = keymap_hints.get_hints_for_category(challenge.category)
  enriched.hints = hints

  -- Format hints for display
  enriched.hints_display = keymap_hints.format_hints_for_display(hints)

  return enriched
end

--- Check if a challenge is currently active
---@return boolean
function M.is_challenge_active()
  return M._state.current ~= nil
end

--- Request the next challenge (user-triggered or game-triggered via RPC)
---@param _category? string Optional category hint from the game (reserved for future use)
---@return table|nil challenge The started challenge, or nil if none available
function M.request_next(_category)
  -- Don't start a new challenge if one is active
  if M.is_challenge_active() then
    vim.notify("A challenge is already active. Complete it first!", vim.log.levels.WARN)
    return nil
  end

  local challenge = M.get_next_challenge()
  if not challenge then
    vim.notify("No challenges available!", vim.log.levels.WARN)
    return nil
  end

  -- Remove from available queue
  table.remove(M._state.available, 1)

  -- Enrich with hints
  local enriched = M.get_challenge_with_hints(challenge)

  -- Set as current
  M._state.current = enriched

  -- Create the challenge buffer
  M._create_challenge_buffer(enriched)

  -- Start tracking keystrokes
  challenges.start_tracking()

  -- Notify callback
  if M._on_challenge_start then
    M._on_challenge_start(enriched)
  end

  return enriched
end

--- Create the challenge buffer with content and hints
---@param challenge table
function M._create_challenge_buffer(challenge)
  -- Create a new buffer for the challenge
  M._challenge_buf = vim.api.nvim_create_buf(false, true)

  -- Set buffer options
  vim.api.nvim_buf_set_option(M._challenge_buf, "buftype", "nofile")
  vim.api.nvim_buf_set_option(M._challenge_buf, "bufhidden", "wipe")
  vim.api.nvim_buf_set_option(M._challenge_buf, "swapfile", false)

  -- Set filetype for syntax highlighting
  local filetype = challenge.filetype or "text"
  vim.api.nvim_buf_set_option(M._challenge_buf, "filetype", filetype)

  -- Set the initial content
  local initial = challenge.initial_buffer or ""
  local lines = vim.split(initial, "\n")
  vim.api.nvim_buf_set_lines(M._challenge_buf, 0, -1, false, lines)

  -- Store initial content for validation
  M._initial_content = vim.deepcopy(lines)

  -- Create a floating window for the challenge
  local width = math.min(80, vim.o.columns - 4)
  local height = math.min(20, vim.o.lines - 10)

  local row = math.floor((vim.o.lines - height) / 2)
  local col = math.floor((vim.o.columns - width) / 2)

  M._challenge_win = vim.api.nvim_open_win(M._challenge_buf, true, {
    relative = "editor",
    width = width,
    height = height,
    row = row,
    col = col,
    style = "minimal",
    border = "rounded",
    title = string.format(" %s (Difficulty: %d) ", challenge.name, challenge.difficulty or 1),
    title_pos = "center",
  })

  -- Set window options
  vim.api.nvim_win_set_option(M._challenge_win, "wrap", false)
  vim.api.nvim_win_set_option(M._challenge_win, "cursorline", true)

  -- Set cursor to starting position if specified
  if challenge.cursor_start then
    local row_pos = (challenge.cursor_start[1] or 0) + 1
    local col_pos = challenge.cursor_start[2] or 0
    pcall(vim.api.nvim_win_set_cursor, M._challenge_win, { row_pos, col_pos })
  end

  -- Show challenge info in the command line
  local info_lines = {
    string.format("Challenge: %s", challenge.name),
    string.format("Goal: %s", challenge.description),
    string.format("Par: %d keystrokes | Base reward: %dg", challenge.par_keystrokes or 10, challenge.gold_base or 50),
  }

  -- Add hints if available
  if challenge.hints_display and #challenge.hints_display > 0 then
    table.insert(info_lines, "")
    for _, hint_line in ipairs(challenge.hints_display) do
      table.insert(info_lines, hint_line)
    end
  end

  -- Display info
  vim.notify(table.concat(info_lines, "\n"), vim.log.levels.INFO)
end

--- Complete the current challenge
---@return table|nil result The completion result
function M.complete_current()
  if not M.is_challenge_active() then
    vim.notify("No active challenge to complete!", vim.log.levels.WARN)
    return nil
  end

  local challenge = M._state.current

  -- Get final buffer content
  local final_content = {}
  if M._challenge_buf and vim.api.nvim_buf_is_valid(M._challenge_buf) then
    final_content = vim.api.nvim_buf_get_lines(M._challenge_buf, 0, -1, false)
  end

  -- Validate the challenge
  local result = challenges.validate(challenge, M._initial_content, final_content)

  -- Calculate speed bonus and gold if successful
  if result.success then
    local par_time_ms = (challenge.par_keystrokes or 10) * 1000 -- rough estimate
    result.speed_bonus = M.calculate_speed_bonus(result.time_ms, par_time_ms)
    result.gold_earned = M.calculate_gold(challenge, result.efficiency, result.speed_bonus)
    M._state.total_gold = M._state.total_gold + result.gold_earned

    -- Add to completed list
    table.insert(M._state.completed, {
      challenge = challenge,
      result = result,
    })

    vim.notify(string.format(
      "Challenge complete! +%dg (efficiency: %.0f%%, speed bonus: %.1fx)",
      result.gold_earned,
      result.efficiency * 100,
      result.speed_bonus
    ), vim.log.levels.INFO)

    -- Clean up challenge buffer
    M._cleanup_challenge_buffer()

    -- Clear current challenge
    M._state.current = nil

    -- Notify callback
    if M._on_challenge_complete then
      M._on_challenge_complete(result)
    end

    return result
  else
    result.speed_bonus = 1.0
    result.gold_earned = 0

    -- Show failure feedback instead of just a notification
    failure_feedback.show(challenge, result.failure_details, {
      on_retry = function()
        M._retry_current()
      end,
      on_skip = function()
        M.skip_current()
      end,
    })

    return nil -- Don't return result yet, waiting for user action
  end
end

--- Retry the current challenge (reset buffer to initial state)
function M._retry_current()
  if not M.is_challenge_active() then
    return
  end

  -- Reset buffer content to initial
  if M._challenge_buf and vim.api.nvim_buf_is_valid(M._challenge_buf) then
    local initial = M._state.current.initial_buffer or ""
    local lines = vim.split(initial, "\n")
    vim.api.nvim_buf_set_option(M._challenge_buf, "modifiable", true)
    vim.api.nvim_buf_set_lines(M._challenge_buf, 0, -1, false, lines)
    M._initial_content = vim.deepcopy(lines)

    -- Reset cursor if specified
    local challenge = M._state.current
    if challenge.cursor_start then
      local row_pos = (challenge.cursor_start[1] or 0) + 1
      local col_pos = challenge.cursor_start[2] or 0
      pcall(vim.api.nvim_win_set_cursor, M._challenge_win, { row_pos, col_pos })
    end
  end

  -- Restart keystroke tracking
  challenges.start_tracking()

  vim.notify("Challenge reset. Try again!", vim.log.levels.INFO)
end

--- Skip the current challenge
---@return table|nil result Skip result
function M.skip_current()
  if not M.is_challenge_active() then
    vim.notify("No active challenge to skip!", vim.log.levels.WARN)
    return nil
  end

  local challenge = M._state.current

  -- Stop tracking
  challenges.stop_tracking()

  -- Record as skipped
  table.insert(M._state.skipped, challenge)

  local result = {
    success = false,
    skipped = true,
    keystroke_count = 0,
    time_ms = 0,
    efficiency = 0,
    speed_bonus = 1.0,
    gold_earned = 0,
  }

  -- Clean up
  M._cleanup_challenge_buffer()
  M._state.current = nil

  vim.notify("Challenge skipped.", vim.log.levels.INFO)

  -- Notify callback
  if M._on_challenge_complete then
    M._on_challenge_complete(result)
  end

  return result
end

--- Clean up the challenge buffer and window
function M._cleanup_challenge_buffer()
  if M._challenge_win and vim.api.nvim_win_is_valid(M._challenge_win) then
    vim.api.nvim_win_close(M._challenge_win, true)
  end
  M._challenge_win = nil

  if M._challenge_buf and vim.api.nvim_buf_is_valid(M._challenge_buf) then
    vim.api.nvim_buf_delete(M._challenge_buf, { force = true })
  end
  M._challenge_buf = nil
  M._initial_content = nil
end

--- Get session statistics
---@return table stats
function M.get_stats()
  return {
    total_gold = M._state.total_gold,
    completed_count = #M._state.completed,
    skipped_count = #M._state.skipped,
    available_count = #M._state.available,
    current = M._state.current,
  }
end

--- Reset session state
function M.reset()
  M._cleanup_challenge_buffer()
  M._state = {
    available = {},
    current = nil,
    completed = {},
    skipped = {},
    total_gold = 0,
  }
end

--- Set callback for challenge start
---@param callback function(challenge: table)
function M.on_challenge_start(callback)
  M._on_challenge_start = callback
end

--- Set callback for challenge complete
---@param callback function(result: table)
function M.on_challenge_complete(callback)
  M._on_challenge_complete = callback
end

--- Initialize the queue (call on game start)
function M.init()
  M.reset()
  keymap_hints.discover_keymaps()
  M.load_available_challenges()
end

return M
