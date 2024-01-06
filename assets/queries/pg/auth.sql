-- Model: Auth
-- Table: users

-- SignIn
SELECT users.*, array_to_string(array_agg(DISTINCT permissions.tag), ',') as permission_tags FROM users
    INNER JOIN accounts ON accounts.owner_id = users.id
	INNER JOIN account_roles ON account_roles.account_id = accounts.id
	INNER JOIN roles ON roles.id = account_roles.role_id
	INNER JOIN role_permissions ON role_permissions.role_id = roles.id
	INNER JOIN permissions ON permissions.id = role_permissions.permission_id
	WHERE users.deleted_at IS NULL
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
    AND (users.username = $1 OR users.email = $2)
    GROUP BY users.ID;

