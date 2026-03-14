CREATE TABLE IF NOT EXISTS supportflow.knowledge_base (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company     VARCHAR(255) NOT NULL,
    question    TEXT NOT NULL,
    answer      TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_kb_company ON supportflow.knowledge_base(company);
