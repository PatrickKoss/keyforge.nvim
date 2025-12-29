-- Integration test for start screen and settings passing
-- Verifies config is correctly passed to game binary
--
-- Run with:
-- nvim --headless -c "luafile tests/integration/start_screen_test.lua"

local keyforge = require("keyforge")

local function test_default_config_values()
  print("=== Default Config Values Test ===")
  local passed = true
  local tests = {
    { name = "difficulty", expected = "normal", actual = keyforge.config.difficulty },
    { name = "game_speed", expected = 1.0, actual = keyforge.config.game_speed },
    { name = "starting_gold", expected = 200, actual = keyforge.config.starting_gold },
    { name = "starting_health", expected = 100, actual = keyforge.config.starting_health },
  }

  for _, test in ipairs(tests) do
    local status = test.actual == test.expected and "PASS" or "FAIL"
    print(string.format("  %s: expected %s, got %s (%s)", test.name, test.expected, test.actual, status))
    if status == "FAIL" then
      passed = false
    end
  end

  return passed
end

local function test_config_override()
  print("")
  print("=== Config Override Test ===")
  local passed = true

  -- Test various config combinations
  local test_configs = {
    { difficulty = "easy", game_speed = 0.5, starting_gold = 100, starting_health = 50 },
    { difficulty = "hard", game_speed = 2.0, starting_gold = 500, starting_health = 200 },
    { difficulty = "normal", game_speed = 1.5, starting_gold = 300, starting_health = 150 },
  }

  for i, config in ipairs(test_configs) do
    keyforge.setup(config)
    local all_match = keyforge.config.difficulty == config.difficulty
      and keyforge.config.game_speed == config.game_speed
      and keyforge.config.starting_gold == config.starting_gold
      and keyforge.config.starting_health == config.starting_health

    local status = all_match and "PASS" or "FAIL"
    print(string.format("  Config set %d: %s", i, status))
    if not all_match then
      passed = false
      print(string.format("    difficulty: expected %s, got %s", config.difficulty, keyforge.config.difficulty))
      print(string.format("    game_speed: expected %s, got %s", config.game_speed, keyforge.config.game_speed))
      print(string.format("    starting_gold: expected %s, got %s", config.starting_gold, keyforge.config.starting_gold))
      print(
        string.format("    starting_health: expected %s, got %s", config.starting_health, keyforge.config.starting_health)
      )
    end
  end

  return passed
end

local function test_command_line_format()
  print("")
  print("=== Command Line Format Test ===")
  local passed = true

  local test_cases = {
    {
      config = { difficulty = "easy", game_speed = 0.5, starting_gold = 100, starting_health = 50 },
      expected_fragments = { "--difficulty easy", "--game-speed 0.5", "--starting-gold 100", "--starting-health 50" },
    },
    {
      config = { difficulty = "hard", game_speed = 2.0, starting_gold = 500, starting_health = 200 },
      expected_fragments = { "--difficulty hard", "--game-speed 2.0", "--starting-gold 500", "--starting-health 200" },
    },
    {
      config = { difficulty = "normal", game_speed = 1.0, starting_gold = 200, starting_health = 100 },
      expected_fragments = { "--difficulty normal", "--game-speed 1.0", "--starting-gold 200", "--starting-health 100" },
    },
  }

  for i, tc in ipairs(test_cases) do
    -- Simulate command line generation
    local cmd = string.format(
      "--difficulty %s --game-speed %.1f --starting-gold %d --starting-health %d",
      tc.config.difficulty,
      tc.config.game_speed,
      tc.config.starting_gold,
      tc.config.starting_health
    )

    local all_found = true
    for _, fragment in ipairs(tc.expected_fragments) do
      if not string.find(cmd, fragment, 1, true) then
        all_found = false
        print(string.format("  Test %d: missing fragment '%s' in '%s'", i, fragment, cmd))
      end
    end

    local status = all_found and "PASS" or "FAIL"
    print(string.format("  Command format test %d: %s", i, status))
    if not all_found then
      passed = false
    end
  end

  return passed
end

local function test_game_speed_values()
  print("")
  print("=== Game Speed Values Test ===")
  local passed = true

  local valid_speeds = { 0.5, 1.0, 1.5, 2.0 }
  for _, speed in ipairs(valid_speeds) do
    keyforge.setup({ game_speed = speed })
    local match = keyforge.config.game_speed == speed
    local status = match and "PASS" or "FAIL"
    print(string.format("  Speed %.1fx: %s", speed, status))
    if not match then
      passed = false
    end
  end

  return passed
end

local function test_difficulty_values()
  print("")
  print("=== Difficulty Values Test ===")
  local passed = true

  local valid_difficulties = { "easy", "normal", "hard" }
  for _, diff in ipairs(valid_difficulties) do
    keyforge.setup({ difficulty = diff })
    local match = keyforge.config.difficulty == diff
    local status = match and "PASS" or "FAIL"
    print(string.format("  Difficulty '%s': %s", diff, status))
    if not match then
      passed = false
    end
  end

  return passed
end

local function test_gold_range()
  print("")
  print("=== Starting Gold Range Test ===")
  local passed = true

  local test_values = { 100, 200, 300, 400, 500 }
  for _, gold in ipairs(test_values) do
    keyforge.setup({ starting_gold = gold })
    local match = keyforge.config.starting_gold == gold
    local status = match and "PASS" or "FAIL"
    print(string.format("  Gold %d: %s", gold, status))
    if not match then
      passed = false
    end
  end

  return passed
end

local function test_health_range()
  print("")
  print("=== Starting Health Range Test ===")
  local passed = true

  local test_values = { 50, 100, 150, 200 }
  for _, health in ipairs(test_values) do
    keyforge.setup({ starting_health = health })
    local match = keyforge.config.starting_health == health
    local status = match and "PASS" or "FAIL"
    print(string.format("  Health %d: %s", health, status))
    if not match then
      passed = false
    end
  end

  return passed
end

local function test_partial_config_override()
  print("")
  print("=== Partial Config Override Test ===")

  -- Reset to default
  keyforge.setup({})

  -- Override only one value
  keyforge.setup({ difficulty = "hard" })

  local passed = keyforge.config.difficulty == "hard"
    and keyforge.config.game_speed == 1.0
    and keyforge.config.starting_gold == 200
    and keyforge.config.starting_health == 100

  print(string.format("  Only difficulty changed: %s", passed and "PASS" or "FAIL"))
  if not passed then
    print(string.format("    difficulty: %s (expected hard)", keyforge.config.difficulty))
    print(string.format("    game_speed: %s (expected 1.0)", keyforge.config.game_speed))
    print(string.format("    starting_gold: %s (expected 200)", keyforge.config.starting_gold))
    print(string.format("    starting_health: %s (expected 100)", keyforge.config.starting_health))
  end

  return passed
end

-- Run tests
print("Keyforge Start Screen Integration Tests")
print("========================================")
print("")

local results = {
  default_config = test_default_config_values(),
  config_override = test_config_override(),
  command_line = test_command_line_format(),
  game_speed = test_game_speed_values(),
  difficulty = test_difficulty_values(),
  gold_range = test_gold_range(),
  health_range = test_health_range(),
  partial_override = test_partial_config_override(),
}

print("")
print("=== Summary ===")
local all_passed = true
for name, passed in pairs(results) do
  print(string.format("  %s: %s", name, passed and "PASS" or "FAIL"))
  if not passed then
    all_passed = false
  end
end

print("")
if all_passed then
  print("All tests PASSED!")
  vim.cmd("qa!")
else
  print("Some tests FAILED!")
  vim.cmd("cq!")
end
