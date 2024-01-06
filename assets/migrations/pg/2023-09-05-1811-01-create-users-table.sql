--UP
CREATE TABLE users (
	  id UUID PRIMARY KEY,
	  slug VARCHAR(64) UNIQUE,
	  tenant_id UUID,
	  username VARCHAR(32) UNIQUE,
	  password_digest CHAR(128),
	  email VARCHAR(255) UNIQUE,
	  last_ip INET
	);

--UP
ALTER TABLE users
  ADD COLUMN confirmation_token VARCHAR(36),
  ADD COLUMN is_confirmed BOOLEAN,
  ADD COLUMN last_geolocation geography (Point,4326),
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
