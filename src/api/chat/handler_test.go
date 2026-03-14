package chat

import (
	"testing"

	"supportflow/core/constants"
)

func TestTruncateShortString(t *testing.T) {
	got := truncate("hello", 10)
	if got != "hello" {
		t.Errorf("truncate('hello', 10) = %q, want 'hello'", got)
	}
}

func TestTruncateExactLength(t *testing.T) {
	got := truncate("12345", 5)
	if got != "12345" {
		t.Errorf("truncate('12345', 5) = %q, want '12345'", got)
	}
}

func TestTruncateLongString(t *testing.T) {
	got := truncate("hello world", 5)
	if got != "hello..." {
		t.Errorf("truncate('hello world', 5) = %q, want 'hello...'", got)
	}
}

func TestTruncateEmpty(t *testing.T) {
	got := truncate("", 5)
	if got != "" {
		t.Errorf("truncate('', 5) = %q, want ''", got)
	}
}

func TestTruncateZeroMax(t *testing.T) {
	got := truncate("hello", 0)
	if got != "..." {
		t.Errorf("truncate('hello', 0) = %q, want '...'", got)
	}
}

func TestMapUrgencyToPriorityHigh(t *testing.T) {
	got := mapUrgencyToPriority(constants.UrgencyHigh)
	if got != constants.PriorityHigh {
		t.Errorf("mapUrgencyToPriority(high) = %q, want %q", got, constants.PriorityHigh)
	}
}

func TestMapUrgencyToPriorityLow(t *testing.T) {
	got := mapUrgencyToPriority(constants.UrgencyLow)
	if got != constants.PriorityLow {
		t.Errorf("mapUrgencyToPriority(low) = %q, want %q", got, constants.PriorityLow)
	}
}

func TestMapUrgencyToPriorityMedium(t *testing.T) {
	got := mapUrgencyToPriority(constants.UrgencyMedium)
	if got != constants.PriorityMedium {
		t.Errorf("mapUrgencyToPriority(medium) = %q, want %q", got, constants.PriorityMedium)
	}
}

func TestMapUrgencyToPriorityUnknown(t *testing.T) {
	got := mapUrgencyToPriority("something_else")
	if got != constants.PriorityMedium {
		t.Errorf("mapUrgencyToPriority(unknown) = %q, want %q", got, constants.PriorityMedium)
	}
}

func TestMapUrgencyToPriorityEmpty(t *testing.T) {
	got := mapUrgencyToPriority("")
	if got != constants.PriorityMedium {
		t.Errorf("mapUrgencyToPriority('') = %q, want %q", got, constants.PriorityMedium)
	}
}
