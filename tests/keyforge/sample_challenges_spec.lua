-- Sample challenge validation tests
-- Tests win/fail scenarios for various challenge categories
-- Run with: nvim --headless -c "PlenaryBustedDirectory tests/ {minimal_init = 'tests/minimal_init.lua'}"

local challenges = require("keyforge.challenges")

describe("sample challenge validation", function()
  -- Helper to simulate validation
  local function validate_challenge(challenge, initial, final)
    -- Stop any existing tracking first
    challenges._tracking = false
    challenges._keystroke_count = 0
    challenges._start_time = nil
    -- Start fresh tracking
    challenges._keystroke_count = 5 -- Simulate some keystrokes
    challenges._start_time = vim.loop.hrtime() - 1000000 -- 1ms ago
    challenges._tracking = true

    return challenges.validate(challenge, initial, final)
  end

  describe("movement challenges", function()
    -- 4.3.1 Add win/fail tests for 5 movement challenges

    it("should pass movement_end_of_line when cursor at end", function()
      local challenge = {
        id = "movement_end_of_line",
        validation_type = "cursor_position",
        expected_cursor = { 0, 10 },
        par_keystrokes = 1,
      }
      -- Mock cursor position
      local original_get_cursor = vim.api.nvim_win_get_cursor
      vim.api.nvim_win_get_cursor = function()
        return { 1, 10 }
      end

      local result = validate_challenge(challenge, { "hello world" }, { "hello world" })
      assert.is_true(result.success)

      vim.api.nvim_win_get_cursor = original_get_cursor
    end)

    it("should fail movement_end_of_line when cursor not at end", function()
      local challenge = {
        id = "movement_end_of_line",
        validation_type = "cursor_position",
        expected_cursor = { 0, 10 },
        par_keystrokes = 1,
      }
      local original_get_cursor = vim.api.nvim_win_get_cursor
      vim.api.nvim_win_get_cursor = function()
        return { 1, 0 }
      end

      local result = validate_challenge(challenge, { "hello world" }, { "hello world" })
      assert.is_false(result.success)

      vim.api.nvim_win_get_cursor = original_get_cursor
    end)

    it("should pass movement_find_char when on correct character", function()
      local challenge = {
        id = "movement_find_char",
        validation_type = "cursor_on_char",
        expected_char = "x",
        par_keystrokes = 2,
      }
      local original_get_cursor = vim.api.nvim_win_get_cursor
      vim.api.nvim_win_get_cursor = function()
        return { 1, 6 } -- 'x' is at column 6 (0-indexed) in "the fox"
      end

      local result = validate_challenge(challenge, { "the fox" }, { "the fox" })
      assert.is_true(result.success)

      vim.api.nvim_win_get_cursor = original_get_cursor
    end)

    it("should fail movement_find_char when on wrong character", function()
      local challenge = {
        id = "movement_find_char",
        validation_type = "cursor_on_char",
        expected_char = "x",
        par_keystrokes = 2,
      }
      local original_get_cursor = vim.api.nvim_win_get_cursor
      vim.api.nvim_win_get_cursor = function()
        return { 1, 0 }
      end

      local result = validate_challenge(challenge, { "the fox" }, { "the fox" })
      assert.is_false(result.success)

      vim.api.nvim_win_get_cursor = original_get_cursor
    end)

    it("should pass movement with pattern validation", function()
      local challenge = {
        id = "movement_search",
        validation_type = "pattern",
        pattern = "cursor",
        par_keystrokes = 5,
      }
      local result = validate_challenge(challenge, { "find the cursor here" }, { "find the cursor here" })
      assert.is_true(result.success)
    end)
  end)

  describe("text-object challenges", function()
    -- 4.3.2 Add win/fail tests for 5 text-object challenges

    it("should pass text_change_inner_word with exact match", function()
      local challenge = {
        id = "text_change_inner_word",
        validation_type = "exact_match",
        expected_buffer = "hello universe",
        par_keystrokes = 10,
      }
      local result = validate_challenge(challenge, { "hello world" }, { "hello universe" })
      assert.is_true(result.success)
    end)

    it("should fail text_change_inner_word when not matching", function()
      local challenge = {
        id = "text_change_inner_word",
        validation_type = "exact_match",
        expected_buffer = "hello universe",
        par_keystrokes = 10,
      }
      local result = validate_challenge(challenge, { "hello world" }, { "hello world" })
      assert.is_false(result.success)
    end)

    it("should pass text_delete_inner_parens", function()
      local challenge = {
        id = "text_delete_inner_parens",
        validation_type = "exact_match",
        expected_buffer = "func()",
        par_keystrokes = 3,
      }
      local result = validate_challenge(challenge, { "func(arg1, arg2)" }, { "func()" })
      assert.is_true(result.success)
    end)

    it("should pass text_change_inner_quotes", function()
      local challenge = {
        id = "text_change_inner_quotes",
        validation_type = "exact_match",
        expected_buffer = 'msg = "world"',
        par_keystrokes = 9,
      }
      local result = validate_challenge(challenge, { 'msg = "hello"' }, { 'msg = "world"' })
      assert.is_true(result.success)
    end)

    it("should pass text_delete_line", function()
      local challenge = {
        id = "text_delete_line",
        validation_type = "exact_match",
        expected_buffer = "line1\nline3",
        par_keystrokes = 2,
      }
      local result = validate_challenge(challenge, { "line1", "line2", "line3" }, { "line1", "line3" })
      assert.is_true(result.success)
    end)
  end)

  describe("lsp-navigation challenges", function()
    -- 4.3.3 Add win/fail tests for 5 LSP challenges

    it("should pass lsp_goto_definition with cursor position", function()
      local challenge = {
        id = "lsp_goto_definition",
        validation_type = "cursor_position",
        expected_cursor = { 0, 9 },
        par_keystrokes = 2,
      }
      local original_get_cursor = vim.api.nvim_win_get_cursor
      vim.api.nvim_win_get_cursor = function()
        return { 1, 9 }
      end

      local result = validate_challenge(challenge, { "function test() {}" }, { "function test() {}" })
      assert.is_true(result.success)

      vim.api.nvim_win_get_cursor = original_get_cursor
    end)

    it("should fail lsp_goto_definition when cursor wrong", function()
      local challenge = {
        id = "lsp_goto_definition",
        validation_type = "cursor_position",
        expected_cursor = { 0, 9 },
        par_keystrokes = 2,
      }
      local original_get_cursor = vim.api.nvim_win_get_cursor
      vim.api.nvim_win_get_cursor = function()
        return { 1, 0 }
      end

      local result = validate_challenge(challenge, { "function test() {}" }, { "function test() {}" })
      assert.is_false(result.success)

      vim.api.nvim_win_get_cursor = original_get_cursor
    end)

    it("should pass lsp challenge with different validation", function()
      local challenge = {
        id = "lsp_hover",
        validation_type = "different",
        par_keystrokes = 1,
      }
      local result = validate_challenge(challenge, { "original" }, { "modified" })
      assert.is_true(result.success)
    end)

    it("should pass lsp_rename with contains validation", function()
      local challenge = {
        id = "lsp_rename",
        validation_type = "contains",
        expected_content = "newName",
        par_keystrokes = 10,
      }
      local result = validate_challenge(challenge, { "oldName" }, { "const newName = 1;" })
      assert.is_true(result.success)
    end)

    it("should fail lsp_rename when content not found", function()
      local challenge = {
        id = "lsp_rename",
        validation_type = "contains",
        expected_content = "newName",
        par_keystrokes = 10,
      }
      local result = validate_challenge(challenge, { "oldName" }, { "oldName" })
      assert.is_false(result.success)
    end)
  end)

  describe("search-replace challenges", function()
    -- 4.3.4 Add win/fail tests for 3 search-replace challenges

    it("should pass search_replace_line with exact match", function()
      local challenge = {
        id = "search_replace_line",
        validation_type = "exact_match",
        expected_buffer = "The bar is here",
        par_keystrokes = 13,
      }
      local result = validate_challenge(challenge, { "The foo is here" }, { "The bar is here" })
      assert.is_true(result.success)
    end)

    it("should pass search_replace_global with multiline", function()
      local challenge = {
        id = "search_replace_global",
        validation_type = "exact_match",
        expected_buffer = "new value\nnew again",
        par_keystrokes = 16,
      }
      local result = validate_challenge(challenge, { "old value", "old again" }, { "new value", "new again" })
      assert.is_true(result.success)
    end)

    it("should fail search_replace when incomplete", function()
      local challenge = {
        id = "search_replace_global",
        validation_type = "exact_match",
        expected_buffer = "new value\nnew again",
        par_keystrokes = 16,
      }
      local result = validate_challenge(challenge, { "old value", "old again" }, { "new value", "old again" })
      assert.is_false(result.success)
    end)
  end)

  describe("refactoring challenges", function()
    -- 4.3.5 Add win/fail tests for 3 refactoring challenges

    it("should pass extract_function with function_exists", function()
      local challenge = {
        id = "refactor_extract_function",
        validation_type = "function_exists",
        function_name = "validateEmail",
        par_keystrokes = 50,
      }
      local result = validate_challenge(challenge, { "// inline code" }, { "function validateEmail(email) {", "  return email.includes('@');", "}" })
      assert.is_true(result.success)
    end)

    it("should fail extract_function when function missing", function()
      local challenge = {
        id = "refactor_extract_function",
        validation_type = "function_exists",
        function_name = "validateEmail",
        par_keystrokes = 50,
      }
      local result = validate_challenge(challenge, { "// inline code" }, { "// still inline code" })
      assert.is_false(result.success)
    end)

    it("should pass refactor_join_lines", function()
      local challenge = {
        id = "refactor_join_lines",
        validation_type = "exact_match",
        expected_buffer = "line1 line2",
        par_keystrokes = 1,
      }
      local result = validate_challenge(challenge, { "line1", "line2" }, { "line1 line2" })
      assert.is_true(result.success)
    end)
  end)

  describe("git-operations challenges", function()
    -- 4.3.6 Add win/fail tests for 3 git-operations challenges

    it("should pass git challenge with different validation", function()
      local challenge = {
        id = "git_stage_hunk",
        validation_type = "different",
        par_keystrokes = 3,
      }
      local result = validate_challenge(challenge, { "unstaged" }, { "staged" })
      assert.is_true(result.success)
    end)

    it("should fail git challenge when unchanged", function()
      local challenge = {
        id = "git_stage_hunk",
        validation_type = "different",
        par_keystrokes = 3,
      }
      local result = validate_challenge(challenge, { "unchanged" }, { "unchanged" })
      assert.is_false(result.success)
    end)

    it("should pass git_preview_hunk with different", function()
      local challenge = {
        id = "git_preview_hunk",
        validation_type = "different",
        par_keystrokes = 3,
      }
      local result = validate_challenge(challenge, { "before" }, { "after preview" })
      assert.is_true(result.success)
    end)
  end)

  describe("window-management challenges", function()
    -- 4.3.7 Add win/fail tests for 2 window-management challenges

    it("should pass window challenge with different validation", function()
      local challenge = {
        id = "window_split_horizontal",
        validation_type = "different",
        par_keystrokes = 3,
      }
      local result = validate_challenge(challenge, { "single window" }, { "split window" })
      assert.is_true(result.success)
    end)

    it("should fail window challenge when unchanged", function()
      local challenge = {
        id = "window_go_left",
        validation_type = "different",
        par_keystrokes = 1,
      }
      local result = validate_challenge(challenge, { "same" }, { "same" })
      assert.is_false(result.success)
    end)
  end)

  describe("buffer-management challenges", function()
    -- 4.3.8 Add win/fail tests for 2 buffer-management challenges

    it("should pass buffer challenge with different validation", function()
      local challenge = {
        id = "buffer_next",
        validation_type = "different",
        par_keystrokes = 1,
      }
      local result = validate_challenge(challenge, { "buffer1" }, { "buffer2" })
      assert.is_true(result.success)
    end)

    it("should pass buffer_save with contains", function()
      local challenge = {
        id = "buffer_save",
        validation_type = "contains",
        expected_content = "saved",
        par_keystrokes = 3,
      }
      local result = validate_challenge(challenge, { "unsaved" }, { "content saved successfully" })
      assert.is_true(result.success)
    end)
  end)

  describe("folding challenges", function()
    -- 4.3.9 Add win/fail tests for 2 folding challenges

    it("should pass fold challenge with different validation", function()
      local challenge = {
        id = "fold_toggle",
        validation_type = "different",
        par_keystrokes = 1,
      }
      local result = validate_challenge(challenge, { "unfolded" }, { "folded" })
      assert.is_true(result.success)
    end)

    it("should fail fold challenge when unchanged", function()
      local challenge = {
        id = "fold_open",
        validation_type = "different",
        par_keystrokes = 1,
      }
      local result = validate_challenge(challenge, { "same state" }, { "same state" })
      assert.is_false(result.success)
    end)
  end)

  describe("quickfix challenges", function()
    -- 4.3.10 Add win/fail tests for 2 quickfix challenges

    it("should pass quickfix challenge with different validation", function()
      local challenge = {
        id = "quickfix_next",
        validation_type = "different",
        par_keystrokes = 2,
      }
      local result = validate_challenge(challenge, { "item1" }, { "item2" })
      assert.is_true(result.success)
    end)

    it("should pass quickfix_open with contains", function()
      local challenge = {
        id = "quickfix_open",
        validation_type = "contains",
        expected_content = "quickfix",
        par_keystrokes = 3,
      }
      local result = validate_challenge(challenge, { "closed" }, { "quickfix list open" })
      assert.is_true(result.success)
    end)
  end)

  describe("telescope challenges", function()
    -- 4.3.11 Add win/fail tests for 2 telescope challenges

    it("should pass telescope challenge with different validation", function()
      local challenge = {
        id = "telescope_find_files",
        validation_type = "different",
        par_keystrokes = 3,
        required_plugin = "telescope",
      }
      local result = validate_challenge(challenge, { "no file" }, { "file selected" })
      assert.is_true(result.success)
    end)

    it("should fail telescope challenge when unchanged", function()
      local challenge = {
        id = "telescope_buffers",
        validation_type = "different",
        par_keystrokes = 3,
        required_plugin = "telescope",
      }
      local result = validate_challenge(challenge, { "same buffer" }, { "same buffer" })
      assert.is_false(result.success)
    end)
  end)

  describe("surround challenges", function()
    -- 4.3.12 Add win/fail tests for 2 surround challenges

    it("should pass surround_change_quotes with exact match", function()
      local challenge = {
        id = "surround_change_quotes",
        validation_type = "exact_match",
        expected_buffer = 'const msg = "hello";',
        par_keystrokes = 4,
        required_plugin = "nvim-surround",
      }
      local result = validate_challenge(challenge, { "const msg = 'hello';" }, { 'const msg = "hello";' })
      assert.is_true(result.success)
    end)

    it("should fail surround_change_quotes when unchanged", function()
      local challenge = {
        id = "surround_change_quotes",
        validation_type = "exact_match",
        expected_buffer = 'const msg = "hello";',
        par_keystrokes = 4,
        required_plugin = "nvim-surround",
      }
      local result = validate_challenge(challenge, { "const msg = 'hello';" }, { "const msg = 'hello';" })
      assert.is_false(result.success)
    end)
  end)

  describe("timeout scenarios", function()
    -- 4.6.4 Test challenge timeout scenarios

    it("should track time elapsed during challenge", function()
      local challenge = {
        id = "timed_challenge",
        validation_type = "different",
        par_keystrokes = 5,
      }

      -- Start tracking
      challenges.start_tracking()

      -- Simulate some time passing (at least a tiny bit)
      local start = vim.loop.hrtime()
      while vim.loop.hrtime() - start < 1000000 do
        -- Wait ~1ms
      end

      local result = challenges.validate(challenge, { "before" }, { "after" })

      assert.is_true(result.success)
      assert.is_true(result.time_ms >= 0)
    end)

    it("should calculate efficiency based on keystrokes", function()
      local challenge = {
        id = "efficiency_test",
        validation_type = "different",
        par_keystrokes = 10,
      }

      -- Manually set keystroke count
      challenges._keystroke_count = 10
      challenges._start_time = vim.loop.hrtime()
      challenges._tracking = true

      local result = challenges.validate(challenge, { "before" }, { "after" })

      assert.is_true(result.success)
      assert.equals(1.0, result.efficiency) -- Par = actual, so 100% efficiency
    end)

    it("should give lower efficiency for more keystrokes", function()
      local challenge = {
        id = "efficiency_test",
        validation_type = "different",
        par_keystrokes = 5,
      }

      -- Manually set keystroke count higher than par
      challenges._keystroke_count = 10
      challenges._start_time = vim.loop.hrtime()
      challenges._tracking = true

      local result = challenges.validate(challenge, { "before" }, { "after" })

      assert.is_true(result.success)
      assert.equals(0.5, result.efficiency) -- 5/10 = 50% efficiency
    end)
  end)
end)
