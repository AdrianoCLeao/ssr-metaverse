CREATE TABLE IF NOT EXISTS object_position (
    object_id UUID PRIMARY KEY REFERENCES objects(object_id) ON DELETE CASCADE,
    x DOUBLE PRECISION NOT NULL,  -- X coordinate of the object
    y DOUBLE PRECISION NOT NULL,  -- Y coordinate of the object
    z DOUBLE PRECISION NOT NULL   -- Z coordinate of the object
);

CREATE TABLE IF NOT EXISTS object_rotation (
    object_id UUID PRIMARY KEY REFERENCES objects(object_id) ON DELETE CASCADE,
    rx DOUBLE PRECISION NOT NULL,  -- Rotation around X axis
    ry DOUBLE PRECISION NOT NULL,  -- Rotation around Y axis
    rz DOUBLE PRECISION NOT NULL   -- Rotation around Z axis
);

CREATE TABLE IF NOT EXISTS object_scale (
    object_id UUID PRIMARY KEY REFERENCES objects(object_id) ON DELETE CASCADE,
    scale_x DOUBLE PRECISION NOT NULL,  -- Scale along X axis
    scale_y DOUBLE PRECISION NOT NULL,  -- Scale along Y axis
    scale_z DOUBLE PRECISION NOT NULL   -- Scale along Z axis
);
