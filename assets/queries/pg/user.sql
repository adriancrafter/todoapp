-- Model: User
-- Table: users

-- GetAll
SELECT * FROM users WHERE deleted_at IS NULL;

-- GetOne
SELECT * FROM users WHERE id = :id AND deleted_at IS NULL LIMIT 1;

-- Create
INSERT INTO users (id, tenant_id, slug, username, password, password_digest, email, email_confirmation, last_ip, confirmation_token, is_confirmed, last_geo_location, since, until, is_active, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :tenant_id, :slug, :username, :password, :password_digest, :email, :email_confirmation, :last_ip, :confirmation_token, :is_confirmed, :last_geo_location, :since, :until, :is_active, :created_by_id, :updated_by_id, :created_at, :updated_at);

-- Update
UPDATE users SET username = :username, password = :password, password_digest = :password_digest, email = :email, email_confirmation = :email_confirmation, last_ip = :last_ip, confirmation_token = :confirmation_token, is_confirmed = :is_confirmed, last_geo_location = :last_geo_location, since = :since, until = :until, is_active = :is_active, updated_by_id = :updated_by_id, updated_at = :updated_at
WHERE id = :id;

-- SoftDelete
UPDATE users SET deleted_by = :deleted_by, deleted_at = :deleted_at WHERE id = :id;

-- Delete
DELETE FROM users WHERE id = :id;

-- Purge
DELETE FROM users WHERE deleted_at IS NOT NULL;
