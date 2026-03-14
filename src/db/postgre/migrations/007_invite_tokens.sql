CREATE TABLE IF NOT EXISTS supportflow.invite_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token       UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    created_by  UUID NOT NULL REFERENCES supportflow.users(id),
    used_by     UUID REFERENCES supportflow.users(id),
    used        BOOLEAN NOT NULL DEFAULT false,
    expires_at  TIMESTAMPTZ NOT NULL DEFAULT now() + interval '7 days',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_invite_tokens_token ON supportflow.invite_tokens(token);

ALTER TABLE supportflow.invite_tokens DROP COLUMN IF EXISTS company;
