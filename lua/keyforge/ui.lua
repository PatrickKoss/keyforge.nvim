--- UI helpers for Keyforge challenge buffers
local M = {}

-- Active challenge state
M._challenge_buf = nil
M._challenge_win = nil
M._challenge_data = nil
M._initial_content = nil

--- Create a challenge buffer with the given content
---@param challenge table Challenge data from the game
---@return number bufnr The buffer number
function M.create_challenge_buffer(challenge)
  -- Close any existing challenge buffer
  M.close_challenge_buffer()

  -- Create new buffer
  local buf = vim.api.nvim_create_buf(false, true)

  -- Set buffer content
  local lines = vim.split(challenge.initial_buffer or "", "\n")
  vim.api.nvim_buf_set_lines(buf, 0, -1, false, lines)

  -- Store initial content for validation
  M._initial_content = vim.api.nvim_buf_get_lines(buf, 0, -1, false)

  -- Set buffer options
  vim.bo[buf].buftype = "nofile"
  vim.bo[buf].bufhidden = "wipe"
  vim.bo[buf].swapfile = false
  vim.bo[buf].modifiable = true

  -- Set filetype for syntax highlighting
  local filetype = challenge.filetype or "text"
  vim.bo[buf].filetype = filetype

  -- Set buffer name
  local name = string.format("keyforge://challenge/%s", challenge.id or "unnamed")
  vim.api.nvim_buf_set_name(buf, name)

  -- Store state
  M._challenge_buf = buf
  M._challenge_data = challenge

  return buf
end

--- Open the challenge buffer in a window
---@param buf? number Buffer number (uses current challenge buffer if not provided)
---@return number winid The window ID
function M.open_challenge_window(buf)
  buf = buf or M._challenge_buf
  if not buf or not vim.api.nvim_buf_is_valid(buf) then
    vim.notify("No challenge buffer to open", vim.log.levels.ERROR)
    return -1
  end

  -- Create a split below the game
  vim.cmd("botright split")
  vim.cmd("resize 12")

  local win = vim.api.nvim_get_current_win()
  vim.api.nvim_win_set_buf(win, buf)

  -- Set window options
  vim.wo[win].number = true
  vim.wo[win].relativenumber = true
  vim.wo[win].signcolumn = "yes"
  vim.wo[win].cursorline = true

  M._challenge_win = win

  -- Set up keymaps for the challenge buffer
  M._setup_challenge_keymaps(buf)

  -- Display challenge description
  M._show_challenge_info()

  return win
end

--- Set up keymaps for challenge buffer
---@param buf number
function M._setup_challenge_keymaps(buf)
  -- Complete challenge with <leader>kc
  vim.keymap.set("n", "<leader>kc", function()
    M.complete_challenge()
  end, { buffer = buf, desc = "Complete Keyforge challenge" })

  -- Skip challenge with <leader>ks
  vim.keymap.set("n", "<leader>ks", function()
    M.skip_challenge()
  end, { buffer = buf, desc = "Skip Keyforge challenge" })
end

--- Show challenge information in a floating window
function M._show_challenge_info()
  if not M._challenge_data then
    return
  end

  local challenge = M._challenge_data
  local lines = {
    "╭─── CHALLENGE ───╮",
    string.format("│ %s", challenge.name or "Unnamed"),
    "├──────────────────",
  }

  -- Add description lines
  local desc = challenge.description or "Complete the editing task"
  for _, line in ipairs(vim.split(desc, "\n")) do
    table.insert(lines, string.format("│ %s", line))
  end

  table.insert(lines, "├──────────────────")
  table.insert(lines, string.format("│ Category: %s", challenge.category or "general"))
  table.insert(lines, string.format("│ Difficulty: %d", challenge.difficulty or 1))

  if challenge.par_keystrokes then
    table.insert(lines, string.format("│ Par: %d keystrokes", challenge.par_keystrokes))
  end

  table.insert(lines, "├──────────────────")
  table.insert(lines, "│ <leader>kc - Complete")
  table.insert(lines, "│ <leader>ks - Skip")
  table.insert(lines, "╰──────────────────╯")

  -- Create floating window for info
  local info_buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_buf_set_lines(info_buf, 0, -1, false, lines)
  vim.bo[info_buf].modifiable = false

  local width = 24
  local height = #lines
  local win_config = {
    relative = "editor",
    width = width,
    height = height,
    row = 2,
    col = vim.o.columns - width - 2,
    style = "minimal",
    border = "none",
  }

  local info_win = vim.api.nvim_open_win(info_buf, false, win_config)
  vim.wo[info_win].winblend = 10

  -- Close info window when challenge buffer is closed
  vim.api.nvim_create_autocmd("BufWipeout", {
    buffer = M._challenge_buf,
    callback = function()
      if vim.api.nvim_win_is_valid(info_win) then
        vim.api.nvim_win_close(info_win, true)
      end
    end,
    once = true,
  })
end

--- Close the challenge buffer
function M.close_challenge_buffer()
  if M._challenge_win and vim.api.nvim_win_is_valid(M._challenge_win) then
    vim.api.nvim_win_close(M._challenge_win, true)
  end

  if M._challenge_buf and vim.api.nvim_buf_is_valid(M._challenge_buf) then
    vim.api.nvim_buf_delete(M._challenge_buf, { force = true })
  end

  M._challenge_buf = nil
  M._challenge_win = nil
  M._challenge_data = nil
  M._initial_content = nil
end

--- Get the current buffer content
---@return string[]
function M.get_buffer_content()
  if not M._challenge_buf or not vim.api.nvim_buf_is_valid(M._challenge_buf) then
    return {}
  end
  return vim.api.nvim_buf_get_lines(M._challenge_buf, 0, -1, false)
end

--- Complete the current challenge
---@return table result Completion result
function M.complete_challenge()
  if not M._challenge_buf then
    vim.notify("No active challenge", vim.log.levels.WARN)
    return { success = false, error = "No active challenge" }
  end

  local final_content = M.get_buffer_content()
  local challenges = require("keyforge.challenges")

  -- Validate the challenge
  local result = challenges.validate(M._challenge_data, M._initial_content, final_content)

  if result.success then
    vim.notify(string.format("Challenge complete! Score: %.0f%%", result.efficiency * 100), vim.log.levels.INFO)
  else
    vim.notify("Challenge failed: " .. (result.error or "Unknown error"), vim.log.levels.WARN)
  end

  -- Send result to game via RPC (if connected)
  local rpc = require("keyforge.rpc")
  if rpc.is_connected() then
    rpc.notify("challenge_complete", {
      request_id = M._challenge_data.request_id,
      success = result.success,
      keystroke_count = result.keystroke_count or 0,
      time_ms = result.time_ms or 0,
      efficiency = result.efficiency or 0,
    })
  end

  -- Close the challenge buffer
  M.close_challenge_buffer()

  return result
end

--- Skip the current challenge
function M.skip_challenge()
  if not M._challenge_buf then
    vim.notify("No active challenge", vim.log.levels.WARN)
    return
  end

  vim.notify("Challenge skipped", vim.log.levels.INFO)

  -- Notify game
  local rpc = require("keyforge.rpc")
  if rpc.is_connected() then
    rpc.notify("challenge_complete", {
      request_id = M._challenge_data.request_id,
      success = false,
      skipped = true,
    })
  end

  M.close_challenge_buffer()
end

--- Check if a challenge is active
---@return boolean
function M.is_challenge_active()
  return M._challenge_buf ~= nil and vim.api.nvim_buf_is_valid(M._challenge_buf)
end

return M
