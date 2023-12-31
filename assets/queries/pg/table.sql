--Model:Table
--GetAll
SELECT * FROM %s.tables WHERE tenant_id = :tenant_id
                           AND branch_id = :branch_id
                           AND deleted_at IS NULL;
--GetOne
SELECT * FROM %s.tables WHERE tenant_id = :tenant_id
                               AND branch_slug = :branch_slug
                               AND slug = :slug
                               AND deleted_at IS NULL
                               LIMIT 1;
--Create
INSERT INTO %s.tables (id, slug, tenant_id, branch_id, name, description, status, created_by_id, updated_by_id, deleted_by, created_at, updated_at, deleted_at)
VALUES (:id, :slug, :tenant_id, :branch_id, :name, :description, :status, :created_by_id, :updated_by_id, :deleted_by, :created_at, :updated_at, :deleted_at);

--Update
UPDATE %s.tables SET name = :name, description = :description, status = :status, updated_by_id = :updated_by_id, updated_at = :updated_at WHERE tenant_id = :tenant_id
                                AND branch_id = :branch_id
                                AND slug = :slug;

--SoftDelete
UPDATE %s.tables SET deleted_by = :deleted_by, deleted_at = :deleted_at
                 WHERE tenant_id = :tenant_id
                   AND branch_id = :branch_id
                   AND slug = :slug;

--Delete
DELETE FROM %s.tables WHERE tenant_id = :tenant_id
                     AND branch_id = :branch_id
                     AND slug = :slug;

--Purge
DELETE FROM %s.tables WHERE tenant_id = :tenant_id AND deleted_at IS NOT NULL;
