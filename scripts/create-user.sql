CREATE TABLE users (
    id_user SERIAL PRIMARY KEY,          -- Auto-increment primary key for identifying users
    username VARCHAR(50) NOT NULL,       -- Unique username for each user
    email VARCHAR(50) NOT NULL,          -- Unique username for each user
    password TEXT NOT NULL,              -- Hashed password for secure authentication
    created_at TIMESTAMP DEFAULT NOW()   -- Timestamp indicating when the user was created
);

CREATE TABLE roles (
    id_role SERIAL PRIMARY KEY,              -- Auto-increment primary key for identifying roles
    role_name VARCHAR(50) NOT NULL UNIQUE,   -- Unique name of the role (e.g., "admin", "moderator")
    description TEXT                         -- Optional description explaining the purpose or permissions of the role
);

CREATE TABLE account_roles (
    id_account_roles SERIAL PRIMARY KEY,            -- Auto-increment primary key for identifying the relationship
    id_user INT NOT NULL,                           -- Reference to the user ID in the 'users' table
    id_role INT NOT NULL,                           -- Reference to the role ID in the 'roles' table
    granted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),  -- Timestamp when the role was assigned
    revoked_at TIMESTAMPTZ,                         -- Timestamp when the role was revoked, NULL if active
    CONSTRAINT fk_account FOREIGN KEY (id_user) REFERENCES users (id_user),
    CONSTRAINT fk_role FOREIGN KEY (id_role) REFERENCES roles (id_role)
);

INSERT INTO roles (role_name, description) VALUES 
('user', 'Default role for standard users'),
('admin', 'Role for system administrators')
ON CONFLICT (role_name) DO NOTHING;