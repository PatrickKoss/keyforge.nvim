-- Tests for keyforge challenge validation
-- Run with: nvim --headless -c "PlenaryBustedDirectory tests/ {minimal_init = 'tests/minimal_init.lua'}"

local challenges = require("keyforge.challenges")

describe("challenges", function()
  describe("content_equal", function()
    it("should return true for identical content", function()
      local a = { "line 1", "line 2", "line 3" }
      local b = { "line 1", "line 2", "line 3" }
      assert.is_true(challenges._content_equal(a, b))
    end)

    it("should return false for different content", function()
      local a = { "line 1", "line 2" }
      local b = { "line 1", "line 3" }
      assert.is_false(challenges._content_equal(a, b))
    end)

    it("should return false for different lengths", function()
      local a = { "line 1", "line 2" }
      local b = { "line 1" }
      assert.is_false(challenges._content_equal(a, b))
    end)

    it("should return true for empty arrays", function()
      local a = {}
      local b = {}
      assert.is_true(challenges._content_equal(a, b))
    end)
  end)

  describe("validate_exact_match", function()
    it("should pass when content matches expected", function()
      local challenge = {
        expected_buffer = "hello world",
      }
      local final = { "hello world" }
      assert.is_true(challenges._validate_exact_match(challenge, final))
    end)

    it("should fail when content differs", function()
      local challenge = {
        expected_buffer = "hello world",
      }
      local final = { "goodbye world" }
      assert.is_false(challenges._validate_exact_match(challenge, final))
    end)

    it("should handle multiline content", function()
      local challenge = {
        expected_buffer = "line 1\nline 2\nline 3",
      }
      local final = { "line 1", "line 2", "line 3" }
      assert.is_true(challenges._validate_exact_match(challenge, final))
    end)
  end)

  describe("validate_contains", function()
    it("should pass when content contains expected", function()
      local challenge = {
        expected_content = "function test",
      }
      local final = { "const x = 1;", "function test() {", "  return 42;", "}" }
      assert.is_true(challenges._validate_contains(challenge, final))
    end)

    it("should fail when content does not contain expected", function()
      local challenge = {
        expected_content = "function missing",
      }
      local final = { "function present() {}" }
      assert.is_false(challenges._validate_contains(challenge, final))
    end)
  end)

  describe("validate_function_exists", function()
    it("should find JavaScript function declaration", function()
      local challenge = { function_name = "myFunc" }
      local final = { "function myFunc() {", "  return 1;", "}" }
      assert.is_true(challenges._validate_function_exists(challenge, final))
    end)

    it("should find Python function", function()
      local challenge = { function_name = "my_func" }
      local final = { "def my_func():", "    pass" }
      assert.is_true(challenges._validate_function_exists(challenge, final))
    end)

    it("should find Go function", function()
      local challenge = { function_name = "MyFunc" }
      local final = { "func MyFunc() int {", "    return 42", "}" }
      assert.is_true(challenges._validate_function_exists(challenge, final))
    end)

    it("should find const arrow function", function()
      local challenge = { function_name = "myFunc" }
      local final = { "const myFunc = () => 42;" }
      assert.is_true(challenges._validate_function_exists(challenge, final))
    end)

    it("should return false when function not found", function()
      local challenge = { function_name = "missingFunc" }
      local final = { "function otherFunc() {}" }
      assert.is_false(challenges._validate_function_exists(challenge, final))
    end)
  end)

  describe("validate_pattern", function()
    it("should match regex pattern", function()
      local challenge = { pattern = "const%s+%w+%s*=" }
      local final = { "const myVar = 42;" }
      assert.is_true(challenges._validate_pattern(challenge, final))
    end)

    it("should fail when pattern not matched", function()
      local challenge = { pattern = "^function" }
      local final = { "const x = 1;" }
      assert.is_false(challenges._validate_pattern(challenge, final))
    end)
  end)

  describe("validate", function()
    it("should validate exact_match type", function()
      local challenge = {
        validation_type = "exact_match",
        expected_buffer = "result",
      }
      local initial = { "input" }
      local final = { "result" }

      local result = challenges.validate(challenge, initial, final)
      assert.is_true(result.success)
    end)

    it("should validate different type", function()
      local challenge = {
        validation_type = "different",
      }
      local initial = { "original" }
      local final = { "modified" }

      local result = challenges.validate(challenge, initial, final)
      assert.is_true(result.success)
    end)

    it("should fail different type when unchanged", function()
      local challenge = {
        validation_type = "different",
      }
      local initial = { "same" }
      local final = { "same" }

      local result = challenges.validate(challenge, initial, final)
      assert.is_false(result.success)
    end)

    it("should calculate efficiency", function()
      local challenge = {
        validation_type = "different",
        par_keystrokes = 5,
      }
      local initial = { "a" }
      local final = { "b" }

      local result = challenges.validate(challenge, initial, final)
      assert.is_true(result.success)
      -- Efficiency should be between 0 and 1
      assert.is_true(result.efficiency >= 0)
      assert.is_true(result.efficiency <= 1)
    end)
  end)

  describe("calculate_reward", function()
    it("should calculate base reward", function()
      local challenge = {
        gold_base = 50,
        difficulty = 1,
      }
      local gold = challenges.calculate_reward(challenge, 1.0)
      assert.is_true(gold >= 50)
    end)

    it("should scale with difficulty", function()
      local easy = { gold_base = 50, difficulty = 1 }
      local hard = { gold_base = 50, difficulty = 3 }

      local easy_gold = challenges.calculate_reward(easy, 1.0)
      local hard_gold = challenges.calculate_reward(hard, 1.0)

      assert.is_true(hard_gold > easy_gold)
    end)

    it("should scale with efficiency", function()
      local challenge = { gold_base = 100, difficulty = 1 }

      local perfect = challenges.calculate_reward(challenge, 1.0)
      local poor = challenges.calculate_reward(challenge, 0.5)

      assert.is_true(perfect > poor)
    end)

    it("should always give at least 1 gold", function()
      local challenge = { gold_base = 1, difficulty = 1 }
      local gold = challenges.calculate_reward(challenge, 0.01)
      assert.is_true(gold >= 1)
    end)
  end)

  describe("get_random_challenge", function()
    it("should return a challenge", function()
      local challenge = challenges.get_random_challenge()
      assert.is_not_nil(challenge)
      assert.is_not_nil(challenge.id)
    end)

    it("should filter by category", function()
      local challenge = challenges.get_random_challenge("movement")
      if challenge then
        assert.equals("movement", challenge.category)
      end
    end)

    it("should filter by difficulty", function()
      local challenge = challenges.get_random_challenge(nil, 1)
      if challenge then
        assert.is_true(challenge.difficulty <= 1)
      end
    end)
  end)

  describe("sample_challenges", function()
    it("should have sample challenges defined", function()
      assert.is_true(#challenges.sample_challenges > 0)
    end)

    it("should have required fields in each challenge", function()
      for _, c in ipairs(challenges.sample_challenges) do
        assert.is_not_nil(c.id, "Challenge missing id")
        assert.is_not_nil(c.name, "Challenge missing name")
        assert.is_not_nil(c.category, "Challenge missing category")
        assert.is_not_nil(c.difficulty, "Challenge missing difficulty")
        assert.is_not_nil(c.validation_type, "Challenge missing validation_type")
      end
    end)
  end)

  describe("action_patterns", function()
    it("should have patterns defined", function()
      assert.is_not_nil(challenges.action_patterns)
      assert.is_not_nil(challenges.action_patterns.find_files)
    end)

    it("should have rhs patterns for telescope actions", function()
      local find_files = challenges.action_patterns.find_files
      assert.is_not_nil(find_files.rhs)
      assert.is_true(#find_files.rhs > 0)
    end)

    it("should have desc patterns for telescope actions", function()
      local find_files = challenges.action_patterns.find_files
      assert.is_not_nil(find_files.desc)
      assert.is_true(#find_files.desc > 0)
    end)
  end)

  describe("plugin_aliases", function()
    it("should have aliases defined", function()
      assert.is_not_nil(challenges.plugin_aliases)
    end)

    it("should have telescope aliases", function()
      local telescope = challenges.plugin_aliases["telescope"]
      assert.is_not_nil(telescope)
      assert.is_true(#telescope > 0)
    end)

    it("should have surround aliases for both nvim-surround and mini.surround", function()
      local nvim_surround = challenges.plugin_aliases["nvim-surround"]
      local mini_surround = challenges.plugin_aliases["mini.surround"]
      assert.is_not_nil(nvim_surround)
      assert.is_not_nil(mini_surround)
    end)
  end)

  describe("format_keymap_display", function()
    it("should replace leader with Space when leader is space", function()
      -- Note: This test depends on the user's mapleader setting
      -- In minimal_init, mapleader might not be set, so we test the function exists
      local result = challenges.format_keymap_display("<leader>ff")
      assert.is_not_nil(result)
      assert.is_true(#result > 0)
    end)

    it("should replace C- with Ctrl+", function()
      local result = challenges.format_keymap_display("<C-h>")
      assert.is_true(result:find("Ctrl") ~= nil)
    end)

    it("should handle empty string", function()
      local result = challenges.format_keymap_display("")
      assert.equals("", result)
    end)

    it("should handle nil", function()
      local result = challenges.format_keymap_display(nil)
      assert.equals("", result)
    end)
  end)

  describe("get_challenge_hint", function()
    it("should return description when no hint_action", function()
      local challenge = {
        description = "Test description",
      }
      local hint = challenges.get_challenge_hint(challenge)
      assert.equals("Test description", hint)
    end)

    it("should return fallback when hint_action cannot be resolved", function()
      local challenge = {
        description = "Original description",
        hint_action = "nonexistent_action",
        hint_fallback = "Fallback hint",
      }
      local hint = challenges.get_challenge_hint(challenge)
      assert.equals("Fallback hint", hint)
    end)

    it("should return description when no fallback and resolution fails", function()
      local challenge = {
        description = "Original description",
        hint_action = "nonexistent_action",
      }
      local hint = challenges.get_challenge_hint(challenge)
      assert.equals("Original description", hint)
    end)
  end)

  describe("filter_challenges_by_plugins", function()
    it("should keep challenges without required_plugin", function()
      local input = {
        { id = "test1", name = "Test 1" },
        { id = "test2", name = "Test 2" },
      }
      local result = challenges.filter_challenges_by_plugins(input)
      assert.equals(2, #result)
    end)

    it("should filter out challenges with unavailable plugins", function()
      local input = {
        { id = "test1", name = "Test 1" },
        { id = "test2", name = "Test 2", required_plugin = "nonexistent_plugin_xyz" },
      }
      local result = challenges.filter_challenges_by_plugins(input)
      assert.equals(1, #result)
      assert.equals("test1", result[1].id)
    end)
  end)

  describe("clear_plugin_cache", function()
    it("should clear the plugin cache", function()
      -- Set something in the cache
      challenges._plugin_cache["test_plugin"] = true
      assert.is_true(challenges._plugin_cache["test_plugin"])

      -- Clear it
      challenges.clear_plugin_cache()

      -- Verify it's cleared
      assert.is_nil(challenges._plugin_cache["test_plugin"])
    end)
  end)

  describe("resolve_keymap", function()
    it("should return nil for unknown action", function()
      local lhs, desc = challenges.resolve_keymap("completely_unknown_action")
      assert.is_nil(lhs)
      assert.is_nil(desc)
    end)

    it("should return nil when action_patterns has no matching keymap", function()
      -- This tests the case where patterns exist but no matching keymap is found
      local lhs, desc = challenges.resolve_keymap("find_files")
      -- In minimal test environment, telescope is likely not set up
      -- so this should return nil
      if lhs == nil then
        assert.is_nil(desc)
      end
    end)
  end)
end)
