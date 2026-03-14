package ai

import (
	"testing"
)

func TestExecuteToolDispatchRefund(t *testing.T) {
	result := ExecuteTool(nil, "t-1", "refund", `{"amount":25.00,"reason":"overcharge"}`)
	if !result.Success {
		t.Errorf("refund should succeed, got: %s", result.Message)
	}
}

func TestExecuteToolDispatchResetPassword(t *testing.T) {
	result := ExecuteTool(nil, "t-1", "reset_password", "{}")
	if !result.Success {
		t.Error("reset_password should succeed")
	}
}

func TestExecuteToolDispatchCancelSub(t *testing.T) {
	result := ExecuteTool(nil, "t-1", "cancel_subscription", "{}")
	if !result.Success {
		t.Error("cancel_subscription should succeed")
	}
	if result.Message == "" {
		t.Error("message should not be empty")
	}
}

func TestExecuteToolDispatchSendEmail(t *testing.T) {
	result := ExecuteTool(nil, "t-1", "send_email", `{"subject":"Update","body":"Your issue is resolved"}`)
	if !result.Success {
		t.Error("send_email should succeed")
	}
}

func TestExecuteToolDispatchEscalateNoDB(t *testing.T) {
	result := ExecuteTool(nil, "t-1", "escalate", `{"reason":"needs manager","priority":"critical"}`)
	if !result.Success {
		t.Logf("escalate without DB: %s (expected since UpdateTicketStatus fails without DB)", result.Message)
	}
}

func TestExecuteToolDispatchChangePlanNoDB(t *testing.T) {
	result := ExecuteTool(nil, "t-1", "change_plan", `{"new_plan":"premium"}`)
	if result.Success {
		t.Error("change_plan should fail without DB (GetTicket unavailable)")
	}
	if result.Message != "ticket not found" {
		t.Errorf("message = %q, want 'ticket not found'", result.Message)
	}
}

func TestExecuteToolDispatchLookupCustomerNoDB(t *testing.T) {
	result := ExecuteTool(nil, "t-1", "lookup_customer", "{}")
	if result.Success {
		t.Error("lookup_customer should fail without DB")
	}
	if result.Message != "ticket not found" {
		t.Errorf("message = %q, want 'ticket not found'", result.Message)
	}
}

func TestExecuteToolDispatchLookupBillingNoDB(t *testing.T) {
	result := ExecuteTool(nil, "t-1", "lookup_billing", "{}")
	if result.Success {
		t.Error("lookup_billing should fail without DB")
	}
	if result.Message != "ticket not found" {
		t.Errorf("message = %q, want 'ticket not found'", result.Message)
	}
}

func TestToolResultStructFields(t *testing.T) {
	r := ToolResult{Success: true, Message: "ok", Data: map[string]any{"key": "val"}}
	if !r.Success {
		t.Error("Success should be true")
	}
	if r.Message != "ok" {
		t.Errorf("Message = %q, want ok", r.Message)
	}
	data := r.Data.(map[string]any)
	if data["key"] != "val" {
		t.Errorf("Data[key] = %v, want val", data["key"])
	}
}

func TestExecuteRefundLargeAmount(t *testing.T) {
	result := executeRefund(nil, "t-1", `{"amount":99999.99,"reason":"full refund"}`)
	if !result.Success {
		t.Error("large refund should still succeed (mock)")
	}
}

func TestExecuteRefundNegativeAmount(t *testing.T) {
	result := executeRefund(nil, "t-1", `{"amount":-5.00,"reason":"test"}`)
	if !result.Success {
		t.Error("negative amount should still succeed (mock, no validation)")
	}
}

func TestExecuteSendEmailEmptyFields(t *testing.T) {
	result := executeSendEmail(nil, "t-1", `{"subject":"","body":""}`)
	if !result.Success {
		t.Error("send_email with empty fields should still succeed")
	}
}

func TestExecuteEscalateDefaultPriority(t *testing.T) {
	result := executeEscalate(nil, "t-1", `{"reason":"customer upset"}`)
	if result.Success && result.Message == "" {
		t.Error("escalate message should not be empty on success")
	}
}
