-- WARNING: This rollback will assign new integer IDs to all records
-- The new IDs will be different from the original ones
-- Only use this for development/testing purposes

-- Step 1: Add temporary integer ID columns
ALTER TABLE roles ADD COLUMN int_id SERIAL;
ALTER TABLE users ADD COLUMN int_id SERIAL;
ALTER TABLE user_details ADD COLUMN int_id SERIAL;
ALTER TABLE user_details ADD COLUMN int_user_id INTEGER;
ALTER TABLE user_roles ADD COLUMN int_id SERIAL;
ALTER TABLE user_roles ADD COLUMN int_user_id INTEGER;
ALTER TABLE user_roles ADD COLUMN int_role_id INTEGER;
ALTER TABLE user_roles ADD COLUMN int_assigned_by INTEGER;

-- Step 2: Map UUID foreign keys back to integers
UPDATE user_details SET int_user_id = u.int_id FROM users u WHERE user_details.user_id = u.id;
UPDATE user_roles SET int_user_id = u.int_id FROM users u WHERE user_roles.user_id = u.id;
UPDATE user_roles SET int_role_id = r.int_id FROM roles r WHERE user_roles.role_id = r.id;
UPDATE user_roles SET int_assigned_by = u.int_id FROM users u WHERE user_roles.assigned_by = u.id;

-- Step 3: Drop UUID foreign key constraints
ALTER TABLE user_details DROP CONSTRAINT user_details_user_id_fkey;
ALTER TABLE user_roles DROP CONSTRAINT user_roles_user_id_fkey;
ALTER TABLE user_roles DROP CONSTRAINT user_roles_role_id_fkey;
ALTER TABLE user_roles DROP CONSTRAINT user_roles_assigned_by_fkey;

-- Step 4: Drop UUID unique constraints
ALTER TABLE user_details DROP CONSTRAINT user_details_user_id_unique;
ALTER TABLE user_roles DROP CONSTRAINT user_roles_user_id_role_id_unique;

-- Step 5: Drop indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_is_active;
DROP INDEX IF EXISTS idx_user_details_user_id;
DROP INDEX IF EXISTS idx_user_roles_user_id;
DROP INDEX IF EXISTS idx_user_roles_role_id;

-- Step 6: Drop UUID PRIMARY KEY constraints
ALTER TABLE roles DROP CONSTRAINT roles_pkey;
ALTER TABLE users DROP CONSTRAINT users_pkey;
ALTER TABLE user_details DROP CONSTRAINT user_details_pkey;
ALTER TABLE user_roles DROP CONSTRAINT user_roles_pkey;

-- Step 7: Drop UUID columns
ALTER TABLE roles DROP COLUMN id;
ALTER TABLE users DROP COLUMN id;
ALTER TABLE user_details DROP COLUMN id, DROP COLUMN user_id;
ALTER TABLE user_roles DROP COLUMN id, DROP COLUMN user_id, DROP COLUMN role_id, DROP COLUMN assigned_by;

-- Step 8: Rename integer columns to standard names
ALTER TABLE roles RENAME COLUMN int_id TO id;
ALTER TABLE users RENAME COLUMN int_id TO id;
ALTER TABLE user_details RENAME COLUMN int_id TO id;
ALTER TABLE user_details RENAME COLUMN int_user_id TO user_id;
ALTER TABLE user_roles RENAME COLUMN int_id TO id;
ALTER TABLE user_roles RENAME COLUMN int_user_id TO user_id;
ALTER TABLE user_roles RENAME COLUMN int_role_id TO role_id;
ALTER TABLE user_roles RENAME COLUMN int_assigned_by TO assigned_by;

-- Step 9: Add PRIMARY KEY constraints with INTEGER
ALTER TABLE roles ADD PRIMARY KEY (id);
ALTER TABLE users ADD PRIMARY KEY (id);
ALTER TABLE user_details ADD PRIMARY KEY (id);
ALTER TABLE user_roles ADD PRIMARY KEY (id);

-- Step 10: Add FOREIGN KEY constraints with INTEGER
ALTER TABLE user_details ADD CONSTRAINT user_details_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE user_roles ADD CONSTRAINT user_roles_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE user_roles ADD CONSTRAINT user_roles_role_id_fkey
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;
ALTER TABLE user_roles ADD CONSTRAINT user_roles_assigned_by_fkey
    FOREIGN KEY (assigned_by) REFERENCES users(id);

-- Step 11: Add UNIQUE constraints back
ALTER TABLE user_details ADD CONSTRAINT user_details_user_id_unique UNIQUE (user_id);
ALTER TABLE user_roles ADD CONSTRAINT user_roles_user_id_role_id_unique UNIQUE (user_id, role_id);

-- Step 12: Recreate indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_user_details_user_id ON user_details(user_id);
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);

-- Step 13: Set NOT NULL constraints
ALTER TABLE roles ALTER COLUMN id SET NOT NULL;
ALTER TABLE users ALTER COLUMN id SET NOT NULL;
ALTER TABLE user_details ALTER COLUMN id SET NOT NULL;
ALTER TABLE user_details ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE user_roles ALTER COLUMN id SET NOT NULL;
ALTER TABLE user_roles ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE user_roles ALTER COLUMN role_id SET NOT NULL;
