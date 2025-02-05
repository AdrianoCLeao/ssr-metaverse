# Application Data Model

This document describes the data model, including the tables and their relationships.

## Database Tables

### 1. `users` Table

Stores user information.

| Column      | Type         | Constraints                      | Description |
|------------|------------|---------------------------------|-------------|
| id_user    | SERIAL     | PRIMARY KEY                     | Unique identifier for the user. |
| username   | VARCHAR(50) | NOT NULL                        | Unique username. |
| email      | VARCHAR(50) | NOT NULL                        | User's email address. |
| password   | TEXT       | NOT NULL                        | Hashed password for secure authentication. |
| created_at | TIMESTAMP  | DEFAULT NOW()                   | Timestamp indicating when the user was created. |

---

### 2. `roles` Table

Stores different user roles.

| Column      | Type         | Constraints                      | Description |
|------------|------------|---------------------------------|-------------|
| id_role    | SERIAL     | PRIMARY KEY                     | Unique identifier for the role. |
| role_name  | VARCHAR(50) | NOT NULL, UNIQUE                | Unique role name (e.g., "admin", "user"). |
| description| TEXT       |                                 | Optional description of the role's purpose. |

---

### 3. `account_roles` Table

Defines the many-to-many relationship between users and roles.

| Column           | Type            | Constraints                               | Description |
|-----------------|---------------|-----------------------------------------|-------------|
| id_account_roles| SERIAL        | PRIMARY KEY                            | Unique identifier for the user-role relationship. |
| id_user         | INT           | NOT NULL, FOREIGN KEY → `users(id_user)` | References the user assigned to the role. |
| id_role         | INT           | NOT NULL, FOREIGN KEY → `roles(id_role)` | References the assigned role. |
| granted_at      | TIMESTAMPTZ    | NOT NULL, DEFAULT NOW()                  | Timestamp when the role was assigned. |
| revoked_at      | TIMESTAMPTZ    | NULLABLE                                 | Timestamp when the role was revoked (NULL if active). |

---

### 4. `objects` Table

Stores information about 3D objects in the application.

| Column              | Type         | Constraints                               | Description |
|---------------------|------------|-----------------------------------------|-------------|
| object_id          | UUID        | PRIMARY KEY, DEFAULT `uuid_generate_v4()` | Unique object identifier. |
| object_name        | VARCHAR(25) | NOT NULL                                | Object name (max 25 characters). |
| object_description | VARCHAR(256)|                                         | Object description (max 256 characters). |
| object_file        | BYTEA       |                                         | Binary storage for `.glb` files. |
| owner             | INT          | FOREIGN KEY → `users(id_user)`          | ID of the user who owns the object. |
| movable           | BOOLEAN      | NOT NULL, DEFAULT FALSE                 | Indicates if the object is movable. |
| printable         | BOOLEAN      | NOT NULL, DEFAULT FALSE                 | Indicates if the object is printable. |

---

### 5. `object_position` Table

Stores the position of each object in a 3D space.

| Column    | Type              | Constraints                         | Description |
|----------|------------------|---------------------------------|-------------|
| object_id | UUID            | PRIMARY KEY, FOREIGN KEY → `objects(object_id)`, ON DELETE CASCADE | Unique object reference. |
| x         | DOUBLE PRECISION | NOT NULL                         | X-coordinate of the object. |
| y         | DOUBLE PRECISION | NOT NULL                         | Y-coordinate of the object. |
| z         | DOUBLE PRECISION | NOT NULL                         | Z-coordinate of the object. |

---

### 6. `object_rotation` Table

Stores the rotation of each object in 3D space.

| Column    | Type              | Constraints                         | Description |
|----------|------------------|---------------------------------|-------------|
| object_id | UUID            | PRIMARY KEY, FOREIGN KEY → `objects(object_id)`, ON DELETE CASCADE | Unique object reference. |
| rx        | DOUBLE PRECISION | NOT NULL                         | Rotation around the X-axis. |
| ry        | DOUBLE PRECISION | NOT NULL                         | Rotation around the Y-axis. |
| rz        | DOUBLE PRECISION | NOT NULL                         | Rotation around the Z-axis. |

---

### 7. `object_scale` Table

Stores the scale of each object.

| Column    | Type              | Constraints                         | Description |
|----------|------------------|---------------------------------|-------------|
| object_id | UUID            | PRIMARY KEY, FOREIGN KEY → `objects(object_id)`, ON DELETE CASCADE | Unique object reference. |
| scale_x   | DOUBLE PRECISION | NOT NULL                         | Scale factor along the X-axis. |
| scale_y   | DOUBLE PRECISION | NOT NULL                         | Scale factor along the Y-axis. |
| scale_z   | DOUBLE PRECISION | NOT NULL                         | Scale factor along the Z-axis. |

---

