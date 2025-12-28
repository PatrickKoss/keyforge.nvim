-- Minimal init for running tests
-- Add the plugin to the runtime path
local plugin_path = vim.fn.fnamemodify(debug.getinfo(1, "S").source:sub(2), ":h:h")
vim.opt.runtimepath:prepend(plugin_path)

-- Load plenary if available
local ok, _ = pcall(require, "plenary")
if not ok then
  print("Warning: plenary.nvim not found, tests may not run")
end
