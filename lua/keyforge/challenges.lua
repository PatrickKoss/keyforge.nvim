--- Challenge validation and scoring for Keyforge
local M = {}

-- Plugin detection cache (persists for session)
M._plugin_cache = {}

-- Action patterns for runtime keymap resolution
-- Maps hint_action values to search patterns for rhs and desc
M.action_patterns = {
  -- Telescope actions
  find_files = {
    rhs = { "telescope%.builtin%.find_files", "Telescope find_files", "fzf#" },
    desc = { "find file", "find files", "file finder" },
  },
  live_grep = {
    rhs = { "telescope%.builtin%.live_grep", "Telescope live_grep", "live grep" },
    desc = { "live grep", "search in files", "grep" },
  },
  buffers = {
    rhs = { "telescope%.builtin%.buffers", "Telescope buffers" },
    desc = { "buffers", "find buffer" },
  },
  oldfiles = {
    rhs = { "telescope%.builtin%.oldfiles", "Telescope oldfiles", "recent" },
    desc = { "recent", "oldfiles", "old files" },
  },
  help_tags = {
    rhs = { "telescope%.builtin%.help_tags", "Telescope help_tags" },
    desc = { "help", "help tags" },
  },
  keymaps = {
    rhs = { "telescope%.builtin%.keymaps", "Telescope keymaps" },
    desc = { "keymap", "keymaps" },
  },
  current_buffer_fuzzy = {
    rhs = { "telescope%.builtin%.current_buffer_fuzzy_find", "current_buffer_fuzzy" },
    desc = { "fuzzy find", "current buffer" },
  },
  resume = {
    rhs = { "telescope%.builtin%.resume", "Telescope resume" },
    desc = { "resume" },
  },

  -- LSP actions
  goto_definition = {
    rhs = { "vim%.lsp%.buf%.definition", "lsp.*definition" },
    desc = { "definition", "go to def" },
  },
  goto_declaration = {
    rhs = { "vim%.lsp%.buf%.declaration", "lsp.*declaration" },
    desc = { "declaration" },
  },
  goto_implementation = {
    rhs = { "vim%.lsp%.buf%.implementation", "lsp.*implementation" },
    desc = { "implementation" },
  },
  find_references = {
    rhs = { "vim%.lsp%.buf%.references", "lsp.*references" },
    desc = { "references", "find ref" },
  },
  hover = {
    rhs = { "vim%.lsp%.buf%.hover" },
    desc = { "hover" },
  },
  signature_help = {
    rhs = { "vim%.lsp%.buf%.signature_help" },
    desc = { "signature" },
  },
  code_action = {
    rhs = { "vim%.lsp%.buf%.code_action", "lsp.*code_action" },
    desc = { "code action" },
  },
  rename = {
    rhs = { "vim%.lsp%.buf%.rename", "lsp.*rename" },
    desc = { "rename" },
  },
  format = {
    rhs = { "vim%.lsp%.buf%.format", "conform%.format", "format" },
    desc = { "format" },
  },
  type_definition = {
    rhs = { "vim%.lsp%.buf%.type_definition" },
    desc = { "type def" },
  },
  document_symbols = {
    rhs = { "document_symbols", "DocumentSymbol" },
    desc = { "document symbol" },
  },
  workspace_symbols = {
    rhs = { "workspace_symbols", "WorkspaceSymbol" },
    desc = { "workspace symbol" },
  },

  -- Diagnostics
  diagnostic_float = {
    rhs = { "vim%.diagnostic%.open_float", "diagnostic.*float" },
    desc = { "diagnostic", "show error" },
  },
  diagnostic_next = {
    rhs = { "vim%.diagnostic%.goto_next", "diagnostic.*next" },
    desc = { "next diagnostic" },
  },
  diagnostic_prev = {
    rhs = { "vim%.diagnostic%.goto_prev", "diagnostic.*prev" },
    desc = { "prev diagnostic" },
  },

  -- Git/gitsigns actions
  git_next_hunk = {
    rhs = { "gitsigns%.next_hunk", "next_hunk" },
    desc = { "next hunk" },
  },
  git_prev_hunk = {
    rhs = { "gitsigns%.prev_hunk", "prev_hunk" },
    desc = { "prev hunk" },
  },
  git_stage_hunk = {
    rhs = { "gitsigns%.stage_hunk", "stage_hunk" },
    desc = { "stage hunk" },
  },
  git_reset_hunk = {
    rhs = { "gitsigns%.reset_hunk", "reset_hunk" },
    desc = { "reset hunk" },
  },
  git_preview_hunk = {
    rhs = { "gitsigns%.preview_hunk", "preview_hunk" },
    desc = { "preview hunk" },
  },
  git_blame_line = {
    rhs = { "gitsigns%.blame_line", "blame_line" },
    desc = { "blame", "git blame" },
  },
  git_toggle_blame = {
    rhs = { "gitsigns%.toggle_current_line_blame", "toggle.*blame" },
    desc = { "toggle blame" },
  },
  git_diff_this = {
    rhs = { "gitsigns%.diffthis", "diff_this" },
    desc = { "diff" },
  },
  git_stage_buffer = {
    rhs = { "gitsigns%.stage_buffer", "stage_buffer" },
    desc = { "stage buffer" },
  },

  -- Harpoon actions
  harpoon_add = {
    rhs = { "harpoon.*add", "mark%.add_file" },
    desc = { "add.*harpoon", "harpoon add" },
  },
  harpoon_menu = {
    rhs = { "harpoon.*menu", "harpoon.*toggle" },
    desc = { "harpoon menu", "harpoon toggle" },
  },
  harpoon_nav_1 = {
    rhs = { "harpoon.*1", "nav_file%(1%)" },
    desc = { "harpoon 1" },
  },
  harpoon_nav_2 = {
    rhs = { "harpoon.*2", "nav_file%(2%)" },
    desc = { "harpoon 2" },
  },

  -- Surround actions (mini.surround or nvim-surround)
  surround_add = {
    rhs = { "surround%.add", "MiniSurround%.add" },
    desc = { "surround add", "add surround" },
  },
  surround_delete = {
    rhs = { "surround%.delete", "MiniSurround%.delete" },
    desc = { "surround delete", "delete surround" },
  },
  surround_replace = {
    rhs = { "surround%.replace", "MiniSurround%.replace" },
    desc = { "surround replace", "change surround" },
  },

  -- Window management
  window_split_horizontal = {
    rhs = { "split" },
    desc = { "split horizontal", "horizontal split" },
  },
  window_split_vertical = {
    rhs = { "vsplit" },
    desc = { "split vertical", "vertical split" },
  },

  -- Buffer management
  buffer_close = {
    rhs = { "bdelete", "bd", "Bdelete" },
    desc = { "close buffer", "delete buffer" },
  },
  buffer_save = {
    rhs = { ":w<", "write" },
    desc = { "save", "write" },
  },

  -- Folding (nvim-ufo)
  fold_peek = {
    rhs = { "ufo%.peekFoldedLines", "preview.*fold" },
    desc = { "peek fold", "preview fold" },
  },
}

-- Plugin name aliases for detection
M.plugin_aliases = {
  ["telescope"] = { "telescope", "telescope.nvim", "telescope.builtin" },
  ["nvim-surround"] = { "nvim-surround", "mini.surround", "surround" },
  ["mini.surround"] = { "mini.surround", "nvim-surround", "surround" },
  ["harpoon"] = { "harpoon" },
  ["gitsigns"] = { "gitsigns", "gitsigns.nvim" },
  ["nvim-tree"] = { "nvim-tree", "nvim-tree.lua", "neo-tree" },
  ["nvim-ufo"] = { "ufo", "nvim-ufo" },
  ["conform"] = { "conform", "conform.nvim" },
  ["diffview"] = { "diffview", "diffview.nvim" },
}

--- Check if a plugin is available
---@param name string Plugin name to check
---@return boolean
function M.plugin_available(name)
  -- Check cache first
  if M._plugin_cache[name] ~= nil then
    return M._plugin_cache[name]
  end

  -- Get aliases to try
  local aliases = M.plugin_aliases[name] or { name, name:gsub("-", "_"), name:gsub("%.nvim$", "") }

  for _, alias in ipairs(aliases) do
    -- Try requiring the module
    local ok = pcall(require, alias)
    if ok then
      M._plugin_cache[name] = true
      return true
    end
  end

  -- Try checking if command exists (for vimscript plugins)
  local cmd_name = name:gsub("%-", ""):gsub("^%l", string.upper)
  if vim.fn.exists(":" .. cmd_name) == 2 then
    M._plugin_cache[name] = true
    return true
  end

  M._plugin_cache[name] = false
  return false
end

--- Resolve a hint_action to the user's actual keymap
---@param action string The action identifier (e.g., "find_files")
---@return string|nil lhs The keymap (e.g., "<leader>ff") or nil if not found
---@return string|nil desc The keymap description or nil
function M.resolve_keymap(action)
  local patterns = M.action_patterns[action]
  if not patterns then
    return nil, nil
  end

  -- Search all normal mode keymaps
  for _, keymap in ipairs(vim.api.nvim_get_keymap("n")) do
    -- Check rhs patterns
    if keymap.rhs then
      for _, pattern in ipairs(patterns.rhs or {}) do
        if keymap.rhs:lower():find(pattern:lower()) then
          return keymap.lhs, keymap.desc
        end
      end
    end
    -- Check callback (for lua functions, check desc)
    if keymap.desc then
      for _, pattern in ipairs(patterns.desc or {}) do
        if keymap.desc:lower():find(pattern:lower()) then
          return keymap.lhs, keymap.desc
        end
      end
    end
  end

  return nil, nil -- Fallback to hint_fallback
end

--- Format a keymap for display (convert <leader> to actual key, etc.)
---@param lhs string The raw keymap lhs (e.g., "<leader>ff")
---@return string formatted The formatted keymap (e.g., "Space f f")
function M.format_keymap_display(lhs)
  if not lhs then
    return ""
  end

  -- Get the mapleader value
  local leader = vim.g.mapleader or "\\"
  local leader_display = "Space"
  if leader == "\\" then
    leader_display = "\\"
  elseif leader == "," then
    leader_display = ","
  end

  local formatted = lhs
  -- Replace <leader> with actual leader key display
  formatted = formatted:gsub("<[Ll]eader>", leader_display .. " ")
  -- Replace common special keys
  formatted = formatted:gsub("<[Cc]%-(%w)>", "Ctrl+%1 ")
  formatted = formatted:gsub("<[Ss]%-(%w)>", "Shift+%1 ")
  formatted = formatted:gsub("<[Aa]%-(%w)>", "Alt+%1 ")
  formatted = formatted:gsub("<[Cc][Rr]>", "Enter")
  formatted = formatted:gsub("<[Ee][Ss][Cc]>", "Esc")
  formatted = formatted:gsub("<[Tt][Aa][Bb]>", "Tab")

  -- Add spaces between characters for readability (but not for special keys)
  -- Only if it's a simple sequence like "ff" -> "f f"
  if not formatted:find("[%+%s]") and #formatted > 1 then
    formatted = formatted:gsub("(.)", "%1 "):gsub("%s+$", "")
  end

  return formatted
end

--- Get the hint for a challenge (resolved keymap or fallback)
---@param challenge table Challenge data
---@return string hint The hint to display
function M.get_challenge_hint(challenge)
  -- If no hint_action, return the description as-is
  if not challenge.hint_action then
    return challenge.description or ""
  end

  -- Try to resolve the keymap
  local lhs, desc = M.resolve_keymap(challenge.hint_action)
  if lhs then
    local formatted = M.format_keymap_display(lhs)
    local hint = string.format("Use %s", formatted)
    if desc then
      hint = hint .. " (" .. desc .. ")"
    end
    return hint
  end

  -- Fallback to hint_fallback or description
  return challenge.hint_fallback or challenge.description or ""
end

--- Filter challenges by available plugins
---@param challenges table[] List of challenges
---@return table[] filtered Challenges the user can complete
function M.filter_challenges_by_plugins(challenges)
  local available = {}
  for _, challenge in ipairs(challenges) do
    if not challenge.required_plugin or M.plugin_available(challenge.required_plugin) then
      table.insert(available, challenge)
    end
  end
  return available
end

--- Clear the plugin cache (useful for testing or when plugins change)
function M.clear_plugin_cache()
  M._plugin_cache = {}
end

-- Keystroke tracking state
M._tracking = false
M._keystroke_count = 0
M._start_time = nil
M._on_key_ns = nil

--- Start tracking keystrokes
function M.start_tracking()
  M._keystroke_count = 0
  M._start_time = vim.loop.hrtime()
  M._tracking = true

  -- Set up keystroke tracking via vim.on_key
  M._on_key_ns = vim.on_key(function(key)
    if M._tracking and key ~= "" then
      M._keystroke_count = M._keystroke_count + 1
    end
  end)
end

--- Stop tracking keystrokes
---@return number keystrokes Total keystroke count
---@return number time_ms Time elapsed in milliseconds
function M.stop_tracking()
  M._tracking = false

  local keystrokes = M._keystroke_count
  local time_ms = 0

  if M._start_time then
    local elapsed = vim.loop.hrtime() - M._start_time
    time_ms = math.floor(elapsed / 1000000) -- Convert to ms
  end

  -- Remove keystroke handler
  if M._on_key_ns then
    vim.on_key(nil, M._on_key_ns)
    M._on_key_ns = nil
  end

  M._keystroke_count = 0
  M._start_time = nil

  return keystrokes, time_ms
end

--- Validate a challenge completion
---@param challenge table Challenge data
---@param initial string[] Initial buffer content
---@param final string[] Final buffer content
---@return table result Validation result
function M.validate(challenge, initial, final)
  local keystrokes, time_ms = M.stop_tracking()

  local result = {
    success = false,
    keystroke_count = keystrokes,
    time_ms = time_ms,
    efficiency = 0,
    error = nil,
  }

  -- Determine validation type
  local validation_type = challenge.validation_type or "exact_match"

  if validation_type == "exact_match" then
    result.success = M._validate_exact_match(challenge, final)
  elseif validation_type == "contains" then
    result.success = M._validate_contains(challenge, final)
  elseif validation_type == "function_exists" then
    result.success = M._validate_function_exists(challenge, final)
  elseif validation_type == "pattern" then
    result.success = M._validate_pattern(challenge, final)
  elseif validation_type == "different" then
    -- Just check that the content changed
    result.success = not M._content_equal(initial, final)
  elseif validation_type == "cursor_position" then
    -- Check cursor is at expected position
    result.success = M._validate_cursor_position(challenge)
  elseif validation_type == "cursor_on_char" then
    -- Check cursor is on a specific character
    result.success = M._validate_cursor_on_char(challenge, final)
  else
    result.error = "Unknown validation type: " .. validation_type
    return result
  end

  -- Calculate efficiency if successful
  if result.success then
    local par = challenge.par_keystrokes or keystrokes
    if keystrokes > 0 then
      result.efficiency = math.min(1.0, par / keystrokes)
    else
      result.efficiency = 1.0
    end
  end

  return result
end

--- Validate exact match
---@param challenge table
---@param final string[]
---@return boolean
function M._validate_exact_match(challenge, final)
  local expected = challenge.expected_buffer
  if not expected then
    return false
  end

  local expected_lines = vim.split(expected, "\n")
  return M._content_equal(expected_lines, final)
end

--- Validate contains (final must contain expected)
---@param challenge table
---@param final string[]
---@return boolean
function M._validate_contains(challenge, final)
  local expected = challenge.expected_content
  if not expected then
    return false
  end

  local content = table.concat(final, "\n")
  return content:find(expected, 1, true) ~= nil
end

--- Validate function exists
---@param challenge table
---@param final string[]
---@return boolean
function M._validate_function_exists(challenge, final)
  local func_name = challenge.function_name
  if not func_name then
    return false
  end

  local content = table.concat(final, "\n")
  -- Check for common function patterns
  local patterns = {
    "function%s+" .. func_name .. "%s*%(", -- Lua/JS: function name(
    "def%s+" .. func_name .. "%s*%(", -- Python: def name(
    "func%s+" .. func_name .. "%s*%(", -- Go: func name(
    func_name .. "%s*=%s*function", -- JS: name = function
    func_name .. "%s*:%s*function", -- JS method: name: function
    "const%s+" .. func_name .. "%s*=", -- JS const: const name =
    "let%s+" .. func_name .. "%s*=", -- JS let: let name =
  }

  for _, pattern in ipairs(patterns) do
    if content:match(pattern) then
      return true
    end
  end

  return false
end

--- Validate pattern match
---@param challenge table
---@param final string[]
---@return boolean
function M._validate_pattern(challenge, final)
  local pattern = challenge.pattern
  if not pattern then
    return false
  end

  local content = table.concat(final, "\n")
  return content:match(pattern) ~= nil
end

--- Validate cursor position (for movement challenges)
---@param challenge table
---@return boolean
function M._validate_cursor_position(challenge)
  local expected = challenge.expected_cursor
  if not expected or #expected ~= 2 then
    return false
  end

  local cursor = vim.api.nvim_win_get_cursor(0)
  -- expected is 0-indexed [row, col], cursor is 1-indexed [row, col]
  local expected_row = expected[1] + 1
  local expected_col = expected[2]

  return cursor[1] == expected_row and cursor[2] == expected_col
end

--- Validate cursor is on a specific character (for find/search challenges)
---@param challenge table
---@param final string[]
---@return boolean
function M._validate_cursor_on_char(challenge, final)
  local expected_char = challenge.expected_char
  if not expected_char then
    return false
  end

  local cursor = vim.api.nvim_win_get_cursor(0)
  local row = cursor[1]
  local col = cursor[2]

  -- Get the character at cursor position
  if row <= #final then
    local line = final[row]
    if col < #line then
      local char = line:sub(col + 1, col + 1)
      return char == expected_char
    end
  end

  return false
end

--- Check if two content arrays are equal
---@param a string[]
---@param b string[]
---@return boolean
function M._content_equal(a, b)
  if #a ~= #b then
    return false
  end

  for i, line in ipairs(a) do
    if line ~= b[i] then
      return false
    end
  end

  return true
end

--- Calculate gold reward for a challenge
---@param challenge table Challenge data
---@param efficiency number Efficiency score (0-1)
---@return number gold Gold reward
function M.calculate_reward(challenge, efficiency)
  local base_gold = challenge.gold_base or 50
  local difficulty_mult = 1 + (challenge.difficulty or 1) * 0.2
  local efficiency_mult = 0.5 + efficiency * 0.5 -- 50% base + up to 50% for efficiency

  local gold = math.floor(base_gold * difficulty_mult * efficiency_mult)
  return math.max(1, gold) -- Minimum 1 gold
end

--- Sample challenges for testing
M.sample_challenges = {
  -- Movement challenges
  {
    id = "movement_basics_1",
    name = "Jump to End",
    category = "movement",
    difficulty = 1,
    description = "Move the cursor to the end of the line using $",
    initial_buffer = "The quick brown fox jumps over the lazy dog",
    validation_type = "cursor_position",
    expected_cursor = { 0, 42 }, -- End of line (0-indexed)
    par_keystrokes = 1,
    gold_base = 25,
  },
  {
    id = "movement_word_hop",
    name = "Word Hop",
    category = "movement",
    difficulty = 1,
    description = "Move forward 5 words using 5w",
    initial_buffer = "one two three four five six seven eight",
    validation_type = "cursor_on_char",
    expected_char = "s", -- Start of "six"
    par_keystrokes = 2,
    gold_base = 25,
  },
  {
    id = "movement_find_char",
    name = "Find the X",
    category = "movement",
    difficulty = 1,
    description = "Jump to the letter 'x' using fx",
    initial_buffer = "The fox jumped over the box",
    validation_type = "cursor_on_char",
    expected_char = "x",
    par_keystrokes = 2,
    gold_base = 25,
  },
  -- Text object challenges
  {
    id = "text_object_1",
    name = "Change Inside Quotes",
    category = "text-objects",
    difficulty = 2,
    description = 'Change the text inside the quotes to "world"',
    initial_buffer = 'message = "hello"',
    expected_buffer = 'message = "world"',
    validation_type = "exact_match",
    par_keystrokes = 9, -- ci"world<Esc>
    gold_base = 50,
  },
  {
    id = "text_object_2",
    name = "Delete Inside Parens",
    category = "text-objects",
    difficulty = 2,
    description = "Delete everything inside the parentheses using di(",
    initial_buffer = "console.log(getValue());",
    expected_buffer = "console.log();",
    validation_type = "exact_match",
    par_keystrokes = 3,
    gold_base = 40,
  },
  {
    id = "delete_line_1",
    name = "Delete the Comment",
    category = "movement",
    difficulty = 1,
    description = "Delete the commented line",
    initial_buffer = [[
function hello() {
  // TODO: remove this
  console.log("hello");
}]],
    expected_buffer = [[
function hello() {
  console.log("hello");
}]],
    validation_type = "exact_match",
    par_keystrokes = 3, -- jdd
    gold_base = 30,
  },
  -- Search and replace
  {
    id = "search_replace_1",
    name = "Simple Replace",
    category = "search-replace",
    difficulty = 1,
    description = "Replace 'foo' with 'bar' using :s/foo/bar/",
    initial_buffer = "The foo is here",
    expected_buffer = "The bar is here",
    validation_type = "exact_match",
    par_keystrokes = 13,
    gold_base = 35,
  },
  {
    id = "search_replace_global",
    name = "Global Replace",
    category = "search-replace",
    difficulty = 2,
    description = "Replace all 'old' with 'new' using :%s/old/new/g",
    initial_buffer = "old value here\nanother old one\nold again",
    expected_buffer = "new value here\nanother new one\nnew again",
    validation_type = "exact_match",
    par_keystrokes = 16,
    gold_base = 50,
  },
  -- Refactoring
  {
    id = "extract_function_1",
    name = "Extract Function",
    category = "refactoring",
    difficulty = 3,
    description = "Extract the validation logic into a function called 'validateEmail'",
    initial_buffer = [[
function processForm(data) {
  if (!data.email || !data.email.includes('@')) {
    throw new Error('Invalid email');
  }
  saveData(data);
}]],
    validation_type = "function_exists",
    function_name = "validateEmail",
    par_keystrokes = 50,
    gold_base = 100,
  },
  -- LSP Navigation
  {
    id = "lsp_goto_def",
    name = "Go to Definition",
    category = "lsp-navigation",
    difficulty = 2,
    description = "Use gd to go to the definition of getUserName",
    filetype = "javascript",
    initial_buffer = [[
function getUserName(user) {
  return user.name;
}

const name = getUserName(currentUser);]],
    cursor_start = { 4, 13 }, -- On "getUserName" in the call
    validation_type = "cursor_position",
    expected_cursor = { 0, 9 }, -- At "getUserName" in the definition (line 0, col 9)
    par_keystrokes = 2,
    gold_base = 50,
  },
  -- Telescope challenges (require plugin)
  {
    id = "telescope_find_files",
    name = "Fuzzy Find File",
    category = "lsp-navigation",
    difficulty = 1,
    description = "Use Telescope to find a file",
    initial_buffer = "Use your fuzzy finder to locate any file.",
    validation_type = "different",
    par_keystrokes = 3,
    gold_base = 40,
    required_plugin = "telescope",
  },
  {
    id = "telescope_live_grep",
    name = "Search in Files",
    category = "search-replace",
    difficulty = 2,
    description = "Use Telescope live_grep to search across files",
    initial_buffer = "Use live grep to search for patterns.",
    validation_type = "different",
    par_keystrokes = 4,
    gold_base = 50,
    required_plugin = "telescope",
  },
  -- Surround challenges (require plugin)
  {
    id = "surround_change_quotes",
    name = "Change Quotes",
    category = "text-objects",
    difficulty = 2,
    description = "Change single quotes to double quotes using cs'\"",
    initial_buffer = "const msg = 'hello world';",
    expected_buffer = 'const msg = "hello world";',
    validation_type = "exact_match",
    par_keystrokes = 4,
    gold_base = 45,
    required_plugin = "nvim-surround",
  },
  -- Git operations
  {
    id = "git_status",
    name = "Git Status",
    category = "git-operations",
    difficulty = 1,
    description = "View git status using :Git or gitsigns",
    initial_buffer = "Check the current git status.",
    validation_type = "different",
    par_keystrokes = 4,
    gold_base = 35,
  },
}

--- Get a random challenge by category and difficulty
---@param category? string
---@param difficulty? number
---@return table|nil
function M.get_random_challenge(category, difficulty)
  local matching = {}

  for _, challenge in ipairs(M.sample_challenges) do
    local matches = true
    if category and challenge.category ~= category then
      matches = false
    end
    if difficulty and challenge.difficulty > difficulty then
      matches = false
    end
    if matches then
      table.insert(matching, challenge)
    end
  end

  if #matching == 0 then
    return nil
  end

  return matching[math.random(#matching)]
end

return M
