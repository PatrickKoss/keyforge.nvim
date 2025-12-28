-- Tests for keyforge RPC module
-- Run with: nvim --headless -c "PlenaryBustedDirectory tests/ {minimal_init = 'tests/minimal_init.lua'}"

local rpc = require("keyforge.rpc")

describe("rpc", function()
  before_each(function()
    -- Reset RPC state between tests
    rpc._job_id = nil
    rpc._buffer = ""
    rpc._pending_requests = {}
    rpc._request_id = 0
  end)

  describe("is_connected", function()
    it("should return false when not connected", function()
      assert.is_false(rpc.is_connected())
    end)

    it("should return true when connected", function()
      rpc.connect(123)
      assert.is_true(rpc.is_connected())
    end)
  end)

  describe("connect", function()
    it("should set job_id", function()
      rpc.connect(456)
      assert.equals(456, rpc._job_id)
    end)

    it("should reset buffer", function()
      rpc._buffer = "leftover"
      rpc.connect(789)
      assert.equals("", rpc._buffer)
    end)

    it("should clear pending requests", function()
      rpc._pending_requests[1] = function() end
      rpc.connect(101)
      assert.equals(0, vim.tbl_count(rpc._pending_requests))
    end)
  end)

  describe("disconnect", function()
    it("should clear job_id", function()
      rpc.connect(123)
      rpc.disconnect()
      assert.is_nil(rpc._job_id)
    end)

    it("should call pending callbacks with error", function()
      local callback_called = false
      local received_error = nil

      rpc._job_id = 123
      rpc._pending_requests[1] = function(err, result)
        callback_called = true
        received_error = err
      end

      rpc.disconnect()

      assert.is_true(callback_called)
      assert.is_not_nil(received_error)
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

      rpc._pending_requests[42] = function(err, result)
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

      rpc._pending_requests[43] = function(err, result)
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
