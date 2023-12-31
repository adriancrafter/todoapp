-- Model: Account
-- Table: accounts

-- GetAll
SELECT * FROM accounts WHERE tenant_id = :tenant_id AND deleted_at IS NULL;

-- GetOne
SELECT * FROM accounts WHERE tenant_id = :tenant_id AND slug = :slug AND deleted_at IS NULL LIMIT 1;

-- Create
INSERT INTO accounts (tenant_id, user_id, slug, username, description, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:tenant_id, :user_id, :slug, :username, :description, :created_by_id, :updated_by_id, :created_at, :updated_at);

-- Update
UPDATE accounts SET username = :username, description = :description, updated_by_id = :updated_by_id, updated_at = :updated_at
WHERE tenant_id = :tenant_id AND slug = :slug;

-- SoftDelete
UPDATE accounts SET deleted_by = :deleted_by, deleted_at = :deleted_at WHERE tenant_id = :tenant_id AND slug = :slug;

-- Delete
DELETE FROM accounts WHERE tenant_id = :tenant_id AND slug = :slug;

-- Purge
DELETE FROM accounts WHERE tenant_id = :tenant_id AND deleted_at IS NOT NULL;

