-- Adding the 'email' column to the 'users' table
ALTER TABLE users ADD COLUMN email VARCHAR(50) NOT NULL;

-- Changing the type of 'id_user' and 'id_role' in 'account_roles' from SERIAL to INT
ALTER TABLE account_roles ALTER COLUMN id_user TYPE INT USING id_user::INT;
ALTER TABLE account_roles ALTER COLUMN id_role TYPE INT USING id_role::INT;

-- Creating the UUID extension if it does not exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Creating the 'objects' table
CREATE TABLE IF NOT EXISTS objects (
    object_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- Unique object identifier
    object_name VARCHAR(25) NOT NULL,                      -- Object name (max 25 characters)
    object_description VARCHAR(256),                       -- Object description (max 256 characters)
    object_file BYTEA,                                     -- .glb file stored in binary format
    owner INT REFERENCES users(id_user),                   -- Reference to the object owner in 'users' table
    movable BOOLEAN NOT NULL DEFAULT FALSE,                -- Indicates whether the object is movable
    printable BOOLEAN NOT NULL DEFAULT FALSE               -- Indicates whether the object is in the printing stand
);

-- Creating the 'object_position' table to store object coordinates
CREATE TABLE IF NOT EXISTS object_position (
    object_id UUID PRIMARY KEY REFERENCES objects(object_id) ON DELETE CASCADE,
    x DOUBLE PRECISION NOT NULL,  -- X coordinate of the object
    y DOUBLE PRECISION NOT NULL,  -- Y coordinate of the object
    z DOUBLE PRECISION NOT NULL   -- Z coordinate of the object
);

-- Creating the 'object_rotation' table to store object rotation values
CREATE TABLE IF NOT EXISTS object_rotation (
    object_id UUID PRIMARY KEY REFERENCES objects(object_id) ON DELETE CASCADE,
    rx DOUBLE PRECISION NOT NULL,  -- Rotation around X axis
    ry DOUBLE PRECISION NOT NULL,  -- Rotation around Y axis
    rz DOUBLE PRECISION NOT NULL   -- Rotation around Z axis
);

-- Creating the 'object_scale' table to store object scaling values
CREATE TABLE IF NOT EXISTS object_scale (
    object_id UUID PRIMARY KEY REFERENCES objects(object_id) ON DELETE CASCADE,
    scale_x DOUBLE PRECISION NOT NULL,  -- Scale along X axis
    scale_y DOUBLE PRECISION NOT NULL,  -- Scale along Y axis
    scale_z DOUBLE PRECISION NOT NULL   -- Scale along Z axis
);
