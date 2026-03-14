package constants

const (
	TicketStatusOpen       = "open"
	TicketStatusInProgress = "in_progress"
	TicketStatusWaiting    = "waiting"
	TicketStatusResolved   = "resolved"
	TicketStatusClosed     = "closed"

	PriorityLow      = "low"
	PriorityMedium   = "medium"
	PriorityHigh     = "high"
	PriorityCritical = "critical"

	RoleCustomer  = "customer"
	RoleAI        = "ai"
	RoleAgent     = "agent"
	RoleSystem    = "system"

	ActionStatusPending  = "pending"
	ActionStatusApproved = "approved"
	ActionStatusExecuted = "executed"
	ActionStatusRejected = "rejected"

	ActionRefund        = "refund"
	ActionChangePlan    = "change_plan"
	ActionResetPassword = "reset_password"
	ActionEscalate      = "escalate"
	ActionSendEmail     = "send_email"
	ActionCancelSub     = "cancel_subscription"

	SentimentPositive = "positive"
	SentimentNeutral  = "neutral"
	SentimentNegative = "negative"
	SentimentAngry    = "angry"

	UrgencyLow    = "low"
	UrgencyMedium = "medium"
	UrgencyHigh   = "high"

	AutoExecConfidenceThreshold = 0.85
)
