CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS objects (
    object_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),  -- Unique object ID (UUID)
    object_name VARCHAR(25) NOT NULL,                      -- Object name (max 25 characters)
    object_description VARCHAR(256),                       -- Object description (max 256 characters)
    object_file BYTEA,                                     -- .glb file stored in binary format
     owner INT REFERENCES users(id_user),                  -- Object owner's ID referencing users table
    movable BOOLEAN NOT NULL DEFAULT FALSE,               -- Indicates if the object is movable (true or false)
    printable BOOLEAN NOT NULL DEFAULT FALSE              -- Indicates if the object is in the printing stand (true or false)
);
