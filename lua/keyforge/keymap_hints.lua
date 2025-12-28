--- Keymap discovery and hints for Keyforge
--- Discovers user's keybindings and provides context-aware hints for challenges
local M = {}

-- Cache structure for keymaps
M._cache = {
  normal = {},      -- mode mappings: lhs -> {rhs, desc}
  visual = {},
  insert = {},
  plugins = {},     -- detected plugins
  last_refresh = 0,
}

-- Cache timeout in seconds (5 minutes)
local CACHE_TIMEOUT = 300

-- Plugins to detect with their require paths or check functions
local plugins_to_detect = {
  { name = "telescope", require_path = "telescope" },
  { name = "nvim-tree", require_path = "nvim-tree" },
  { name = "neo-tree", require_path = "neo-tree" },
  { name = "oil", require_path = "oil" },
  { name = "fugitive", check = function() return vim.fn.exists(":Git") > 0 end },
  { name = "gitsigns", require_path = "gitsigns" },
  { name = "nvim-surround", require_path = "nvim-surround" },
  { name = "mini.surround", require_path = "mini.surround" },
  { name = "flash", require_path = "flash" },
  { name = "leap", require_path = "leap" },
  { name = "hop", require_path = "hop" },
  { name = "trouble", require_path = "trouble" },
  { name = "harpoon", require_path = "harpoon" },
}

-- Action mappings: maps common actions to patterns to search for in keymaps
-- Each action can have multiple patterns that might indicate that keybinding
local action_patterns = {
  -- Telescope / fuzzy finder
  find_files = {
    patterns = { "find_files", "files", "Telescope find" },
    category = "lsp-navigation",
    description = "Find files",
    default_binding = "<leader>ff",
  },
  live_grep = {
    patterns = { "live_grep", "grep", "Telescope grep", "search" },
    category = "search-replace",
    description = "Search in files",
    default_binding = "<leader>fg",
  },
  buffers = {
    patterns = { "buffers", "Telescope buffer" },
    category = "movement",
    description = "List buffers",
    default_binding = "<leader>fb",
  },
  help_tags = {
    patterns = { "help_tags", "Telescope help" },
    category = "lsp-navigation",
    description = "Search help",
    default_binding = "<leader>fh",
  },
  -- LSP
  goto_definition = {
    patterns = { "definition", "lsp.*definition", "vim.lsp.*definition" },
    category = "lsp-navigation",
    description = "Go to definition",
    default_binding = "gd",
  },
  goto_references = {
    patterns = { "references", "lsp.*references" },
    category = "lsp-navigation",
    description = "Find references",
    default_binding = "gr",
  },
  hover = {
    patterns = { "hover", "lsp.*hover" },
    category = "lsp-navigation",
    description = "Hover documentation",
    default_binding = "K",
  },
  rename = {
    patterns = { "rename", "lsp.*rename" },
    category = "lsp-navigation",
    description = "Rename symbol",
    default_binding = "<leader>rn",
  },
  code_action = {
    patterns = { "code_action", "lsp.*code_action" },
    category = "lsp-navigation",
    description = "Code actions",
    default_binding = "<leader>ca",
  },
  -- File explorer
  file_explorer = {
    patterns = { "NvimTree", "Neotree", "Oil", "Explore", "tree" },
    category = "movement",
    description = "File explorer",
    default_binding = "<leader>e",
  },
  -- Git
  git_status = {
    patterns = { "Git status", "Gitsigns", "git.*status" },
    category = "git-operations",
    description = "Git status",
    default_binding = "<leader>gs",
  },
  git_blame = {
    patterns = { "blame", "Gitsigns.*blame" },
    category = "git-operations",
    description = "Git blame",
    default_binding = "<leader>gb",
  },
  -- Surround
  surround_add = {
    patterns = { "surround", "ys" },
    category = "text-objects",
    description = "Add surround",
    default_binding = "ys",
  },
  surround_delete = {
    patterns = { "surround.*delete", "ds" },
    category = "text-objects",
    description = "Delete surround",
    default_binding = "ds",
  },
  surround_change = {
    patterns = { "surround.*change", "cs" },
    category = "text-objects",
    description = "Change surround",
    default_binding = "cs",
  },
  -- Window navigation
  window_left = {
    patterns = { "wincmd h", "<C-w>h" },
    category = "movement",
    description = "Window left",
    default_binding = "<C-w>h",
  },
  window_down = {
    patterns = { "wincmd j", "<C-w>j" },
    category = "movement",
    description = "Window down",
    default_binding = "<C-w>j",
  },
  window_up = {
    patterns = { "wincmd k", "<C-w>k" },
    category = "movement",
    description = "Window up",
    default_binding = "<C-w>k",
  },
  window_right = {
    patterns = { "wincmd l", "<C-w>l" },
    category = "movement",
    description = "Window right",
    default_binding = "<C-w>l",
  },
}

--- Check if a string matches any of the patterns
---@param str string
---@param patterns string[]
---@return boolean
local function matches_any_pattern(str, patterns)
  if not str then return false end
  local lower_str = str:lower()
  for _, pattern in ipairs(patterns) do
    if lower_str:find(pattern:lower(), 1, true) then
      return true
    end
  end
  return false
end

--- Detect installed plugins
---@return table<string, boolean>
function M.detect_plugins()
  local detected = {}

  for _, plugin in ipairs(plugins_to_detect) do
    local ok = false
    if plugin.check then
      ok = plugin.check()
    elseif plugin.require_path then
      ok = pcall(require, plugin.require_path)
    end
    detected[plugin.name] = ok
  end

  return detected
end

--- Discover keymaps for a given mode
---@param mode string Mode character ('n', 'v', 'i')
---@return table<string, table> Keymaps indexed by lhs
local function discover_mode_keymaps(mode)
  local keymaps = {}
  local mappings = vim.api.nvim_get_keymap(mode)

  for _, map in ipairs(mappings) do
    local lhs = map.lhs
    local rhs = map.rhs or ""
    local desc = map.desc or ""
    local callback = map.callback

    -- Store both rhs and description for pattern matching
    keymaps[lhs] = {
      rhs = rhs,
      desc = desc,
      has_callback = callback ~= nil,
    }
  end

  -- Also get buffer-local mappings for current buffer
  local buf_mappings = vim.api.nvim_buf_get_keymap(0, mode)
  for _, map in ipairs(buf_mappings) do
    local lhs = map.lhs
    keymaps[lhs] = {
      rhs = map.rhs or "",
      desc = map.desc or "",
      has_callback = map.callback ~= nil,
      buffer_local = true,
    }
  end

  return keymaps
end

--- Discover all keymaps and cache them
function M.discover_keymaps()
  M._cache.normal = discover_mode_keymaps("n")
  M._cache.visual = discover_mode_keymaps("v")
  M._cache.insert = discover_mode_keymaps("i")
  M._cache.plugins = M.detect_plugins()
  M._cache.last_refresh = os.time()
end

--- Refresh cache if stale
function M.refresh_cache()
  local now = os.time()
  if now - M._cache.last_refresh > CACHE_TIMEOUT then
    M.discover_keymaps()
  end
end

--- Find keybinding for an action by searching keymaps
---@param action string Action name from action_patterns
---@param mode? string Mode to search ('n', 'v', 'i'), defaults to 'n'
---@return string|nil lhs The keybinding if found
---@return string|nil source Description of where it was found
function M.find_binding_for_action(action, mode)
  mode = mode or "n"
  M.refresh_cache()

  local action_info = action_patterns[action]
  if not action_info then
    return nil, nil
  end

  local mode_cache = M._cache[({ n = "normal", v = "visual", i = "insert" })[mode]]
  if not mode_cache then
    return nil, nil
  end

  -- Search through keymaps for matching patterns
  for lhs, map_info in pairs(mode_cache) do
    -- Check rhs and desc against patterns
    local search_str = (map_info.rhs or "") .. " " .. (map_info.desc or "")
    if matches_any_pattern(search_str, action_info.patterns) then
      return lhs, map_info.desc or action_info.description
    end
  end

  return nil, nil
end

--- Get hint for an action, falling back to default if not found
---@param action string Action name
---@return table hint { binding = string, description = string, is_default = boolean }
function M.get_hint_for_action(action)
  local action_info = action_patterns[action]
  if not action_info then
    return {
      binding = "?",
      description = "Unknown action: " .. action,
      is_default = true,
    }
  end

  local binding, source = M.find_binding_for_action(action)

  if binding then
    return {
      binding = binding,
      description = source or action_info.description,
      is_default = false,
    }
  end

  -- Fall back to default
  return {
    binding = action_info.default_binding,
    description = action_info.description,
    is_default = true,
  }
end

--- Get hints for a category (returns all relevant actions)
---@param category string Challenge category
---@return table[] hints List of hint tables
function M.get_hints_for_category(category)
  M.refresh_cache()

  local hints = {}

  for action, info in pairs(action_patterns) do
    if info.category == category then
      local hint = M.get_hint_for_action(action)
      hint.action = action
      table.insert(hints, hint)
    end
  end

  return hints
end

--- Get all detected plugins
---@return table<string, boolean>
function M.get_detected_plugins()
  M.refresh_cache()
  return vim.deepcopy(M._cache.plugins)
end

--- Check if a specific plugin is detected
---@param plugin_name string
---@return boolean
function M.has_plugin(plugin_name)
  M.refresh_cache()
  return M._cache.plugins[plugin_name] or false
end

--- Force refresh the cache
function M.force_refresh()
  M._cache.last_refresh = 0
  M.discover_keymaps()
end

--- Get cache info (for debugging)
---@return table
function M.get_cache_info()
  return {
    normal_count = vim.tbl_count(M._cache.normal),
    visual_count = vim.tbl_count(M._cache.visual),
    insert_count = vim.tbl_count(M._cache.insert),
    plugins = M._cache.plugins,
    last_refresh = M._cache.last_refresh,
    age_seconds = os.time() - M._cache.last_refresh,
  }
end

--- Format hints for display in challenge buffer
---@param hints table[] List of hints
---@return string[] lines Formatted lines for display
function M.format_hints_for_display(hints)
  local lines = {}

  if #hints == 0 then
    return lines
  end

  table.insert(lines, "Your Keybindings:")

  for _, hint in ipairs(hints) do
    local marker = hint.is_default and "(default)" or ""
    local line = string.format("  %s: %s %s", hint.description, hint.binding, marker)
    table.insert(lines, line)
  end

  return lines
end

return M
