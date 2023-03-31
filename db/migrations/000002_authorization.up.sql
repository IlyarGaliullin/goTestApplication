CREATE TABLE IF NOT EXISTS users
(
    id       INTEGER GENERATED ALWAYS AS IDENTITY
        CONSTRAINT users_pkey
            PRIMARY KEY,
    name     VARCHAR(255) NOT NULL,
    email    VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS roles
(
    id   INTEGER GENERATED ALWAYS AS IDENTITY
        CONSTRAINT roles_pkey
            PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS userRoles
(
    id     INTEGER GENERATED ALWAYS AS IDENTITY
        CONSTRAINT userRoles_pkey
            PRIMARY KEY,
    userId INTEGER,
    roleId INTEGER,
    CONSTRAINT fk_user
        FOREIGN KEY (userId) REFERENCES users (id),
    CONSTRAINT fk_role
        FOREIGN KEY (roleId) REFERENCES roles (id)
);

CREATE TABLE IF NOT EXISTS grants
(
    id       INTEGER GENERATED ALWAYS AS IDENTITY
        CONSTRAINT pk_id
            PRIMARY KEY,
    roleId   INTEGER      NOT NULL,
    onTable  VARCHAR(255) NOT NULL,
    read     BOOLEAN DEFAULT FALSE,
    "create" BOOLEAN DEFAULT FALSE,
    "update" BOOLEAN DEFAULT FALSE,
    "delete" BOOLEAN DEFAULT FALSE,
    CONSTRAINT fk_role
        FOREIGN KEY (roleId) REFERENCES roles (id)
)