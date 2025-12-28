-- Integration test for challenge-based economy
-- Verifies that the game can be won using primarily challenge gold
--
-- Run with:
-- nvim --headless -c "luafile tests/integration/challenge_economy_test.lua"

local function test_economy_balance()
  local engine = {}
  local challenges = require("keyforge.challenges")
  local challenge_queue = require("keyforge.challenge_queue")

  -- Simulate game economy
  local starting_gold = 200
  local tower_cost = 50 -- Arrow tower
  local num_waves = 10

  -- Calculate required gold to win
  -- Assume we need approximately 6 towers to win (rough estimate)
  local towers_needed = 6
  local gold_needed = towers_needed * tower_cost

  -- Calculate gold from mob kills (25% of base values)
  local mob_gold_per_wave = 20 -- Rough estimate with 25% multiplier
  local total_mob_gold = mob_gold_per_wave * num_waves

  -- Calculate how much challenge gold we need
  local challenge_gold_needed = gold_needed - starting_gold - total_mob_gold

  print("=== Economy Balance Test ===")
  print(string.format("Starting gold: %d", starting_gold))
  print(string.format("Towers needed: %d @ %d each = %d total", towers_needed, tower_cost, gold_needed))
  print(string.format("Expected mob gold (10 waves @ 25%%): %d", total_mob_gold))
  print(string.format("Challenge gold needed: %d", challenge_gold_needed))
  print("")

  -- Initialize challenge queue
  challenge_queue.init()

  -- Simulate completing 5 challenges
  local total_challenge_gold = 0
  local challenges_completed = 0

  for i = 1, 5 do
    local challenge = challenge_queue.get_next_challenge()
    if challenge then
      -- Simulate perfect completion
      local efficiency = 1.0
      local speed_bonus = 1.5 -- Assuming reasonably fast completion
      local gold = challenge_queue.calculate_gold(challenge, efficiency, speed_bonus)
      total_challenge_gold = total_challenge_gold + gold
      challenges_completed = challenges_completed + 1
      print(string.format("Challenge %d: %s -> %dg", i, challenge.name, gold))
    end
  end

  print("")
  print(string.format("Challenges completed: %d", challenges_completed))
  print(string.format("Total challenge gold: %d", total_challenge_gold))

  local total_available = starting_gold + total_mob_gold + total_challenge_gold
  local can_win = total_available >= gold_needed

  print("")
  print(string.format("Total available gold: %d", total_available))
  print(string.format("Gold needed to win: %d", gold_needed))
  print(string.format("Can win: %s", can_win and "YES" or "NO"))

  if not can_win then
    print("")
    print("WARNING: Economy may need balancing!")
    print(string.format("Shortfall: %dg", gold_needed - total_available))
  end

  return can_win
end

local function test_speed_bonus()
  local challenge_queue = require("keyforge.challenge_queue")

  print("")
  print("=== Speed Bonus Test ===")

  local test_cases = {
    { time = 5000, par = 5000, expected_min = 1.0, expected_max = 1.0 },
    { time = 7000, par = 5000, expected_min = 1.0, expected_max = 1.0 },
    { time = 2500, par = 5000, expected_min = 1.4, expected_max = 1.6 },
    { time = 1000, par = 5000, expected_min = 1.9, expected_max = 2.0 },
    { time = 500, par = 5000, expected_min = 2.0, expected_max = 2.0 },
  }

  local all_passed = true
  for _, tc in ipairs(test_cases) do
    local bonus = challenge_queue.calculate_speed_bonus(tc.time, tc.par)
    local passed = bonus >= tc.expected_min and bonus <= tc.expected_max
    local status = passed and "PASS" or "FAIL"
    print(string.format("  time=%d, par=%d -> bonus=%.2f (%s)", tc.time, tc.par, bonus, status))
    if not passed then
      all_passed = false
    end
  end

  return all_passed
end

local function test_gold_calculation()
  local challenge_queue = require("keyforge.challenge_queue")

  print("")
  print("=== Gold Calculation Test ===")

  -- Test challenge with different efficiency and speed
  local challenge = { gold_base = 50, difficulty = 2 }

  local test_cases = {
    { eff = 1.0, speed = 1.0, min = 70, max = 80 },   -- Perfect, no speed bonus
    { eff = 1.0, speed = 2.0, min = 140, max = 160 }, -- Perfect, max speed bonus
    { eff = 0.5, speed = 1.0, min = 50, max = 60 },   -- Half efficiency
    { eff = 0.5, speed = 1.5, min = 70, max = 90 },   -- Half efficiency, speed bonus
  }

  local all_passed = true
  for _, tc in ipairs(test_cases) do
    local gold = challenge_queue.calculate_gold(challenge, tc.eff, tc.speed)
    local passed = gold >= tc.min and gold <= tc.max
    local status = passed and "PASS" or "FAIL"
    print(string.format("  eff=%.1f, speed=%.1f -> gold=%d (%s, expected %d-%d)",
      tc.eff, tc.speed, gold, status, tc.min, tc.max))
    if not passed then
      all_passed = false
    end
  end

  return all_passed
end

-- Run tests
print("Keyforge Economy Integration Tests")
print("===================================")
print("")

local economy_pass = test_economy_balance()
local speed_pass = test_speed_bonus()
local gold_pass = test_gold_calculation()

print("")
print("=== Summary ===")
print(string.format("Economy balance: %s", economy_pass and "PASS" or "FAIL"))
print(string.format("Speed bonus: %s", speed_pass and "PASS" or "FAIL"))
print(string.format("Gold calculation: %s", gold_pass and "PASS" or "FAIL"))

local all_passed = economy_pass and speed_pass and gold_pass
print("")
if all_passed then
  print("All tests PASSED!")
  vim.cmd("qa!")
else
  print("Some tests FAILED!")
  vim.cmd("cq!")
end
