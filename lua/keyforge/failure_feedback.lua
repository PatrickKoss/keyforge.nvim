--- Failure feedback floating window for Keyforge
--- Shows detailed feedback when a challenge fails with diff, hints, and retry/skip controls
local M = {}

M._win = nil
M._buf = nil
M._callbacks = nil
M._challenge = nil

--- Add diff lines to content (for exact_match failures)
---@param lines string[] Lines to append to
---@param failure_details table Failure details
local function add_diff_section(lines, failure_details)
  if not failure_details or not failure_details.diff_lines then
    return
  end
  if #failure_details.diff_lines == 0 then
    return
  end

  table.insert(lines, "  Expected vs Actual:")
  local diff_count = 0
  for _, diff_line in ipairs(failure_details.diff_lines) do
    if diff_count >= 8 then
      table.insert(lines, "    ... (truncated)")
      break
    end
    table.insert(lines, "    " .. diff_line)
    diff_count = diff_count + 1
  end
  table.insert(lines, "")
end

--- Add cursor position comparison (for cursor_position failures)
---@param lines string[] Lines to append to
---@param failure_details table Failure details
local function add_cursor_position_section(lines, failure_details)
  if not failure_details or failure_details.validation_type ~= "cursor_position" then
    return
  end

  local exp = failure_details.expected
  local act = failure_details.actual
  if exp and act then
    table.insert(lines, string.format("  Expected: row %d, column %d", exp.row, exp.col))
    table.insert(lines, string.format("  Actual:   row %d, column %d", act.row, act.col))
    table.insert(lines, "")
  end
end

--- Add character comparison (for cursor_on_char failures)
---@param lines string[] Lines to append to
---@param failure_details table Failure details
local function add_cursor_char_section(lines, failure_details)
  if not failure_details or failure_details.validation_type ~= "cursor_on_char" then
    return
  end

  table.insert(lines, string.format("  Expected char: '%s'", failure_details.expected or "?"))
  table.insert(lines, string.format("  Cursor on:     '%s'", failure_details.actual or "?"))
  table.insert(lines, "")
end

--- Add hint section
---@param lines string[] Lines to append to
---@param challenge table Challenge data
local function add_hint_section(lines, challenge)
  local challenges = require("keyforge.challenges")

  table.insert(lines, "  " .. string.rep("-", 40))

  local hint = challenges.get_challenge_hint(challenge)
  if hint and hint ~= "" then
    if #hint > 45 then
      table.insert(lines, "  Hint: " .. hint:sub(1, 45))
      table.insert(lines, "        " .. hint:sub(46))
    else
      table.insert(lines, "  Hint: " .. hint)
    end
  else
    table.insert(lines, "  Hint: " .. (challenge.description or "Complete the challenge"))
  end
  table.insert(lines, "")
end

--- Build the content lines for the failure feedback window
---@param challenge table Challenge data
---@param failure_details table Failure details from validation
---@return string[] lines Content lines for the window
local function build_content(challenge, failure_details)
  local lines = {
    "",
    "  Challenge: " .. (challenge.name or "Unknown"),
    "",
  }

  -- Validation-specific failure message
  local message = failure_details and failure_details.message or "Challenge requirements not met"
  table.insert(lines, "  " .. message)
  table.insert(lines, "")

  -- Add validation-specific sections
  add_diff_section(lines, failure_details)
  add_cursor_position_section(lines, failure_details)
  add_cursor_char_section(lines, failure_details)

  -- Add hint and controls
  add_hint_section(lines, challenge)
  table.insert(lines, "  [r] Retry   [s] Skip   [Esc] Close")
  table.insert(lines, "")

  return lines
end

--- Show the failure feedback window
---@param challenge table Challenge data
---@param failure_details table|nil Failure details from validation
---@param callbacks table { on_retry: function, on_skip: function }
function M.show(challenge, failure_details, callbacks)
  -- Close existing window if any
  M.close()

  M._callbacks = callbacks
  M._challenge = challenge

  local lines = build_content(challenge, failure_details)

  -- Calculate window size
  local width = 0
  for _, line in ipairs(lines) do
    width = math.max(width, vim.fn.strdisplaywidth(line))
  end
  width = math.min(width + 4, vim.o.columns - 4)
  local height = #lines

  -- Center position
  local row = math.floor((vim.o.lines - height) / 2)
  local col = math.floor((vim.o.columns - width) / 2)

  -- Create buffer
  M._buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_buf_set_lines(M._buf, 0, -1, false, lines)
  vim.api.nvim_buf_set_option(M._buf, "modifiable", false)
  vim.api.nvim_buf_set_option(M._buf, "buftype", "nofile")
  vim.api.nvim_buf_set_option(M._buf, "bufhidden", "wipe")

  -- Create floating window
  M._win = vim.api.nvim_open_win(M._buf, true, {
    relative = "editor",
    width = width,
    height = height,
    row = row,
    col = col,
    style = "minimal",
    border = "rounded",
    title = " Challenge Failed ",
    title_pos = "center",
  })

  -- Red-ish highlight for failure
  vim.api.nvim_win_set_option(M._win, "winhl", "Normal:DiffDelete,FloatBorder:DiffDelete")

  -- Set up keymaps
  M._setup_keymaps()
end

--- Set up keymaps for the failure window
function M._setup_keymaps()
  local buf = M._buf

  vim.keymap.set("n", "r", function()
    local cb = M._callbacks
    M.close()
    if cb and cb.on_retry then
      vim.schedule(function()
        cb.on_retry()
      end)
    end
  end, { buffer = buf, desc = "Retry challenge" })

  vim.keymap.set("n", "s", function()
    local cb = M._callbacks
    M.close()
    if cb and cb.on_skip then
      vim.schedule(function()
        cb.on_skip()
      end)
    end
  end, { buffer = buf, desc = "Skip challenge" })

  vim.keymap.set("n", "<Esc>", function()
    M.close()
  end, { buffer = buf, desc = "Close feedback" })

  vim.keymap.set("n", "q", function()
    M.close()
  end, { buffer = buf, desc = "Close feedback" })
end

--- Close the failure feedback window
function M.close()
  if M._win and vim.api.nvim_win_is_valid(M._win) then
    vim.api.nvim_win_close(M._win, true)
  end
  M._win = nil

  if M._buf and vim.api.nvim_buf_is_valid(M._buf) then
    vim.api.nvim_buf_delete(M._buf, { force = true })
  end
  M._buf = nil
  M._callbacks = nil
  M._challenge = nil
end

--- Check if feedback window is currently shown
---@return boolean
function M.is_showing()
  return M._win ~= nil and vim.api.nvim_win_is_valid(M._win)
end

return M
