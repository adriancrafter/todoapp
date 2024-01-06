--UP
CREATE EXTENSION IF NOT EXISTS postgis;

--DOWN
DROP EXTENSION IF EXISTS postgis;
--UP
CREATE TABLE users (
	  id UUID PRIMARY KEY,
	  slug VARCHAR(64) UNIQUE,
	  tenant_id VARCHAR(128),
	  username VARCHAR(32) UNIQUE,
	  password_digest CHAR(128),
	  email VARCHAR(255) UNIQUE,
	  last_ip INET
	);

--UP
ALTER TABLE users
  ADD COLUMN confirmation_token VARCHAR(36),
  ADD COLUMN is_confirmed BOOLEAN,
  ADD COLUMN geolocation geography (Point,4326),
  ADD COLUMN since TIMESTAMP,
  ADD COLUMN until TIMESTAMP,
  ADD COLUMN is_active BOOLEAN,
  ADD COLUMN created_by_id UUID,
  ADD COLUMN updated_by_id UUID,
  ADD COLUMN deleted_by_id UUID,
  ADD COLUMN created_at TIMESTAMP,
  ADD COLUMN updated_at TIMESTAMP,
  ADD COLUMN deleted_at TIMESTAMP;

--DOWN
DROP TABLE users;
--UP
CREATE TABLE accounts (
                          id UUID PRIMARY KEY,
                          slug VARCHAR(64) UNIQUE,
                          tenant_id VARCHAR(128),
                          owner_id UUID,
                          parent_id UUID,
                          account_type VARCHAR(36),
                          username VARCHAR(32) UNIQUE,
                          email VARCHAR(255),
                          given_names VARCHAR(32),
                          middle_names VARCHAR(32) NULL,
                          family_names VARCHAR(64)
);

--UP
ALTER TABLE accounts
    ADD COLUMN locale VARCHAR(32),
    ADD COLUMN base_tz VARCHAR(64),
    ADD COLUMN current_tz VARCHAR(64),
    ADD COLUMN is_active BOOLEAN,
    ADD COLUMN is_deleted BOOLEAN,
    ADD COLUMN created_by_id UUID,
    ADD COLUMN updated_by_id UUID,
    ADD COLUMN deleted_by_id UUID,
    ADD COLUMN created_at TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP,
    ADD COLUMN deleted_at TIMESTAMP;

--DOWN
DROP TABLE accounts;
--UP
CREATE TABLE profiles (
                          id UUID PRIMARY KEY,
                          tenant_id VARCHAR(128),
                          slug VARCHAR(64) UNIQUE,
                          owner_id UUID REFERENCES users(id),
                          account_type VARCHAR(36),
                          name VARCHAR(64),
                          email VARCHAR(255) UNIQUE,
                          location VARCHAR(255) NULL,
                          bio VARCHAR(255),
                          moto VARCHAR(255),
                          website VARCHAR(255),
                          aniversary_date TIMESTAMP,
                          avatar_small bytea,
                          header_small VARCHAR(255),
                          avatar_path VARCHAR(255),
                          header_path VARCHAR(255),
                          extended_data jsonb
);

--UP
ALTER TABLE profiles
    ADD COLUMN geolocation geography (Point,4326),
    ADD COLUMN is_active BOOLEAN,
    ADD COLUMN is_deleted BOOLEAN,
    ADD COLUMN created_by_id UUID REFERENCES users(id),
    ADD COLUMN updated_by_id UUID REFERENCES users(id),
    ADD COLUMN created_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE;

--DOWN
DROP TABLE profiles;
--UP
CREATE TABLE resources (
                           id UUID PRIMARY KEY,
                           slug VARCHAR(64) UNIQUE,
                           tenant_id VARCHAR(128),
                           name VARCHAR(32) UNIQUE,
                           description TEXT NULL,
                           tag VARCHAR(16) UNIQUE,
                           path VARCHAR(512) UNIQUE
);

--UP
ALTER TABLE resources
    ADD COLUMN is_active BOOLEAN,
    ADD COLUMN is_deleted BOOLEAN,
    ADD COLUMN created_by_id UUID,
    ADD COLUMN updated_by_id UUID,
    ADD COLUMN created_at TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP;

--DOWN
DROP TABLE resources;
--UP
CREATE TABLE roles (
                       id UUID PRIMARY KEY,
                       slug VARCHAR(64) UNIQUE,
                       tenant_id VARCHAR(128),
                       name VARCHAR(32) UNIQUE,
                       description TEXT NULL
);

--UP
ALTER TABLE roles
    ADD COLUMN is_active BOOLEAN,
    ADD COLUMN is_deleted BOOLEAN,
    ADD COLUMN created_by_id UUID,
    ADD COLUMN updated_by_id UUID,
    ADD COLUMN created_at TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP;

--DOWN
DROP TABLE roles;
--UP
CREATE TABLE permissions (
                             id UUID PRIMARY KEY,
                             slug VARCHAR(64) UNIQUE,
                             tenant_id VARCHAR(128),
                             name VARCHAR(32) UNIQUE,
                             description TEXT NULL,
                             tag VARCHAR(16) UNIQUE,
                             path VARCHAR(512) UNIQUE
);

--UP
ALTER TABLE permissions
    ADD COLUMN is_active BOOLEAN,
    ADD COLUMN is_deleted BOOLEAN,
    ADD COLUMN created_by_id UUID,
    ADD COLUMN updated_by_id UUID,
    ADD COLUMN created_at TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP;

--DOWN
DROP TABLE permissions;
--UP
CREATE TABLE account_roles (
                               id UUID PRIMARY KEY,
                               slug VARCHAR(64) UNIQUE,
                               tenant_id VARCHAR(128),
                               name VARCHAR(32) UNIQUE,
                               account_id UUID NOT NULL,
                               role_id UUID NOT NULL
);

--UP
ALTER TABLE account_roles
    ADD COLUMN is_active BOOLEAN,
    ADD COLUMN is_deleted BOOLEAN,
    ADD COLUMN created_by_id UUID,
    ADD COLUMN updated_by_id UUID,
    ADD COLUMN created_at TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP,
    ADD UNIQUE (tenant_id, account_id, role_id);

--DOWN
DROP TABLE account_roles;
--UP
CREATE TABLE role_permissions (
                                  id UUID PRIMARY KEY,
                                  slug VARCHAR(64) UNIQUE,
                                  tenant_id VARCHAR(128),
                                  name VARCHAR(32) UNIQUE,
                                  role_id UUID NOT NULL,
                                  permission_id UUID NOT NULL
);

--UP
ALTER TABLE role_permissions
    ADD COLUMN is_active BOOLEAN,
    ADD COLUMN is_deleted BOOLEAN,
    ADD COLUMN created_by_id UUID,
    ADD COLUMN updated_by_id UUID,
    ADD COLUMN created_at TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP,
    ADD UNIQUE (tenant_id, role_id, permission_id);

--DOWN
DROP TABLE role_permissions;
--UP
CREATE TABLE resource_permissions (
                                      id UUID PRIMARY KEY,
                                      slug VARCHAR(64) UNIQUE,
                                      tenant_id VARCHAR(128),
                                      name VARCHAR(32) UNIQUE,
                                      resource_id UUID NOT NULL,
                                      permission_id UUID NOT NULL
);

--UP
ALTER TABLE resource_permissions
    ADD COLUMN is_active BOOLEAN,
    ADD COLUMN is_deleted BOOLEAN,
    ADD COLUMN created_by_id UUID,
    ADD COLUMN updated_by_id UUID,
    ADD COLUMN created_at TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP,
    ADD UNIQUE (tenant_id, resource_id, permission_id);

--DOWN
DROP TABLE resource_permissions;
