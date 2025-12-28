-- Tests for keyforge UI module
-- Run with: nvim --headless -c "PlenaryBustedDirectory tests/ {minimal_init = 'tests/minimal_init.lua'}"

local ui = require("keyforge.ui")

describe("ui", function()
  before_each(function()
    -- Reset UI state between tests
    ui._challenge_buf = nil
    ui._challenge_win = nil
    ui._challenge_data = nil
    ui._initial_content = nil
  end)

  after_each(function()
    -- Clean up any created buffers
    ui.close_challenge_buffer()
  end)

  describe("is_challenge_active", function()
    it("should return false when no challenge", function()
      assert.is_false(ui.is_challenge_active())
    end)
  end)

  describe("create_challenge_buffer", function()
    it("should create a buffer with initial content", function()
      local challenge = {
        id = "test_challenge",
        initial_buffer = "line 1\nline 2\nline 3",
        filetype = "text",
      }

      local buf = ui.create_challenge_buffer(challenge)

      assert.is_not_nil(buf)
      assert.is_true(vim.api.nvim_buf_is_valid(buf))

      local lines = vim.api.nvim_buf_get_lines(buf, 0, -1, false)
      assert.equals(3, #lines)
      assert.equals("line 1", lines[1])
    end)

    it("should store initial content for validation", function()
      local challenge = {
        id = "test",
        initial_buffer = "original",
      }

      ui.create_challenge_buffer(challenge)

      assert.is_not_nil(ui._initial_content)
      assert.equals("original", ui._initial_content[1])
    end)

    it("should set buffer filetype", function()
      local challenge = {
        id = "test",
        initial_buffer = "const x = 1;",
        filetype = "javascript",
      }

      local buf = ui.create_challenge_buffer(challenge)

      assert.equals("javascript", vim.bo[buf].filetype)
    end)

    it("should close existing challenge before creating new one", function()
      local challenge1 = { id = "test1", initial_buffer = "first" }
      local challenge2 = { id = "test2", initial_buffer = "second" }

      local buf1 = ui.create_challenge_buffer(challenge1)
      local buf2 = ui.create_challenge_buffer(challenge2)

      assert.is_false(vim.api.nvim_buf_is_valid(buf1))
      assert.is_true(vim.api.nvim_buf_is_valid(buf2))
    end)
  end)

  describe("get_buffer_content", function()
    it("should return empty when no challenge", function()
      local content = ui.get_buffer_content()
      assert.equals(0, #content)
    end)

    it("should return buffer content", function()
      local challenge = {
        id = "test",
        initial_buffer = "hello\nworld",
      }

      ui.create_challenge_buffer(challenge)
      local content = ui.get_buffer_content()

      assert.equals(2, #content)
      assert.equals("hello", content[1])
      assert.equals("world", content[2])
    end)
  end)

  describe("close_challenge_buffer", function()
    it("should clean up state", function()
      local challenge = { id = "test", initial_buffer = "test" }
      ui.create_challenge_buffer(challenge)

      ui.close_challenge_buffer()

      assert.is_nil(ui._challenge_buf)
      assert.is_nil(ui._challenge_win)
      assert.is_nil(ui._challenge_data)
      assert.is_nil(ui._initial_content)
    end)

    it("should handle being called when no challenge active", function()
      -- Should not error
      ui.close_challenge_buffer()
      assert.is_true(true)
    end)
  end)
end)
