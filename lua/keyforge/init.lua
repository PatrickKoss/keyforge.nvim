---@class KeyforgeConfig
---@field keybind string Keybind to launch game (default: "<leader>kf")
---@field keybind_next_challenge string Keybind to start next challenge (default: "<leader>kn")
---@field keybind_complete string Keybind to complete challenge (default: "<leader>kc")
---@field keybind_skip string Keybind to skip challenge (default: "<leader>ks")
---@field keybind_submit string Keybind to submit challenge in buffer (default: "<CR>")
---@field keybind_cancel string Keybind to cancel challenge in buffer (default: "<Esc>")
---@field difficulty string Difficulty level: "easy", "normal", "hard" (default: "normal")
---@field game_speed number Game speed multiplier: 0.5, 1.0, 1.5, 2.0 (default: 1.0)
---@field use_nerd_fonts boolean Use Nerd Font icons (default: true)
---@field starting_gold number Initial gold amount (default: 200)
---@field starting_health number Initial health (default: 100)
---@field auto_build boolean Auto-build binary on first run (default: true)
---@field challenge_timeout number Challenge timeout in seconds (default: 300)

local M = {}

---@type KeyforgeConfig
local default_config = {
  keybind = "<leader>kf",
  keybind_next_challenge = "<leader>kn",
  keybind_complete = "<leader>kc",
  keybind_skip = "<leader>ks",
  keybind_submit = "<CR>",
  keybind_cancel = "<Esc>",
  difficulty = "normal",
  game_speed = 1.0,
  use_nerd_fonts = true,
  starting_gold = 200,
  starting_health = 100,
  auto_build = true,
  challenge_timeout = 300,
}

---@type KeyforgeConfig
M.config = vim.deepcopy(default_config)

-- Module state
M._job_id = nil
M._term_buf = nil
M._term_win = nil
M._term_tab = nil
M._plugin_dir = nil
M._current_challenge_id = nil
M._game_state = "idle" -- idle, playing, paused, challenge_waiting, game_over, victory
M._socket_path = nil   -- Unix socket path for RPC

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

--- Generate a unique socket path for RPC
---@return string
local function generate_socket_path()
  local pid = vim.fn.getpid()
  local timestamp = vim.fn.localtime()
  return string.format("/tmp/keyforge-%d-%d.sock", pid, timestamp)
end

--- Clean up stale socket files
local function cleanup_stale_sockets()
  local pattern = "/tmp/keyforge-*.sock"
  local files = vim.fn.glob(pattern, false, true)
  for _, file in ipairs(files) do
    -- Try to remove stale sockets (ignore errors)
    pcall(vim.fn.delete, file)
  end
end

--- Connect to the game's RPC socket with retry
---@param socket_path string
---@param max_attempts? number
local function connect_rpc_with_retry(socket_path, max_attempts)
  max_attempts = max_attempts or 20
  local attempt = 0
  local delay = 100 -- Start with 100ms

  local rpc = require("keyforge.rpc")
  rpc.register_handlers()

  local function try_connect()
    attempt = attempt + 1
    rpc.connect(socket_path, function()
      -- Success
      vim.notify("Keyforge RPC connected!", vim.log.levels.INFO)
    end, function(err)
      -- Error
      if attempt < max_attempts then
        -- Retry with exponential backoff
        delay = math.min(delay * 1.5, 1000)
        vim.defer_fn(try_connect, delay)
      else
        vim.notify("Failed to connect to Keyforge RPC: " .. tostring(err), vim.log.levels.WARN)
      end
    end)
  end

  -- Initial delay to give the game time to start its socket server
  vim.defer_fn(try_connect, 200)
end

--- Launch the game in a fullscreen tab
function M._launch_game()
  local binary = get_binary_path()

  -- Clean up any stale socket files
  cleanup_stale_sockets()

  -- Generate unique socket path
  M._socket_path = generate_socket_path()

  -- Create a new tab for fullscreen game
  vim.cmd("tabnew")
  M._term_tab = vim.api.nvim_get_current_tabpage()

  -- Create terminal buffer
  M._term_buf = vim.api.nvim_get_current_buf()
  M._term_win = vim.api.nvim_get_current_win()

  -- Start the game process in nvim mode with socket RPC
  -- Pass config values as command-line flags
  local cmd = string.format(
    "%s --nvim-mode --rpc-socket %s --difficulty %s --game-speed %.1f --starting-gold %d --starting-health %d",
    binary,
    M._socket_path,
    M.config.difficulty,
    M.config.game_speed,
    M.config.starting_gold,
    M.config.starting_health
  )
  M._job_id = vim.fn.termopen(cmd, {
    on_exit = function(_, code)
      vim.schedule(function()
        M._job_id = nil
        M._game_state = "idle"

        -- Disconnect RPC and clean up socket
        local rpc = require("keyforge.rpc")
        rpc.disconnect()

        if M._socket_path then
          pcall(vim.fn.delete, M._socket_path)
          M._socket_path = nil
        end

        if code ~= 0 then
          vim.notify("Keyforge exited with code " .. code, vim.log.levels.WARN)
        end
      end)
    end,
  })

  -- Set buffer options
  vim.api.nvim_buf_set_name(M._term_buf, "keyforge://game")
  vim.bo[M._term_buf].bufhidden = "wipe"

  -- Enter terminal mode
  vim.cmd("startinsert")

  -- Update game state
  M._game_state = "playing"

  -- Connect to RPC socket with retry
  connect_rpc_with_retry(M._socket_path)

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
  -- Disconnect RPC first
  local rpc = require("keyforge.rpc")
  rpc.disconnect()

  if M._job_id then
    vim.fn.jobstop(M._job_id)
    M._job_id = nil
  end

  if M._term_buf and vim.api.nvim_buf_is_valid(M._term_buf) then
    vim.api.nvim_buf_delete(M._term_buf, { force = true })
  end

  -- Clean up socket file
  if M._socket_path then
    pcall(vim.fn.delete, M._socket_path)
    M._socket_path = nil
  end

  M._term_buf = nil
  M._term_win = nil
  M._term_tab = nil
  M._game_state = "idle"
  M._current_challenge_id = nil
end

--- Get the game tab page
---@return number|nil
function M.get_game_tab()
  if M._term_tab and vim.api.nvim_tabpage_is_valid(M._term_tab) then
    return M._term_tab
  end
  return nil
end

--- Focus the game tab
function M.focus_game_tab()
  local tab = M.get_game_tab()
  if tab then
    vim.api.nvim_set_current_tabpage(tab)
    if M._term_win and vim.api.nvim_win_is_valid(M._term_win) then
      vim.api.nvim_set_current_win(M._term_win)
      vim.cmd("startinsert")
    end
  end
end

--- Send a message to the game via RPC
---@param method string
---@param params? table
function M.send_to_game(method, params)
  local rpc = require("keyforge.rpc")
  rpc.notify(method, params)
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
