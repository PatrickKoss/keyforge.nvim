-- vim: ft=lua tw=80

-- Use Lua 5.1 + LuaJIT + Neovim globals
stds.nvim = {
  read_globals = { "jit" }
}
std = "lua51+nvim"

-- Don't report unused self arguments of methods
self = false

-- Rerun tests only if their modification time changed
cache = true

-- Max line length
max_line_length = 120
max_code_line_length = 120

-- Max cyclomatic complexity
max_cyclomatic_complexity = 15

ignore = {
  "631",      -- max_line_length (handled separately)
  "212/_.*",  -- unused argument, for vars with "_" prefix
  "214",      -- used variable with unused hint ("_" prefix)
  "121",      -- setting read-only global variable 'vim'
  "122",      -- setting read-only field of global variable 'vim'
}

-- Global objects defined by Neovim
read_globals = {
  "vim",
}

-- Mutable globals
globals = {
  "vim.g",
  "vim.b",
  "vim.w",
  "vim.o",
  "vim.bo",
  "vim.wo",
  "vim.go",
  "vim.env",
  "_",
}

-- Files to exclude
exclude_files = {
  ".luarocks",
  "lua_modules",
}

-- Test file specific settings
files["tests/**/*.lua"] = {
  -- Allow test globals
  read_globals = {
    "describe",
    "it",
    "before_each",
    "after_each",
    "setup",
    "teardown",
    "pending",
    "assert",
    "spy",
    "stub",
    "mock",
  },
}

-- Plugin file specific settings
files["plugin/**/*.lua"] = {
  ignore = {
    "122",  -- Allow setting vim fields in plugin loader
  },
}
