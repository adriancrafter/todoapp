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
