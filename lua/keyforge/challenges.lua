--- Challenge validation and scoring for Keyforge
local M = {}

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
    validation_type = "different",
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
    validation_type = "different",
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
    validation_type = "different",
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
    description = "Use gd to go to the definition",
    initial_buffer = [[
function getUserName(user) {
  return user.name;
}

const name = getUserName(currentUser);]],
    validation_type = "different",
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
