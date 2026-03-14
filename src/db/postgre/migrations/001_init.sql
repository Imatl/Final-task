CREATE SCHEMA IF NOT EXISTS supportflow;

CREATE TABLE IF NOT EXISTS supportflow.customers (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    email      VARCHAR(255) UNIQUE NOT NULL,
    plan       VARCHAR(50) NOT NULL DEFAULT 'free',
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS supportflow.agents (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    email      VARCHAR(255) UNIQUE NOT NULL,
    role       VARCHAR(50) NOT NULL DEFAULT 'agent',
    is_online  BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS supportflow.tickets (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES supportflow.customers(id),
    agent_id    UUID REFERENCES supportflow.agents(id),
    subject     VARCHAR(500) NOT NULL,
    status      VARCHAR(50) NOT NULL DEFAULT 'open',
    priority    VARCHAR(20) NOT NULL DEFAULT 'medium',
    category    VARCHAR(100) NOT NULL DEFAULT 'general',
    ai_summary  TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now(),
    closed_at   TIMESTAMP
);

CREATE TABLE IF NOT EXISTS supportflow.messages (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id  UUID NOT NULL REFERENCES supportflow.tickets(id),
    role       VARCHAR(20) NOT NULL,
    content    TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS supportflow.actions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id   UUID NOT NULL REFERENCES supportflow.tickets(id),
    type        VARCHAR(100) NOT NULL,
    params      JSONB NOT NULL DEFAULT '{}',
    status      VARCHAR(50) NOT NULL DEFAULT 'pending',
    result      TEXT,
    confidence  DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    executed_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS supportflow.ai_analyses (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id       UUID NOT NULL REFERENCES supportflow.tickets(id),
    intent          VARCHAR(100) NOT NULL,
    sentiment       VARCHAR(50) NOT NULL,
    urgency         VARCHAR(20) NOT NULL,
    suggested_tools TEXT[],
    reasoning       TEXT,
    confidence      DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at      TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_tickets_status ON supportflow.tickets(status);
CREATE INDEX IF NOT EXISTS idx_tickets_customer ON supportflow.tickets(customer_id);
CREATE INDEX IF NOT EXISTS idx_tickets_agent ON supportflow.tickets(agent_id);
CREATE INDEX IF NOT EXISTS idx_messages_ticket ON supportflow.messages(ticket_id);
CREATE INDEX IF NOT EXISTS idx_actions_ticket ON supportflow.actions(ticket_id);
CREATE INDEX IF NOT EXISTS idx_analyses_ticket ON supportflow.ai_analyses(ticket_id);
