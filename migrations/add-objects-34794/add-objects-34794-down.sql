-- Removing the 'email' column from the 'users' table
ALTER TABLE users DROP COLUMN IF EXISTS email;

-- Reverting 'id_user' and 'id_role' in 'account_roles' back to SERIAL
ALTER TABLE account_roles ALTER COLUMN id_user TYPE SERIAL USING id_user::SERIAL;
ALTER TABLE account_roles ALTER COLUMN id_role TYPE SERIAL USING id_role::SERIAL;

-- Dropping tables added in the migration UP
DROP TABLE IF EXISTS object_scale;
DROP TABLE IF EXISTS object_rotation;
DROP TABLE IF EXISTS object_position;
DROP TABLE IF EXISTS objects;
