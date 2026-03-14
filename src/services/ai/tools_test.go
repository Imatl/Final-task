package ai

import (
	"testing"
)

func TestExecuteToolUnknown(t *testing.T) {
	result := ExecuteTool(nil, "ticket-1", "nonexistent_tool", "{}")
	if result.Success {
		t.Error("unknown tool should return Success=false")
	}
	if result.Message == "" {
		t.Error("unknown tool should return error message")
	}
}

func TestExecuteRefund(t *testing.T) {
	result := executeRefund(nil, "ticket-1", `{"amount": 9.99, "reason": "double charge"}`)
	if !result.Success {
		t.Errorf("refund should succeed, got message: %s", result.Message)
	}
	if result.Data == nil {
		t.Error("refund should return data with transaction_id")
	}
}

func TestExecuteRefundZeroAmount(t *testing.T) {
	result := executeRefund(nil, "ticket-1", `{"amount": 0, "reason": "test"}`)
	if !result.Success {
		t.Error("refund with zero amount should still succeed (mock)")
	}
}

func TestExecuteRefundBadJSON(t *testing.T) {
	result := executeRefund(nil, "ticket-1", `not json`)
	if !result.Success {
		t.Error("refund with bad json should still succeed with zero amount")
	}
}

func TestExecuteResetPassword(t *testing.T) {
	result := executeResetPassword(nil, "ticket-1", "{}")
	if !result.Success {
		t.Error("reset_password should succeed")
	}
}

func TestExecuteCancelSub(t *testing.T) {
	result := executeCancelSub(nil, "ticket-1", "{}")
	if !result.Success {
		t.Error("cancel_subscription should succeed")
	}
}

func TestExecuteSendEmail(t *testing.T) {
	result := executeSendEmail(nil, "ticket-1", `{"subject": "Test", "body": "Hello"}`)
	if !result.Success {
		t.Error("send_email should succeed")
	}
}

func TestExecuteSendEmailBadJSON(t *testing.T) {
	result := executeSendEmail(nil, "ticket-1", `invalid`)
	if !result.Success {
		t.Error("send_email with bad json should still succeed")
	}
}

func TestExecuteEscalateBadJSON(t *testing.T) {
	result := executeEscalate(nil, "ticket-1", `invalid`)
	if result.Success {
		t.Error("escalate without DB should fail (GetTicket not available)")
	}
}
