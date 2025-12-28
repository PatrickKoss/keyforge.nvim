-- Tests for challenge queue module
-- Run with: nvim --headless -c "PlenaryBustedDirectory tests/ {minimal_init = 'tests/minimal_init.lua'}"

local challenge_queue = require("keyforge.challenge_queue")

describe("challenge_queue", function()
  before_each(function()
    -- Reset state before each test
    challenge_queue.reset()
  end)

  describe("calculate_speed_bonus", function()
    it("should return 1.0 for completion at par time", function()
      local bonus = challenge_queue.calculate_speed_bonus(5000, 5000)
      assert.equals(1.0, bonus)
    end)

    it("should return 1.0 for slower than par time", function()
      local bonus = challenge_queue.calculate_speed_bonus(7000, 5000)
      assert.equals(1.0, bonus)
    end)

    it("should return bonus for faster than par time", function()
      local bonus = challenge_queue.calculate_speed_bonus(2500, 5000)
      -- 2x faster should give 1.5x bonus (formula: 1 + (ratio - 1) * 0.5)
      assert.is_true(bonus > 1.0)
      assert.is_true(bonus <= 2.0)
    end)

    it("should cap at 2.0", function()
      local bonus = challenge_queue.calculate_speed_bonus(500, 5000)
      assert.equals(2.0, bonus)
    end)

    it("should handle zero time gracefully", function()
      local bonus = challenge_queue.calculate_speed_bonus(0, 5000)
      assert.equals(1.0, bonus)
    end)

    it("should handle zero par time gracefully", function()
      local bonus = challenge_queue.calculate_speed_bonus(5000, 0)
      assert.equals(1.0, bonus)
    end)
  end)

  describe("calculate_gold", function()
    it("should calculate base gold for perfect run", function()
      local challenge = { gold_base = 50, difficulty = 1 }
      local gold = challenge_queue.calculate_gold(challenge, 1.0, 1.0)
      -- 50 * 1.25 (d1) * 1.0 (eff) * 1.0 (speed) = 62.5 -> 62
      assert.is_true(gold >= 50)
    end)

    it("should increase gold with speed bonus", function()
      local challenge = { gold_base = 50, difficulty = 1 }
      local normal = challenge_queue.calculate_gold(challenge, 1.0, 1.0)
      local fast = challenge_queue.calculate_gold(challenge, 1.0, 2.0)
      assert.is_true(fast > normal)
    end)

    it("should scale with difficulty", function()
      local challenge1 = { gold_base = 50, difficulty = 1 }
      local challenge3 = { gold_base = 50, difficulty = 3 }
      local gold1 = challenge_queue.calculate_gold(challenge1, 1.0, 1.0)
      local gold3 = challenge_queue.calculate_gold(challenge3, 1.0, 1.0)
      assert.is_true(gold3 > gold1)
    end)

    it("should reduce gold for low efficiency", function()
      local challenge = { gold_base = 50, difficulty = 1 }
      local perfect = challenge_queue.calculate_gold(challenge, 1.0, 1.0)
      local poor = challenge_queue.calculate_gold(challenge, 0.5, 1.0)
      assert.is_true(perfect > poor)
    end)

    it("should always return at least 1 gold", function()
      local challenge = { gold_base = 1, difficulty = 1 }
      local gold = challenge_queue.calculate_gold(challenge, 0.1, 1.0)
      assert.is_true(gold >= 1)
    end)
  end)

  describe("init", function()
    it("should reset state and load challenges", function()
      challenge_queue.init()
      local stats = challenge_queue.get_stats()
      assert.equals(0, stats.total_gold)
      assert.equals(0, stats.completed_count)
      assert.is_true(stats.available_count > 0)
    end)
  end)

  describe("load_available_challenges", function()
    it("should load challenges from sample_challenges", function()
      challenge_queue.load_available_challenges()
      local stats = challenge_queue.get_stats()
      assert.is_true(stats.available_count > 0)
    end)
  end)

  describe("get_next_challenge", function()
    it("should return a challenge when available", function()
      challenge_queue.init()
      local challenge = challenge_queue.get_next_challenge()
      assert.is_not_nil(challenge)
      assert.is_not_nil(challenge.id)
    end)
  end)

  describe("is_challenge_active", function()
    it("should return false when no challenge active", function()
      challenge_queue.init()
      assert.is_false(challenge_queue.is_challenge_active())
    end)
  end)

  describe("get_challenge_with_hints", function()
    it("should enrich challenge with hints", function()
      challenge_queue.init()
      local challenge = challenge_queue.get_next_challenge()
      local enriched = challenge_queue.get_challenge_with_hints(challenge)
      assert.is_not_nil(enriched.hints)
      assert.is_table(enriched.hints)
      assert.is_table(enriched.hints_display)
    end)
  end)

  describe("get_stats", function()
    it("should return stats table", function()
      challenge_queue.init()
      local stats = challenge_queue.get_stats()
      assert.is_table(stats)
      assert.is_number(stats.total_gold)
      assert.is_number(stats.completed_count)
      assert.is_number(stats.skipped_count)
      assert.is_number(stats.available_count)
    end)
  end)

  describe("reset", function()
    it("should reset all state", function()
      challenge_queue.init()
      -- Simulate some activity
      challenge_queue._state.total_gold = 100
      challenge_queue._state.completed = { {}, {} }

      challenge_queue.reset()
      local stats = challenge_queue.get_stats()
      assert.equals(0, stats.total_gold)
      assert.equals(0, stats.completed_count)
    end)
  end)

  describe("callbacks", function()
    it("should accept challenge start callback", function()
      local called = false
      challenge_queue.on_challenge_start(function()
        called = true
      end)
      assert.is_false(called) -- Just setting callback shouldn't trigger it
    end)

    it("should accept challenge complete callback", function()
      local called = false
      challenge_queue.on_challenge_complete(function()
        called = true
      end)
      assert.is_false(called) -- Just setting callback shouldn't trigger it
    end)
  end)
end)
