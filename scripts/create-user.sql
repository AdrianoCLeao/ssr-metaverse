/*
Table: users
This table stores user account information.
- id_user: Auto-increment primary key for identifying users.
- username: Unique username for each user.
- password: Hashed password for secure authentication.
- created_at: Timestamp indicating when the user was created, with a default value of the current time.
*/
CREATE TABLE users (
    id_user SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

/*
Table: roles
This table defines the roles that can be assigned to users.
- id_role: Auto-increment primary key for identifying roles.
- role_name: Unique name of the role (e.g., "admin", "moderator").
- description: Optional description explaining the purpose or permissions of the role.
*/
CREATE TABLE roles (
    id_role SERIAL PRIMARY KEY,
    role_name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT
);

/*
Table: account_roles
This table manages the relationships between users and their assigned roles.
- id_account_roles: Auto-increment primary key for identifying the relationship.
- id_user: Reference to the user ID in the 'users' table.
- id_role: Reference to the role ID in the 'roles' table.
- granted_at: Timestamp of when the role was assigned to the user, with a default value of the current time.
- revoked_at: Timestamp of when the role was revoked from the user, NULL if the role is still active.
- fk_account: Foreign key constraint linking id_user to the 'users' table.
- fk_role: Foreign key constraint linking id_role to the 'roles' table.
*/
CREATE TABLE account_roles (
    id_account_roles SERIAL PRIMARY KEY,
    id_user UUID NOT NULL,
    id_role INT NOT NULL,
    granted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ,
    CONSTRAINT fk_account FOREIGN KEY (id_user) REFERENCES users (id_user),
    CONSTRAINT fk_role FOREIGN KEY (id_role) REFERENCES roles (id_role)
);
