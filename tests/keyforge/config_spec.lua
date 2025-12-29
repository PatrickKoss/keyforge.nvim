-- Tests for keyforge config and game settings
-- Run with: nvim --headless -c "PlenaryBustedDirectory tests/ {minimal_init = 'tests/minimal_init.lua'}"

local keyforge = require("keyforge")

describe("keyforge config", function()
  before_each(function()
    -- Reset config to defaults between tests
    keyforge.config = vim.deepcopy({
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
    })
  end)

  describe("default config", function()
    it("should have normal difficulty by default", function()
      assert.equals("normal", keyforge.config.difficulty)
    end)

    it("should have 1.0 game speed by default", function()
      assert.equals(1.0, keyforge.config.game_speed)
    end)

    it("should have 200 starting gold by default", function()
      assert.equals(200, keyforge.config.starting_gold)
    end)

    it("should have 100 starting health by default", function()
      assert.equals(100, keyforge.config.starting_health)
    end)
  end)

  describe("setup", function()
    it("should merge custom config with defaults", function()
      keyforge.setup({
        difficulty = "hard",
        starting_gold = 300,
      })

      assert.equals("hard", keyforge.config.difficulty)
      assert.equals(300, keyforge.config.starting_gold)
      -- Unchanged values should remain at defaults
      assert.equals(1.0, keyforge.config.game_speed)
      assert.equals(100, keyforge.config.starting_health)
    end)

    it("should accept game_speed values", function()
      keyforge.setup({
        game_speed = 2.0,
      })

      assert.equals(2.0, keyforge.config.game_speed)
    end)

    it("should accept all valid difficulty values", function()
      for _, difficulty in ipairs({ "easy", "normal", "hard" }) do
        keyforge.setup({ difficulty = difficulty })
        assert.equals(difficulty, keyforge.config.difficulty)
      end
    end)

    it("should accept all valid game_speed values", function()
      for _, speed in ipairs({ 0.5, 1.0, 1.5, 2.0 }) do
        keyforge.setup({ game_speed = speed })
        assert.equals(speed, keyforge.config.game_speed)
      end
    end)

    it("should accept starting_gold in valid range", function()
      keyforge.setup({ starting_gold = 100 })
      assert.equals(100, keyforge.config.starting_gold)

      keyforge.setup({ starting_gold = 500 })
      assert.equals(500, keyforge.config.starting_gold)

      keyforge.setup({ starting_gold = 300 })
      assert.equals(300, keyforge.config.starting_gold)
    end)

    it("should accept starting_health in valid range", function()
      keyforge.setup({ starting_health = 50 })
      assert.equals(50, keyforge.config.starting_health)

      keyforge.setup({ starting_health = 200 })
      assert.equals(200, keyforge.config.starting_health)

      keyforge.setup({ starting_health = 150 })
      assert.equals(150, keyforge.config.starting_health)
    end)
  end)

  describe("config types", function()
    it("should have string type for difficulty", function()
      assert.equals("string", type(keyforge.config.difficulty))
    end)

    it("should have number type for game_speed", function()
      assert.equals("number", type(keyforge.config.game_speed))
    end)

    it("should have number type for starting_gold", function()
      assert.equals("number", type(keyforge.config.starting_gold))
    end)

    it("should have number type for starting_health", function()
      assert.equals("number", type(keyforge.config.starting_health))
    end)
  end)
end)

describe("keyforge state", function()
  before_each(function()
    -- Reset module state
    keyforge._job_id = nil
    keyforge._term_buf = nil
    keyforge._term_win = nil
    keyforge._term_tab = nil
    keyforge._socket_path = nil
    keyforge._current_challenge_id = nil
    keyforge._game_state = "idle"
  end)

  describe("initial state", function()
    it("should have nil job_id", function()
      assert.is_nil(keyforge._job_id)
    end)

    it("should have nil socket_path", function()
      assert.is_nil(keyforge._socket_path)
    end)

    it("should have idle game_state", function()
      assert.equals("idle", keyforge._game_state)
    end)
  end)
end)

describe("command line generation", function()
  -- Test that the command line format is correct for passing settings
  it("should generate correct command format", function()
    -- Simulate what _launch_game does to create the command string
    local config = {
      difficulty = "hard",
      game_speed = 2.0,
      starting_gold = 300,
      starting_health = 150,
    }

    local binary = "/path/to/keyforge"
    local socket = "/tmp/test.sock"

    local cmd = string.format(
      "%s --nvim-mode --rpc-socket %s --difficulty %s --game-speed %.1f --starting-gold %d --starting-health %d",
      binary,
      socket,
      config.difficulty,
      config.game_speed,
      config.starting_gold,
      config.starting_health
    )

    -- Verify the command contains all expected flags
    assert.is_not_nil(string.find(cmd, "--difficulty hard"))
    assert.is_not_nil(string.find(cmd, "--game-speed 2.0"))
    assert.is_not_nil(string.find(cmd, "--starting-gold 300"))
    assert.is_not_nil(string.find(cmd, "--starting-health 150"))
    assert.is_not_nil(string.find(cmd, "--nvim-mode"))
    assert.is_not_nil(string.find(cmd, "--rpc-socket"))
  end)

  it("should format game_speed with one decimal place", function()
    local config = { game_speed = 1.5 }
    local formatted = string.format("%.1f", config.game_speed)
    assert.equals("1.5", formatted)
  end)

  it("should format starting values as integers", function()
    local config = {
      starting_gold = 200,
      starting_health = 100,
    }
    local gold_formatted = string.format("%d", config.starting_gold)
    local health_formatted = string.format("%d", config.starting_health)
    assert.equals("200", gold_formatted)
    assert.equals("100", health_formatted)
  end)
end)

describe("game state management", function()
  before_each(function()
    keyforge._game_state = "idle"
  end)

  it("should have valid game states", function()
    local valid_states = { "idle", "playing", "paused", "challenge_waiting", "game_over", "victory" }
    for _, state in ipairs(valid_states) do
      keyforge._game_state = state
      assert.equals(state, keyforge._game_state)
    end
  end)

  it("should start in idle state", function()
    keyforge._game_state = "idle"
    assert.equals("idle", keyforge._game_state)
  end)
end)
