---@class KeyforgeConfig
---@field keybind string Keybind to launch game (default: "<leader>kf")
---@field keybind_next_challenge string Keybind to start next challenge (default: "<leader>kn")
---@field keybind_complete string Keybind to complete challenge (default: "<leader>kc")
---@field keybind_skip string Keybind to skip challenge (default: "<leader>ks")
---@field difficulty string Difficulty level: "easy", "normal", "hard" (default: "normal")
---@field use_nerd_fonts boolean Use Nerd Font icons (default: true)
---@field starting_gold number Initial gold amount (default: 200)
---@field starting_health number Initial health (default: 100)
---@field auto_build boolean Auto-build binary on first run (default: true)

local M = {}

---@type KeyforgeConfig
local default_config = {
  keybind = "<leader>kf",
  keybind_next_challenge = "<leader>kn",
  keybind_complete = "<leader>kc",
  keybind_skip = "<leader>ks",
  difficulty = "normal",
  use_nerd_fonts = true,
  starting_gold = 200,
  starting_health = 100,
  auto_build = true,
}

---@type KeyforgeConfig
M.config = vim.deepcopy(default_config)

-- Module state
M._job_id = nil
M._term_buf = nil
M._term_win = nil
M._plugin_dir = nil

--- Get the plugin directory path
---@return string
local function get_plugin_dir()
  if M._plugin_dir then
    return M._plugin_dir
  end
  -- Find the plugin directory by looking for this file's location
  local info = debug.getinfo(1, "S")
  local script_path = info.source:sub(2) -- Remove the @ prefix
  -- Go up from lua/keyforge/init.lua to the plugin root
  M._plugin_dir = vim.fn.fnamemodify(script_path, ":h:h:h")
  return M._plugin_dir
end

--- Get the path to the game binary
---@return string
local function get_binary_path()
  return get_plugin_dir() .. "/game/bin/keyforge"
end

--- Check if the binary exists
---@return boolean
local function binary_exists()
  return vim.fn.filereadable(get_binary_path()) == 1
end

--- Build the game binary
---@param callback? function Callback on completion
function M.build(callback)
  local plugin_dir = get_plugin_dir()
  vim.notify("Building Keyforge...", vim.log.levels.INFO)

  vim.fn.jobstart({ "make", "build" }, {
    cwd = plugin_dir,
    on_exit = function(_, code)
      if code == 0 then
        vim.notify("Keyforge built successfully!", vim.log.levels.INFO)
        if callback then
          callback(true)
        end
      else
        vim.notify("Failed to build Keyforge. Run :KeyforgeBuild for details.", vim.log.levels.ERROR)
        if callback then
          callback(false)
        end
      end
    end,
    on_stderr = function(_, data)
      if data and #data > 0 and data[1] ~= "" then
        vim.notify("Build error: " .. table.concat(data, "\n"), vim.log.levels.ERROR)
      end
    end,
  })
end

--- Start the game
function M.start()
  -- Check if already running
  if M._job_id and vim.fn.jobwait({ M._job_id }, 0)[1] == -1 then
    vim.notify("Keyforge is already running!", vim.log.levels.WARN)
    -- Focus the existing window if possible
    if M._term_win and vim.api.nvim_win_is_valid(M._term_win) then
      vim.api.nvim_set_current_win(M._term_win)
    end
    return
  end

  -- Build if needed
  if not binary_exists() then
    if M.config.auto_build then
      M.build(function(success)
        if success then
          vim.schedule(function()
            M._launch_game()
          end)
        end
      end)
      return
    else
      vim.notify("Keyforge binary not found. Run :KeyforgeBuild first.", vim.log.levels.ERROR)
      return
    end
  end

  M._launch_game()
end

--- Launch the game in a fullscreen tab
function M._launch_game()
  local binary = get_binary_path()

  -- Create a new tab for fullscreen game
  vim.cmd("tabnew")

  -- Create terminal buffer
  M._term_buf = vim.api.nvim_get_current_buf()
  M._term_win = vim.api.nvim_get_current_win()

  -- Start the game process
  M._job_id = vim.fn.termopen(binary, {
    on_exit = function(_, code)
      M._job_id = nil
      if code ~= 0 then
        vim.notify("Keyforge exited with code " .. code, vim.log.levels.WARN)
      end
    end,
  })

  -- Set buffer options
  vim.api.nvim_buf_set_name(M._term_buf, "keyforge://game")
  vim.bo[M._term_buf].bufhidden = "wipe"

  -- Enter terminal mode
  vim.cmd("startinsert")

  -- Set up autocmd to clean up when buffer is closed
  vim.api.nvim_create_autocmd("BufWipeout", {
    buffer = M._term_buf,
    callback = function()
      M.stop()
    end,
    once = true,
  })
end

--- Stop the game
function M.stop()
  if M._job_id then
    vim.fn.jobstop(M._job_id)
    M._job_id = nil
  end

  if M._term_buf and vim.api.nvim_buf_is_valid(M._term_buf) then
    vim.api.nvim_buf_delete(M._term_buf, { force = true })
  end

  M._term_buf = nil
  M._term_win = nil
end

--- Start the next challenge (user-triggered)
function M.next_challenge()
  local challenge_queue = require("keyforge.challenge_queue")
  challenge_queue.request_next()
end

--- Complete the current challenge
function M.complete_challenge()
  local challenge_queue = require("keyforge.challenge_queue")
  challenge_queue.complete_current()
end

--- Skip the current challenge
function M.skip_challenge()
  local challenge_queue = require("keyforge.challenge_queue")
  challenge_queue.skip_current()
end

--- Get challenge statistics
---@return table stats
function M.get_challenge_stats()
  local challenge_queue = require("keyforge.challenge_queue")
  return challenge_queue.get_stats()
end

--- Setup the plugin with user configuration
---@param opts? KeyforgeConfig
function M.setup(opts)
  M.config = vim.tbl_deep_extend("force", default_config, opts or {})

  -- Register keybind for starting the game
  if M.config.keybind and M.config.keybind ~= "" then
    vim.keymap.set("n", M.config.keybind, function()
      M.start()
    end, { desc = "Start Keyforge" })
  end

  -- Register keybind for next challenge
  if M.config.keybind_next_challenge and M.config.keybind_next_challenge ~= "" then
    vim.keymap.set("n", M.config.keybind_next_challenge, function()
      M.next_challenge()
    end, { desc = "Keyforge: Next challenge" })
  end

  -- Register keybind for completing challenge
  if M.config.keybind_complete and M.config.keybind_complete ~= "" then
    vim.keymap.set("n", M.config.keybind_complete, function()
      M.complete_challenge()
    end, { desc = "Keyforge: Complete challenge" })
  end

  -- Register keybind for skipping challenge
  if M.config.keybind_skip and M.config.keybind_skip ~= "" then
    vim.keymap.set("n", M.config.keybind_skip, function()
      M.skip_challenge()
    end, { desc = "Keyforge: Skip challenge" })
  end

  -- Initialize challenge queue
  local challenge_queue = require("keyforge.challenge_queue")
  challenge_queue.init()
end

return M
