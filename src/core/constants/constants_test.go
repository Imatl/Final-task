package constants

import "testing"

func TestTicketStatuses(t *testing.T) {
	statuses := []string{
		TicketStatusOpen,
		TicketStatusInProgress,
		TicketStatusWaiting,
		TicketStatusResolved,
		TicketStatusClosed,
	}
	seen := map[string]bool{}
	for _, s := range statuses {
		if s == "" {
			t.Error("ticket status should not be empty")
		}
		if seen[s] {
			t.Errorf("duplicate ticket status: %q", s)
		}
		seen[s] = true
	}
	if len(statuses) != 5 {
		t.Errorf("expected 5 ticket statuses, got %d", len(statuses))
	}
}

func TestPriorities(t *testing.T) {
	priorities := []string{PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical}
	seen := map[string]bool{}
	for _, p := range priorities {
		if p == "" {
			t.Error("priority should not be empty")
		}
		if seen[p] {
			t.Errorf("duplicate priority: %q", p)
		}
		seen[p] = true
	}
}

func TestRoles(t *testing.T) {
	roles := []string{RoleCustomer, RoleAI, RoleAgent, RoleSystem}
	for _, r := range roles {
		if r == "" {
			t.Error("role should not be empty")
		}
	}
}

func TestActionTypes(t *testing.T) {
	actions := []string{
		ActionRefund,
		ActionChangePlan,
		ActionResetPassword,
		ActionEscalate,
		ActionSendEmail,
		ActionCancelSub,
	}
	seen := map[string]bool{}
	for _, a := range actions {
		if a == "" {
			t.Error("action type should not be empty")
		}
		if seen[a] {
			t.Errorf("duplicate action type: %q", a)
		}
		seen[a] = true
	}
}

func TestActionStatuses(t *testing.T) {
	statuses := []string{ActionStatusPending, ActionStatusApproved, ActionStatusExecuted, ActionStatusRejected}
	seen := map[string]bool{}
	for _, s := range statuses {
		if seen[s] {
			t.Errorf("duplicate action status: %q", s)
		}
		seen[s] = true
	}
}

func TestSentiments(t *testing.T) {
	sentiments := []string{SentimentPositive, SentimentNeutral, SentimentNegative, SentimentAngry}
	for _, s := range sentiments {
		if s == "" {
			t.Error("sentiment should not be empty")
		}
	}
}

func TestUrgencyLevels(t *testing.T) {
	urgencies := []string{UrgencyLow, UrgencyMedium, UrgencyHigh}
	for _, u := range urgencies {
		if u == "" {
			t.Error("urgency should not be empty")
		}
	}
}

func TestAutoExecConfidenceThreshold(t *testing.T) {
	if AutoExecConfidenceThreshold <= 0 || AutoExecConfidenceThreshold > 1.0 {
		t.Errorf("AutoExecConfidenceThreshold = %f, should be in (0, 1]", AutoExecConfidenceThreshold)
	}
	if AutoExecConfidenceThreshold != 0.85 {
		t.Errorf("AutoExecConfidenceThreshold = %f, want 0.85", AutoExecConfidenceThreshold)
	}
}
