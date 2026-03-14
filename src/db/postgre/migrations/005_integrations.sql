CREATE TABLE IF NOT EXISTS supportflow.integrations (
    id         VARCHAR(100) PRIMARY KEY,
    type       VARCHAR(50)  NOT NULL,
    name       VARCHAR(200) NOT NULL DEFAULT '',
    config     JSONB        NOT NULL DEFAULT '{}',
    status     VARCHAR(50)  NOT NULL DEFAULT 'disconnected',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS supportflow.channel_mappings (
    customer_id UUID         NOT NULL REFERENCES supportflow.customers(id),
    channel     VARCHAR(50)  NOT NULL,
    external_id VARCHAR(200) NOT NULL,
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    PRIMARY KEY (customer_id, channel)
);
