-- Tests for keyforge failure feedback module
-- Run with: nvim --headless -c "PlenaryBustedDirectory tests/ {minimal_init = 'tests/minimal_init.lua'}"

local failure_feedback = require("keyforge.failure_feedback")

describe("failure_feedback", function()
  -- Clean up after each test
  after_each(function()
    failure_feedback.close()
  end)

  describe("module structure", function()
    it("should export show function", function()
      assert.is_function(failure_feedback.show)
    end)

    it("should export close function", function()
      assert.is_function(failure_feedback.close)
    end)

    it("should export is_showing function", function()
      assert.is_function(failure_feedback.is_showing)
    end)
  end)

  describe("is_showing", function()
    it("should return false when no window is open", function()
      failure_feedback.close()
      assert.is_false(failure_feedback.is_showing())
    end)
  end)

  describe("close", function()
    it("should be safe to call when nothing is open", function()
      failure_feedback.close()
      failure_feedback.close() -- calling twice should be safe
      assert.is_false(failure_feedback.is_showing())
    end)

    it("should clear internal state", function()
      failure_feedback.close()
      assert.is_nil(failure_feedback._win)
      assert.is_nil(failure_feedback._buf)
      assert.is_nil(failure_feedback._callbacks)
      assert.is_nil(failure_feedback._challenge)
    end)
  end)

  describe("show", function()
    it("should accept challenge and failure_details", function()
      local challenge = {
        name = "Test Challenge",
        description = "Test description",
      }
      local failure_details = {
        message = "Test failure message",
        validation_type = "exact_match",
      }
      local callbacks = {
        on_retry = function() end,
        on_skip = function() end,
      }

      -- This should not error
      local ok = pcall(function()
        failure_feedback.show(challenge, failure_details, callbacks)
      end)
      assert.is_true(ok)

      -- Clean up
      failure_feedback.close()
    end)

    it("should handle nil failure_details", function()
      local challenge = {
        name = "Test Challenge",
        description = "Test description",
      }
      local callbacks = {
        on_retry = function() end,
        on_skip = function() end,
      }

      -- Should not error with nil failure_details
      local ok = pcall(function()
        failure_feedback.show(challenge, nil, callbacks)
      end)
      assert.is_true(ok)

      failure_feedback.close()
    end)

    it("should handle challenge with minimal fields", function()
      local challenge = {}
      local callbacks = {}

      local ok = pcall(function()
        failure_feedback.show(challenge, nil, callbacks)
      end)
      assert.is_true(ok)

      failure_feedback.close()
    end)
  end)

  describe("failure_details handling", function()
    it("should handle exact_match failure details", function()
      local challenge = { name = "Test" }
      local failure_details = {
        validation_type = "exact_match",
        message = "Buffer content does not match",
        diff_lines = { "- expected", "+ actual" },
        expected = { "expected" },
        actual = { "actual" },
      }

      local ok = pcall(function()
        failure_feedback.show(challenge, failure_details, {})
      end)
      assert.is_true(ok)

      failure_feedback.close()
    end)

    it("should handle cursor_position failure details", function()
      local challenge = { name = "Test" }
      local failure_details = {
        validation_type = "cursor_position",
        message = "Cursor at wrong position",
        expected = { row = 1, col = 5 },
        actual = { row = 1, col = 0 },
      }

      local ok = pcall(function()
        failure_feedback.show(challenge, failure_details, {})
      end)
      assert.is_true(ok)

      failure_feedback.close()
    end)

    it("should handle cursor_on_char failure details", function()
      local challenge = { name = "Test" }
      local failure_details = {
        validation_type = "cursor_on_char",
        message = "Cursor on wrong character",
        expected = "x",
        actual = "a",
      }

      local ok = pcall(function()
        failure_feedback.show(challenge, failure_details, {})
      end)
      assert.is_true(ok)

      failure_feedback.close()
    end)

    it("should handle contains failure details", function()
      local challenge = { name = "Test" }
      local failure_details = {
        validation_type = "contains",
        message = "Expected text not found",
        expected = "function test",
      }

      local ok = pcall(function()
        failure_feedback.show(challenge, failure_details, {})
      end)
      assert.is_true(ok)

      failure_feedback.close()
    end)

    it("should handle function_exists failure details", function()
      local challenge = { name = "Test" }
      local failure_details = {
        validation_type = "function_exists",
        message = "Function not found",
        expected = "myFunction",
      }

      local ok = pcall(function()
        failure_feedback.show(challenge, failure_details, {})
      end)
      assert.is_true(ok)

      failure_feedback.close()
    end)

    it("should handle different validation failure details", function()
      local challenge = { name = "Test" }
      local failure_details = {
        validation_type = "different",
        message = "Content must change from initial state",
      }

      local ok = pcall(function()
        failure_feedback.show(challenge, failure_details, {})
      end)
      assert.is_true(ok)

      failure_feedback.close()
    end)

    it("should truncate long diff lists", function()
      local challenge = { name = "Test" }
      local long_diff = {}
      for i = 1, 20 do
        table.insert(long_diff, "- line " .. i)
        table.insert(long_diff, "+ modified line " .. i)
      end

      local failure_details = {
        validation_type = "exact_match",
        message = "Content mismatch",
        diff_lines = long_diff,
      }

      local ok = pcall(function()
        failure_feedback.show(challenge, failure_details, {})
      end)
      assert.is_true(ok)

      failure_feedback.close()
    end)
  end)

  describe("hint display", function()
    it("should show challenge description as hint", function()
      local challenge = {
        name = "Test",
        description = "Use $ to jump to end of line",
      }

      local ok = pcall(function()
        failure_feedback.show(challenge, nil, {})
      end)
      assert.is_true(ok)

      failure_feedback.close()
    end)

    it("should handle very long hints", function()
      local challenge = {
        name = "Test",
        description = "This is a very long description that should be wrapped " ..
            "across multiple lines because it exceeds the normal width limit " ..
            "for hints in the failure feedback window",
      }

      local ok = pcall(function()
        failure_feedback.show(challenge, nil, {})
      end)
      assert.is_true(ok)

      failure_feedback.close()
    end)
  end)
end)
