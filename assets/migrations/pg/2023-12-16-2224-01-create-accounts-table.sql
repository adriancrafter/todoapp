--UP
CREATE TABLE accounts (
                          id UUID PRIMARY KEY,
                          slug VARCHAR(64) UNIQUE,
                          tenant_id UUID,
                          owner_id UUID,
                          parent_id UUID,
                          account_type VARCHAR(36),
                          username VARCHAR(32) UNIQUE,
                          email VARCHAR(255),
                          given_name VARCHAR(32),
                          middle_names VARCHAR(32) NULL,
                          family_names VARCHAR(64)
);

--UP
ALTER TABLE accounts
    ADD COLUMN locale VARCHAR(32),
    ADD COLUMN base_tz VARCHAR(64),
    ADD COLUMN current_tz VARCHAR(64),
    ADD COLUMN is_active BOOLEAN,
    ADD COLUMN created_by_id UUID references users(id),
    ADD COLUMN updated_by_id UUID references users(id),
    ADD COLUMN deleted_by_id UUID references users(id),
    ADD COLUMN created_at TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP,
    ADD COLUMN deleted_at TIMESTAMP;

--DOWN
DROP TABLE accounts;
