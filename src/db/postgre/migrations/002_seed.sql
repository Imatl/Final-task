INSERT INTO supportflow.customers (id, name, email, plan) VALUES
    ('a0000000-0000-0000-0000-000000000001', 'Ivan Petrov', 'ivan@example.com', 'premium'),
    ('a0000000-0000-0000-0000-000000000002', 'Maria Sidorova', 'maria@example.com', 'basic'),
    ('a0000000-0000-0000-0000-000000000003', 'Alex Kozlov', 'alex@example.com', 'premium'),
    ('a0000000-0000-0000-0000-000000000004', 'Elena Novikova', 'elena@example.com', 'free'),
    ('a0000000-0000-0000-0000-000000000005', 'Dmitry Volkov', 'dmitry@example.com', 'basic')
ON CONFLICT (id) DO NOTHING;

INSERT INTO supportflow.agents (id, name, email, role, is_online) VALUES
    ('b0000000-0000-0000-0000-000000000001', 'Support Agent 1', 'agent1@supportflow.com', 'agent', true),
    ('b0000000-0000-0000-0000-000000000002', 'Support Agent 2', 'agent2@supportflow.com', 'agent', true),
    ('b0000000-0000-0000-0000-000000000003', 'Team Lead', 'lead@supportflow.com', 'lead', true)
ON CONFLICT (id) DO NOTHING;
