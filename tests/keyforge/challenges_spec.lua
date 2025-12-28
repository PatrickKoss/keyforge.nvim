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
end)
