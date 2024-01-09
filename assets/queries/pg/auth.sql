-- Model: Auth
-- Table: users

-- Signin
SELECT
    users.id,
    users.slug,
    users.tenant_id,
    users.username,
    users.password_digest,
    users.email,
    users.last_ip,
    users.confirmation_token,
    users.is_confirmed,
    public.ST_AsBinary(last_geolocation) AS last_geolocation,
    users.since,
    users.until,
    users.is_active,
    users.created_by_id,
    users.updated_by_id,
    users.deleted_by_id,
    users.created_at,
    users.updated_at,
    users.deleted_at,
    array_to_string(array_agg(DISTINCT permissions.tag), ',') AS permission_tags
    FROM users
        INNER JOIN accounts ON accounts.owner_id = users.id
        INNER JOIN account_roles ON account_roles.account_id = accounts.id
        INNER JOIN roles ON roles.id = account_roles.role_id
        INNER JOIN role_permissions ON role_permissions.role_id = roles.id
        INNER JOIN permissions ON permissions.id = role_permissions.permission_id
    WHERE (users.tenant_id::text = $1 OR users.tenant_id IS NULL OR users.tenant_id::text = '00000000-0000-0000-0000-000000000000')
        AND (accounts.deleted_at IS NULL)
        AND (account_roles.deleted_at IS NULL)
        AND (roles.deleted_at IS NULL)
        AND (role_permissions.deleted_at IS NULL)
        AND (permissions.deleted_at IS NULL)
        AND (users.is_active IS NULL OR users.is_active)
        AND (accounts.is_active IS NULL OR accounts.is_active)
        AND (account_roles.is_active IS NULL OR account_roles.is_active)
        AND (roles.is_active IS NULL OR roles.is_active)
        AND (role_permissions.is_active IS NULL OR role_permissions.is_active)
        AND (permissions.is_active IS NULL OR permissions.is_active)
        AND (users.username = $2 OR users.email = $3)
    GROUP BY users.ID;

-- UpdateSignin
UPDATE users
SET last_ip = NULLIF($1, '')::inet,
    last_geolocation = public.ST_SetSRID(public.ST_MakePoint($2::double precision, $3::double precision), 4326)
    WHERE (users.tenant_id::uuid = $4 OR users.tenant_id::uuid = '00000000-0000-0000-0000-000000000000')
        AND slug = $5;

