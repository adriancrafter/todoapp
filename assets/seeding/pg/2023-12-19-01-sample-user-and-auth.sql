--UP
INSERT INTO users (id, slug, tenant_id, username, password_digest, email, last_ip, is_confirmed, is_active, created_at, updated_at)
VALUES ('e5b751c8-28b2-49d4-9769-4f5a3c6c7fd3', 'superadmin-slug', 'tenant-id-1', 'superadmin', '$2b$12$kglJGUQlQs2Wizfldpxlk./hP.dNB6GBYEGfBzmhbXzzpFyjv/EPS', 'superadmin@example.com', '192.168.1.1', TRUE, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

--UP
INSERT INTO accounts (id, slug, tenant_id, owner_id, username, email, given_name, is_active, created_at, updated_at)
VALUES ('12e7c7d8-c574-4532-a10a-a1e1d72e2d03', 'account-superadmin-slug', 'tenant-id-1', 'e5b751c8-28b2-49d4-9769-4f5a3c6c7fd3', 'superadmin-account', 'superadmin-account@example.com', 'Super Admin', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

--UP
INSERT INTO roles (id, slug, tenant_id, name, is_active, created_at, updated_at)
VALUES ('04a58d8f-1136-47d3-ba6d-d80d7db57500', 'role-1-slug', 'tenant-id-1', 'Admin Role', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

--UP
INSERT INTO account_roles (id, slug, tenant_id, account_id, role_id, is_active, created_at, updated_at)
VALUES ('be18283d-a084-4c23-bd60-04b50a405c84', 'accountrole-1-slug', 'tenant-id-1', '12e7c7d8-c574-4532-a10a-a1e1d72e2d03', '04a58d8f-1136-47d3-ba6d-d80d7db57500', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

--UP
INSERT INTO permissions (id, slug, tenant_id, name, tag, path, is_active, created_at, updated_at)
VALUES ('a98f3055-8e68-4475-9293-18dac651a21b', 'permission-1-slug', 'tenant-id-1', 'Super Admin Permission', 'superadmin_permission', '/admin/super', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

--UP
INSERT INTO role_permissions (id, slug, tenant_id, role_id, permission_id, is_active, created_at, updated_at)
VALUES ('90408582-aa19-4e86-977e-3e0b69989908', 'rolepermission-1-slug', 'tenant-id-1', '04a58d8f-1136-47d3-ba6d-d80d7db57500', 'a98f3055-8e68-4475-9293-18dac651a21b', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
