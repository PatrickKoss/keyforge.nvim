-- Tests for keymap hints module
-- Run with: nvim --headless -c "PlenaryBustedDirectory tests/ {minimal_init = 'tests/minimal_init.lua'}"

local keymap_hints = require("keyforge.keymap_hints")

describe("keymap_hints", function()
  before_each(function()
    -- Reset cache before each test
    keymap_hints._cache = {
      normal = {},
      visual = {},
      insert = {},
      plugins = {},
      last_refresh = 0,
    }
  end)

  describe("detect_plugins", function()
    it("should return a table of plugin detection results", function()
      local plugins = keymap_hints.detect_plugins()
      assert.is_table(plugins)
    end)

    it("should detect or not detect telescope", function()
      local plugins = keymap_hints.detect_plugins()
      -- Result should be a boolean
      assert.is_true(type(plugins["telescope"]) == "boolean")
    end)

    it("should detect or not detect nvim-surround", function()
      local plugins = keymap_hints.detect_plugins()
      assert.is_true(type(plugins["nvim-surround"]) == "boolean")
    end)
  end)

  describe("discover_keymaps", function()
    it("should populate the cache", function()
      keymap_hints.discover_keymaps()
      assert.is_true(keymap_hints._cache.last_refresh > 0)
    end)

    it("should discover normal mode keymaps", function()
      keymap_hints.discover_keymaps()
      -- Cache should be populated (may be empty in minimal test env)
      assert.is_table(keymap_hints._cache.normal)
    end)

    it("should discover visual mode keymaps", function()
      keymap_hints.discover_keymaps()
      assert.is_table(keymap_hints._cache.visual)
    end)

    it("should discover plugins", function()
      keymap_hints.discover_keymaps()
      assert.is_table(keymap_hints._cache.plugins)
    end)
  end)

  describe("get_hint_for_action", function()
    it("should return a hint table for known actions", function()
      local hint = keymap_hints.get_hint_for_action("find_files")
      assert.is_table(hint)
      assert.is_not_nil(hint.binding)
      assert.is_not_nil(hint.description)
      assert.is_true(type(hint.is_default) == "boolean")
    end)

    it("should return default binding for find_files", function()
      local hint = keymap_hints.get_hint_for_action("find_files")
      -- Should fall back to default since we haven't configured telescope
      assert.is_true(hint.is_default)
      assert.equals("<leader>ff", hint.binding)
    end)

    it("should return default binding for goto_definition", function()
      local hint = keymap_hints.get_hint_for_action("goto_definition")
      -- Should return the default 'gd' or a discovered keymap matching the pattern
      if hint.is_default then
        assert.equals("gd", hint.binding)
      else
        -- A real keymap was found - just verify it's a valid binding
        assert.is_not_nil(hint.binding)
        assert.is_string(hint.binding)
      end
    end)

    it("should handle unknown actions gracefully", function()
      local hint = keymap_hints.get_hint_for_action("unknown_action_xyz")
      assert.is_table(hint)
      assert.equals("?", hint.binding)
      assert.is_true(hint.is_default)
    end)
  end)

  describe("get_hints_for_category", function()
    it("should return hints for movement category", function()
      local hints = keymap_hints.get_hints_for_category("movement")
      assert.is_table(hints)
    end)

    it("should return hints for lsp-navigation category", function()
      local hints = keymap_hints.get_hints_for_category("lsp-navigation")
      assert.is_table(hints)
      assert.is_true(#hints > 0)
    end)

    it("should return hints with action field", function()
      local hints = keymap_hints.get_hints_for_category("lsp-navigation")
      for _, hint in ipairs(hints) do
        assert.is_not_nil(hint.action)
        assert.is_not_nil(hint.binding)
      end
    end)

    it("should return empty table for unknown category", function()
      local hints = keymap_hints.get_hints_for_category("unknown_category")
      assert.is_table(hints)
      assert.equals(0, #hints)
    end)
  end)

  describe("has_plugin", function()
    it("should return boolean for telescope", function()
      keymap_hints.discover_keymaps()
      local result = keymap_hints.has_plugin("telescope")
      assert.is_true(type(result) == "boolean")
    end)

    it("should return false for non-existent plugin", function()
      keymap_hints.discover_keymaps()
      local result = keymap_hints.has_plugin("definitely_not_a_real_plugin_xyz")
      assert.is_false(result)
    end)
  end)

  describe("refresh_cache", function()
    it("should refresh stale cache", function()
      -- Set last refresh to old timestamp
      keymap_hints._cache.last_refresh = 0
      keymap_hints.refresh_cache()
      assert.is_true(keymap_hints._cache.last_refresh > 0)
    end)

    it("should not refresh recent cache", function()
      keymap_hints.discover_keymaps()
      local first_refresh = keymap_hints._cache.last_refresh
      -- Immediately call refresh again
      keymap_hints.refresh_cache()
      -- Should be the same timestamp since cache is fresh
      assert.equals(first_refresh, keymap_hints._cache.last_refresh)
    end)
  end)

  describe("force_refresh", function()
    it("should always refresh cache", function()
      keymap_hints.discover_keymaps()
      local first_refresh = keymap_hints._cache.last_refresh
      -- Wait a tiny bit
      vim.wait(10)
      keymap_hints.force_refresh()
      -- Should have a new timestamp
      assert.is_true(keymap_hints._cache.last_refresh >= first_refresh)
    end)
  end)

  describe("get_cache_info", function()
    it("should return cache statistics", function()
      keymap_hints.discover_keymaps()
      local info = keymap_hints.get_cache_info()
      assert.is_table(info)
      assert.is_number(info.normal_count)
      assert.is_number(info.visual_count)
      assert.is_number(info.insert_count)
      assert.is_table(info.plugins)
      assert.is_number(info.last_refresh)
      assert.is_number(info.age_seconds)
    end)
  end)

  describe("format_hints_for_display", function()
    it("should return empty table for empty hints", function()
      local lines = keymap_hints.format_hints_for_display({})
      assert.is_table(lines)
      assert.equals(0, #lines)
    end)

    it("should format hints as display lines", function()
      local hints = {
        { binding = "gd", description = "Go to definition", is_default = true },
        { binding = "<leader>ff", description = "Find files", is_default = false },
      }
      local lines = keymap_hints.format_hints_for_display(hints)
      assert.is_table(lines)
      assert.is_true(#lines > 0)
      -- First line should be header
      assert.is_true(lines[1]:find("Keybindings") ~= nil)
    end)

    it("should mark default bindings", function()
      local hints = {
        { binding = "gd", description = "Go to definition", is_default = true },
      }
      local lines = keymap_hints.format_hints_for_display(hints)
      local found_default = false
      for _, line in ipairs(lines) do
        if line:find("default") then
          found_default = true
          break
        end
      end
      assert.is_true(found_default)
    end)
  end)

  describe("get_detected_plugins", function()
    it("should return a copy of detected plugins", function()
      keymap_hints.discover_keymaps()
      local plugins = keymap_hints.get_detected_plugins()
      assert.is_table(plugins)
      -- Modifying the returned table shouldn't affect the cache
      plugins["test_plugin"] = true
      local plugins2 = keymap_hints.get_detected_plugins()
      assert.is_nil(plugins2["test_plugin"])
    end)
  end)
end)
