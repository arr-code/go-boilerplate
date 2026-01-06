-- Enable UUID generation (PostgreSQL 13+ has gen_random_uuid built-in)
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Step 1: Add UUID columns to all tables
ALTER TABLE roles ADD COLUMN uuid_id UUID;
ALTER TABLE users ADD COLUMN uuid_id UUID;
ALTER TABLE user_details ADD COLUMN uuid_id UUID;
ALTER TABLE user_details ADD COLUMN uuid_user_id UUID;
ALTER TABLE user_roles ADD COLUMN uuid_id UUID;
ALTER TABLE user_roles ADD COLUMN uuid_user_id UUID;
ALTER TABLE user_roles ADD COLUMN uuid_role_id UUID;
ALTER TABLE user_roles ADD COLUMN uuid_assigned_by UUID;

-- Step 2: Generate UUIDs for existing rows (roles with fixed UUIDs)
UPDATE roles SET uuid_id = '550e8400-e29b-41d4-a716-446655440001'::uuid WHERE id = 1; -- user
UPDATE roles SET uuid_id = '550e8400-e29b-41d4-a716-446655440002'::uuid WHERE id = 2; -- admin
UPDATE roles SET uuid_id = '550e8400-e29b-41d4-a716-446655440003'::uuid WHERE id = 3; -- moderator

-- Generate random UUIDs for other rows
UPDATE users SET uuid_id = gen_random_uuid() WHERE uuid_id IS NULL;
UPDATE user_details SET uuid_id = gen_random_uuid() WHERE uuid_id IS NULL;
UPDATE user_roles SET uuid_id = gen_random_uuid() WHERE uuid_id IS NULL;

-- Step 3: Populate UUID foreign keys by mapping from integer IDs
UPDATE user_details SET uuid_user_id = u.uuid_id FROM users u WHERE user_details.user_id = u.id;
UPDATE user_roles SET uuid_user_id = u.uuid_id FROM users u WHERE user_roles.user_id = u.id;
UPDATE user_roles SET uuid_role_id = r.uuid_id FROM roles r WHERE user_roles.role_id = r.id;
UPDATE user_roles SET uuid_assigned_by = u.uuid_id FROM users u WHERE user_roles.assigned_by = u.id;

-- Step 4: Drop old foreign key constraints
ALTER TABLE user_details DROP CONSTRAINT user_details_user_id_fkey;
ALTER TABLE user_roles DROP CONSTRAINT user_roles_user_id_fkey;
ALTER TABLE user_roles DROP CONSTRAINT user_roles_role_id_fkey;
ALTER TABLE user_roles DROP CONSTRAINT user_roles_assigned_by_fkey;

-- Step 5: Drop old indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_is_active;
DROP INDEX IF EXISTS idx_user_details_user_id;
DROP INDEX IF EXISTS idx_user_roles_user_id;
DROP INDEX IF EXISTS idx_user_roles_role_id;

-- Step 6: Drop old PRIMARY KEY constraints
ALTER TABLE roles DROP CONSTRAINT roles_pkey;
ALTER TABLE users DROP CONSTRAINT users_pkey;
ALTER TABLE user_details DROP CONSTRAINT user_details_pkey;
ALTER TABLE user_roles DROP CONSTRAINT user_roles_pkey;

-- Step 7: Drop old integer columns
ALTER TABLE roles DROP COLUMN id;
ALTER TABLE users DROP COLUMN id;
ALTER TABLE user_details DROP COLUMN id, DROP COLUMN user_id;
ALTER TABLE user_roles DROP COLUMN id, DROP COLUMN user_id, DROP COLUMN role_id, DROP COLUMN assigned_by;

-- Step 8: Rename UUID columns to 'id', 'user_id', etc.
ALTER TABLE roles RENAME COLUMN uuid_id TO id;
ALTER TABLE users RENAME COLUMN uuid_id TO id;
ALTER TABLE user_details RENAME COLUMN uuid_id TO id;
ALTER TABLE user_details RENAME COLUMN uuid_user_id TO user_id;
ALTER TABLE user_roles RENAME COLUMN uuid_id TO id;
ALTER TABLE user_roles RENAME COLUMN uuid_user_id TO user_id;
ALTER TABLE user_roles RENAME COLUMN uuid_role_id TO role_id;
ALTER TABLE user_roles RENAME COLUMN uuid_assigned_by TO assigned_by;

-- Step 9: Add PRIMARY KEY constraints with UUID
ALTER TABLE roles ADD PRIMARY KEY (id);
ALTER TABLE users ADD PRIMARY KEY (id);
ALTER TABLE user_details ADD PRIMARY KEY (id);
ALTER TABLE user_roles ADD PRIMARY KEY (id);

-- Step 10: Add FOREIGN KEY constraints with UUID
ALTER TABLE user_details ADD CONSTRAINT user_details_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE user_roles ADD CONSTRAINT user_roles_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE user_roles ADD CONSTRAINT user_roles_role_id_fkey
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;
ALTER TABLE user_roles ADD CONSTRAINT user_roles_assigned_by_fkey
    FOREIGN KEY (assigned_by) REFERENCES users(id);

-- Step 11: Add UNIQUE constraint back
ALTER TABLE user_details ADD CONSTRAINT user_details_user_id_unique UNIQUE (user_id);
ALTER TABLE user_roles ADD CONSTRAINT user_roles_user_id_role_id_unique UNIQUE (user_id, role_id);

-- Step 12: Recreate indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_user_details_user_id ON user_details(user_id);
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);

-- Step 13: Set NOT NULL constraints on new UUID columns
ALTER TABLE roles ALTER COLUMN id SET NOT NULL;
ALTER TABLE users ALTER COLUMN id SET NOT NULL;
ALTER TABLE user_details ALTER COLUMN id SET NOT NULL;
ALTER TABLE user_details ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE user_roles ALTER COLUMN id SET NOT NULL;
ALTER TABLE user_roles ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE user_roles ALTER COLUMN role_id SET NOT NULL;

-- Step 14: Set default UUID generation for new rows
ALTER TABLE roles ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE users ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE user_details ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE user_roles ALTER COLUMN id SET DEFAULT gen_random_uuid();
