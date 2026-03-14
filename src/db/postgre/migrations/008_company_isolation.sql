ALTER TABLE supportflow.tickets ADD COLUMN IF NOT EXISTS company VARCHAR(255);
ALTER TABLE supportflow.customers ADD COLUMN IF NOT EXISTS company VARCHAR(255);
ALTER TABLE supportflow.agents ADD COLUMN IF NOT EXISTS company VARCHAR(255);

CREATE INDEX IF NOT EXISTS idx_tickets_company ON supportflow.tickets(company);
CREATE INDEX IF NOT EXISTS idx_customers_company ON supportflow.customers(company);
CREATE INDEX IF NOT EXISTS idx_agents_company ON supportflow.agents(company);
