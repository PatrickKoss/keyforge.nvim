--- Game over/victory display for Keyforge
local M = {}

M._win = nil
M._buf = nil

--- Show the game over or victory screen
---@param params table Game state params (state, wave, gold, health, towers)
function M.show(params)
  local state = params.state or "game_over"

  -- Close any existing window
  M.close()

  -- Build content based on state
  local lines
  local title

  if state == "victory" then
    title = " VICTORY! "
    lines = {
      "",
      "  ██╗   ██╗██╗ ██████╗████████╗ ██████╗ ██████╗ ██╗   ██╗██╗",
      "  ██║   ██║██║██╔════╝╚══██╔══╝██╔═══██╗██╔══██╗╚██╗ ██╔╝██║",
      "  ██║   ██║██║██║        ██║   ██║   ██║██████╔╝ ╚████╔╝ ██║",
      "  ╚██╗ ██╔╝██║██║        ██║   ██║   ██║██╔══██╗  ╚██╔╝  ╚═╝",
      "   ╚████╔╝ ██║╚██████╗   ██║   ╚██████╔╝██║  ██║   ██║   ██╗",
      "    ╚═══╝  ╚═╝ ╚═════╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚═╝",
      "",
      "  Congratulations! You defended against all waves!",
      "",
      string.format("  Final Health: %d/%d", params.health or 0, 100),
      string.format("  Final Gold:   %d", params.gold or 0),
      string.format("  Towers Built: %d", params.towers or 0),
      "",
      "  Press [r] to play again, [l] for level select, or [q] to quit",
      "",
    }
  else
    title = " GAME OVER "
    lines = {
      "",
      "   ██████╗  █████╗ ███╗   ███╗███████╗     ██████╗ ██╗   ██╗███████╗██████╗ ",
      "  ██╔════╝ ██╔══██╗████╗ ████║██╔════╝    ██╔═══██╗██║   ██║██╔════╝██╔══██╗",
      "  ██║  ███╗███████║██╔████╔██║█████╗      ██║   ██║██║   ██║█████╗  ██████╔╝",
      "  ██║   ██║██╔══██║██║╚██╔╝██║██╔══╝      ██║   ██║╚██╗ ██╔╝██╔══╝  ██╔══██╗",
      "  ╚██████╔╝██║  ██║██║ ╚═╝ ██║███████╗    ╚██████╔╝ ╚████╔╝ ███████╗██║  ██║",
      "   ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝     ╚═════╝   ╚═══╝  ╚══════╝╚═╝  ╚═╝",
      "",
      string.format("  Wave Reached: %d/10", params.wave or 1),
      string.format("  Final Gold:   %d", params.gold or 0),
      string.format("  Towers Built: %d", params.towers or 0),
      "",
      "  Press [r] to restart, [l] for level select, or [q] to quit",
      "",
    }
  end

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
    title = title,
    title_pos = "center",
  })

  -- Set highlight based on state
  if state == "victory" then
    vim.api.nvim_win_set_option(M._win, "winhl", "Normal:DiffAdd,FloatBorder:DiffAdd")
  else
    vim.api.nvim_win_set_option(M._win, "winhl", "Normal:DiffDelete,FloatBorder:DiffDelete")
  end

  -- Set up keymaps
  vim.keymap.set("n", "r", function()
    M.restart()
  end, { buffer = M._buf, desc = "Restart game" })

  vim.keymap.set("n", "q", function()
    M.quit()
  end, { buffer = M._buf, desc = "Quit game" })

  vim.keymap.set("n", "l", function()
    M.level_select()
  end, { buffer = M._buf, desc = "Go to level select" })

  -- Close on any other key after a delay
  vim.keymap.set("n", "<Esc>", function()
    M.close()
  end, { buffer = M._buf, desc = "Close" })
end

--- Close the game over window
function M.close()
  if M._win and vim.api.nvim_win_is_valid(M._win) then
    vim.api.nvim_win_close(M._win, true)
  end
  M._win = nil

  if M._buf and vim.api.nvim_buf_is_valid(M._buf) then
    vim.api.nvim_buf_delete(M._buf, { force = true })
  end
  M._buf = nil
end

--- Restart the game
function M.restart()
  local keyforge = require("keyforge")
  local rpc = require("keyforge.rpc")

  -- Close popup first
  M.close()

  -- Send restart command to game
  rpc.notify("restart_game", {})

  -- Schedule focus to run after RPC is sent and popup is fully closed
  vim.schedule(function()
    keyforge.focus_game_tab()
  end)
end

--- Quit the game
function M.quit()
  local keyforge = require("keyforge")

  -- Close popup first
  M.close()

  -- Schedule stop to run after current event loop completes
  -- This ensures the popup window is fully closed before buffer cleanup
  vim.schedule(function()
    keyforge.stop()
  end)
end

--- Go to level selection
function M.level_select()
  local keyforge = require("keyforge")
  local rpc = require("keyforge.rpc")

  -- Close popup first
  M.close()

  -- Send level select command to game
  rpc.notify("go_to_level_select", {})

  -- Schedule focus to run after RPC is sent and popup is fully closed
  vim.schedule(function()
    keyforge.focus_game_tab()
  end)
end

return M
