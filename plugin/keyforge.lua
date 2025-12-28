-- Keyforge.nvim - Tower defense game for learning vim keybindings
-- Auto-loads on Neovim start

if vim.g.loaded_keyforge then
  return
end
vim.g.loaded_keyforge = true

-- Lazy-load the main module only when needed
local function load_keyforge()
  return require("keyforge")
end

-- Create user commands
vim.api.nvim_create_user_command("Keyforge", function()
  load_keyforge().start()
end, { desc = "Start Keyforge tower defense game" })

vim.api.nvim_create_user_command("KeyforgeStop", function()
  load_keyforge().stop()
end, { desc = "Stop Keyforge game" })

vim.api.nvim_create_user_command("KeyforgeBuild", function()
  load_keyforge().build()
end, { desc = "Build Keyforge game binary" })

vim.api.nvim_create_user_command("KeyforgeComplete", function()
  load_keyforge().complete_challenge()
end, { desc = "Complete current Keyforge challenge" })
