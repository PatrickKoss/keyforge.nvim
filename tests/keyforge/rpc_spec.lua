-- Tests for keyforge RPC module
-- Run with: nvim --headless -c "PlenaryBustedDirectory tests/ {minimal_init = 'tests/minimal_init.lua'}"

local rpc = require("keyforge.rpc")

describe("rpc", function()
  before_each(function()
    -- Reset RPC state between tests
    rpc._socket = nil
    rpc._connected = false
    rpc._socket_path = nil
    rpc._buffer = ""
    rpc._pending_requests = {}
    rpc._request_id = 0
    rpc._handlers = {}
  end)

  describe("is_connected", function()
    it("should return false when not connected", function()
      assert.is_false(rpc.is_connected())
    end)

    it("should return true when connected", function()
      -- Manually simulate connected state (since actual socket connect is async)
      rpc._connected = true
      rpc._socket = {} -- Mock socket object
      assert.is_true(rpc.is_connected())
    end)
  end)

  describe("disconnect", function()
    it("should clear connected state", function()
      rpc._connected = true
      rpc._socket = { is_closing = function() return false end, read_stop = function() end, close = function() end }
      rpc.disconnect()
      assert.is_false(rpc._connected)
    end)

    it("should call pending callbacks with error", function()
      local callback_called = false
      local received_error = nil

      rpc._connected = true
      rpc._socket = { is_closing = function() return false end, read_stop = function() end, close = function() end }
      rpc._pending_requests[1] = function(err, _result)
        callback_called = true
        received_error = err
      end

      rpc.disconnect()

      assert.is_true(callback_called)
      assert.is_not_nil(received_error)
    end)

    it("should clear socket path", function()
      rpc._connected = true
      rpc._socket_path = "/tmp/test.sock"
      rpc._socket = { is_closing = function() return false end, read_stop = function() end, close = function() end }
      rpc.disconnect()
      assert.is_nil(rpc._socket_path)
    end)
  end)

  describe("on", function()
    it("should register handler", function()
      local handler = function() end
      rpc.on("test_method", handler)
      assert.equals(handler, rpc._handlers["test_method"])
    end)
  end)

  describe("_handle_message", function()
    it("should call handler for notification", function()
      local received_params = nil
      rpc.on("test_notify", function(params)
        received_params = params
      end)

      rpc._handle_message({
        jsonrpc = "2.0",
        method = "test_notify",
        params = { foo = "bar" },
      })

      assert.is_not_nil(received_params)
      assert.equals("bar", received_params.foo)
    end)

    it("should resolve pending request on response", function()
      local callback_called = false
      local received_result = nil

      rpc._pending_requests[42] = function(_err, result)
        callback_called = true
        received_result = result
      end

      rpc._handle_message({
        jsonrpc = "2.0",
        id = 42,
        result = { success = true },
      })

      assert.is_true(callback_called)
      assert.is_not_nil(received_result)
      assert.is_true(received_result.success)
    end)

    it("should pass error to callback on error response", function()
      local received_error = nil

      rpc._pending_requests[43] = function(err, _result)
        received_error = err
      end

      rpc._handle_message({
        jsonrpc = "2.0",
        id = 43,
        error = { code = -32600, message = "Invalid Request" },
      })

      assert.is_not_nil(received_error)
      assert.equals(-32600, received_error.code)
    end)
  end)
end)
