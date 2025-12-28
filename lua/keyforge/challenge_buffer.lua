--- File-based challenge buffer system for Keyforge
--- Creates real temp files for challenges so LSP and vim features work properly
local M = {}

local challenges = require("keyforge.challenges")

-- Challenge state
M._current = nil -- Current challenge data
M._request_id = nil -- RPC request ID
M._filepath = nil -- Path to temp challenge file
M._challenge_buf = nil -- Buffer number
M._challenge_win = nil -- Window number
M._challenge_tab = nil -- Tab page number
M._initial_content = nil -- Initial file content for validation
M._info_win = nil -- Floating window for challenge info
M._info_buf = nil -- Buffer for challenge info
M._timeout_timer = nil -- Timeout timer handle

-- File extension mapping by filetype
local filetype_extensions = {
  lua = ".lua",
  go = ".go",
  python = ".py",
  javascript = ".js",
  typescript = ".ts",
  rust = ".rs",
  c = ".c",
  cpp = ".cpp",
  java = ".java",
  ruby = ".rb",
  text = ".txt",
}

--- Get the temp directory for challenge files
---@return string
local function get_temp_dir()
  local tmpdir = vim.fn.tempname()
  -- tempname returns a file path, we want the directory
  return vim.fn.fnamemodify(tmpdir, ":h")
end

--- Generate a unique temp file path for a challenge
---@param challenge table Challenge data
---@return string filepath
local function generate_temp_filepath(challenge)
  local filetype = challenge.filetype or "text"
  local ext = filetype_extensions[filetype] or ".txt"
  local timestamp = os.time()
  local random = math.random(1000, 9999)
  local filename = string.format("keyforge_challenge_%d_%d%s", timestamp, random, ext)
  return get_temp_dir() .. "/" .. filename
end

--- Create the temp file with challenge content
---@param challenge table Challenge data
---@return string filepath Path to created file
local function create_temp_file(challenge)
  local filepath = generate_temp_filepath(challenge)
  local content = challenge.initial_buffer or ""

  -- Write content to file
  local file = io.open(filepath, "w")
  if file then
    file:write(content)
    file:close()
  end

  return filepath
end

--- Delete the temp file
---@param filepath string
local function delete_temp_file(filepath)
  if filepath and vim.fn.filereadable(filepath) == 1 then
    vim.fn.delete(filepath)
  end
end

--- Create the challenge info floating window
---@param challenge table
local function create_info_window(challenge)
  local keyforge = require("keyforge")
  local config = keyforge.config

  -- Build info content
  local lines = {
    "Challenge: " .. (challenge.name or "Unknown"),
    "",
    "Goal: " .. (challenge.description or "Complete the challenge"),
    "",
    string.format("Par: %d keystrokes | Reward: %dg", challenge.par_keystrokes or 10, challenge.gold_base or 50),
    "",
    string.format("Submit: %s (normal mode) | Cancel: %s", config.keybind_submit, config.keybind_cancel),
  }

  -- Calculate window size
  local width = 0
  for _, line in ipairs(lines) do
    width = math.max(width, #line)
  end
  width = math.min(width + 4, vim.o.columns - 4)
  local height = #lines

  -- Create buffer
  M._info_buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_buf_set_lines(M._info_buf, 0, -1, false, lines)
  vim.api.nvim_buf_set_option(M._info_buf, "modifiable", false)
  vim.api.nvim_buf_set_option(M._info_buf, "buftype", "nofile")

  -- Position at top right
  local row = 1
  local col = vim.o.columns - width - 2

  -- Create floating window
  M._info_win = vim.api.nvim_open_win(M._info_buf, false, {
    relative = "editor",
    width = width,
    height = height,
    row = row,
    col = col,
    style = "minimal",
    border = "rounded",
    title = " Keyforge Challenge ",
    title_pos = "center",
  })

  -- Set window options
  vim.api.nvim_win_set_option(M._info_win, "winblend", 10)
end

--- Close the info window
local function close_info_window()
  if M._info_win and vim.api.nvim_win_is_valid(M._info_win) then
    vim.api.nvim_win_close(M._info_win, true)
  end
  M._info_win = nil

  if M._info_buf and vim.api.nvim_buf_is_valid(M._info_buf) then
    vim.api.nvim_buf_delete(M._info_buf, { force = true })
  end
  M._info_buf = nil
end

--- Set up keymaps for the challenge buffer
---@param buf number Buffer number
local function setup_keymaps(buf)
  local keyforge = require("keyforge")
  local config = keyforge.config

  -- Submit keymap (default <CR>)
  vim.keymap.set("n", config.keybind_submit, function()
    M.submit_challenge()
  end, { buffer = buf, desc = "Submit challenge" })

  -- Cancel keymap (default <Esc>)
  -- Note: We use a different approach for Escape to avoid conflicts
  vim.keymap.set("n", config.keybind_cancel, function()
    M.cancel_challenge()
  end, { buffer = buf, desc = "Cancel challenge" })
end

--- Set up autocmds for the challenge buffer
---@param buf number Buffer number
local function setup_autocmds(buf)
  local group = vim.api.nvim_create_augroup("KeyforgeChallenge", { clear = true })

  -- Handle buffer close (tab close, :q, etc.)
  vim.api.nvim_create_autocmd("BufWipeout", {
    group = group,
    buffer = buf,
    callback = function()
      -- If challenge is still active, treat as cancel
      if M._current then
        M._complete_challenge(false, true) -- skipped
      end
    end,
    once = true,
  })

  -- Track when user leaves the buffer (optional warning)
  vim.api.nvim_create_autocmd("BufLeave", {
    group = group,
    buffer = buf,
    callback = function()
      -- Just update info window position if needed
    end,
  })
end

--- Start a new challenge
---@param request_id string RPC request ID from game
---@param category string Challenge category
---@param difficulty number Difficulty level (1-3)
function M.start_challenge(request_id, category, difficulty)
  local keyforge = require("keyforge")

  -- Don't start if one is already active
  if M._current then
    vim.notify("A challenge is already active!", vim.log.levels.WARN)
    return
  end

  -- Get a matching challenge
  local challenge = challenges.get_random_challenge(category, difficulty)
  if not challenge then
    -- Fallback to any challenge
    challenge = challenges.get_random_challenge(nil, nil)
  end

  if not challenge then
    vim.notify("No challenges available!", vim.log.levels.ERROR)
    -- Send skip result back to game
    M._send_result(request_id, {
      success = false,
      skipped = true,
      keystroke_count = 0,
      time_ms = 0,
      efficiency = 0,
      gold_earned = 0,
    })
    return
  end

  -- Store state
  M._current = challenge
  M._request_id = request_id
  M._initial_content = vim.split(challenge.initial_buffer or "", "\n")

  -- Create temp file
  M._filepath = create_temp_file(challenge)

  -- Open in new tab
  vim.cmd("tabnew " .. vim.fn.fnameescape(M._filepath))
  M._challenge_tab = vim.api.nvim_get_current_tabpage()
  M._challenge_buf = vim.api.nvim_get_current_buf()
  M._challenge_win = vim.api.nvim_get_current_win()

  -- Set filetype for syntax highlighting and LSP
  local filetype = challenge.filetype or "text"
  vim.bo[M._challenge_buf].filetype = filetype

  -- Set cursor position if specified
  if challenge.cursor_start and #challenge.cursor_start == 2 then
    local row = (challenge.cursor_start[1] or 0) + 1
    local col = challenge.cursor_start[2] or 0
    pcall(vim.api.nvim_win_set_cursor, M._challenge_win, { row, col })
  end

  -- Set up keymaps and autocmds
  setup_keymaps(M._challenge_buf)
  setup_autocmds(M._challenge_buf)

  -- Start keystroke tracking
  challenges.start_tracking()

  -- Create info window
  create_info_window(challenge)

  -- Set up timeout using vim.fn.timer_start (returns timer ID)
  local timeout_seconds = keyforge.config.challenge_timeout or 300
  M._timeout_timer = vim.fn.timer_start(timeout_seconds * 1000, function()
    if M._current then
      vim.notify("Challenge timed out!", vim.log.levels.WARN)
      M.cancel_challenge()
    end
  end)

  vim.notify(string.format("Challenge started: %s", challenge.name), vim.log.levels.INFO)
end

--- Submit the current challenge
function M.submit_challenge()
  if not M._current then
    vim.notify("No active challenge!", vim.log.levels.WARN)
    return
  end

  -- Read final content from file
  local final_content = {}
  if M._filepath and vim.fn.filereadable(M._filepath) == 1 then
    -- Save the buffer first to ensure file is up to date
    if M._challenge_buf and vim.api.nvim_buf_is_valid(M._challenge_buf) then
      vim.api.nvim_buf_call(M._challenge_buf, function()
        vim.cmd("silent! write")
      end)
    end

    local file = io.open(M._filepath, "r")
    if file then
      local content = file:read("*all")
      file:close()
      final_content = vim.split(content, "\n")
      -- Remove trailing empty line if present (file:read adds one)
      if #final_content > 0 and final_content[#final_content] == "" then
        table.remove(final_content)
      end
    end
  end

  -- Validate
  local result = challenges.validate(M._current, M._initial_content, final_content)

  if result.success then
    -- Calculate gold
    local gold = challenges.calculate_reward(M._current, result.efficiency)
    result.gold_earned = gold
    vim.notify(
      string.format("Challenge complete! +%dg (efficiency: %.0f%%)", gold, result.efficiency * 100),
      vim.log.levels.INFO
    )
  else
    result.gold_earned = 0
    vim.notify("Challenge failed! Buffer content doesn't match expected result.", vim.log.levels.WARN)
  end

  M._complete_challenge(result.success, false, result)
end

--- Cancel the current challenge
function M.cancel_challenge()
  if not M._current then
    return
  end

  -- Stop tracking and get stats
  local keystrokes, time_ms = challenges.stop_tracking()

  vim.notify("Challenge cancelled.", vim.log.levels.INFO)

  M._complete_challenge(false, true, {
    keystroke_count = keystrokes,
    time_ms = time_ms,
    efficiency = 0,
    gold_earned = 0,
  })
end

--- Internal: Complete the challenge and clean up
---@param success boolean Whether challenge was successful
---@param skipped boolean Whether challenge was skipped/cancelled
---@param result? table Optional result data
function M._complete_challenge(success, skipped, result)
  result = result or {}

  -- Cancel timeout timer
  if M._timeout_timer then
    vim.fn.timer_stop(M._timeout_timer)
    M._timeout_timer = nil
  end

  -- Send result back to game
  if M._request_id then
    M._send_result(M._request_id, {
      request_id = M._request_id,
      success = success,
      skipped = skipped,
      keystroke_count = result.keystroke_count or 0,
      time_ms = result.time_ms or 0,
      efficiency = result.efficiency or 0,
      gold_earned = result.gold_earned or 0,
    })
  end

  -- Clean up
  M._cleanup()
end

--- Send result back to game via RPC
---@param request_id string
---@param result table
function M._send_result(request_id, result)
  local rpc = require("keyforge.rpc")
  result.request_id = request_id

  if not rpc.is_connected() then
    vim.notify("Warning: RPC not connected, game may be stuck", vim.log.levels.WARN)
  end

  rpc.notify("challenge_complete", result)
end

--- Clean up challenge state and return to game
function M._cleanup()
  local keyforge = require("keyforge")

  -- Close info window
  close_info_window()

  -- Stop keystroke tracking (if not already stopped)
  challenges.stop_tracking()

  -- Delete temp file
  if M._filepath then
    delete_temp_file(M._filepath)
  end

  -- Close challenge tab/buffer
  if M._challenge_buf and vim.api.nvim_buf_is_valid(M._challenge_buf) then
    -- Remove autocmds first to prevent recursive calls
    pcall(vim.api.nvim_del_augroup_by_name, "KeyforgeChallenge")
    vim.api.nvim_buf_delete(M._challenge_buf, { force = true })
  end

  -- Reset state
  M._current = nil
  M._request_id = nil
  M._filepath = nil
  M._challenge_buf = nil
  M._challenge_win = nil
  M._challenge_tab = nil
  M._initial_content = nil

  -- Update keyforge state
  keyforge._current_challenge_id = nil
  keyforge._game_state = "playing"

  -- Return to game tab
  keyforge.focus_game_tab()
end

--- Check if a challenge is currently active
---@return boolean
function M.is_active()
  return M._current ~= nil
end

--- Get current challenge info
---@return table|nil
function M.get_current()
  return M._current
end

return M
