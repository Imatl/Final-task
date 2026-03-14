CREATE TABLE IF NOT EXISTS supportflow.users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       VARCHAR(255) UNIQUE NOT NULL,
    name        VARCHAR(255) NOT NULL,
    google_sub  VARCHAR(255),
    password    VARCHAR(255),
    level       INT NOT NULL DEFAULT 1,
    role        VARCHAR(50) NOT NULL DEFAULT 'Support',
    company     VARCHAR(255),
    created_at  TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON supportflow.users(email);
CREATE INDEX IF NOT EXISTS idx_users_google_sub ON supportflow.users(google_sub);

INSERT INTO supportflow.users (email, name, password, level, role)
VALUES ('admin@test.com', 'Admin User', 'password123', 5, 'Owner')
ON CONFLICT (email) DO NOTHING;
