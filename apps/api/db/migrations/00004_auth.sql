-- Accounts (P0-ENG-06, ADR-0007).
--
-- users, sessions, accounts and verifications are Better Auth's tables. Their
-- SQL comes from `pnpm --filter @vernacular/web auth:schema`; copy changes from
-- its output into a new migration when Better Auth is upgraded.

-- +goose Up
CREATE TABLE users (
    id             text        NOT NULL PRIMARY KEY,
    name           text        NOT NULL,
    email          text        NOT NULL UNIQUE,
    email_verified boolean     NOT NULL,
    image          text,
    created_at     timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    role           text        NOT NULL DEFAULT 'developer' CHECK (role IN ('developer', 'contributor', 'admin')), -- contract Role
    banned         boolean,
    ban_reason     text,
    ban_expires    timestamptz
);

CREATE TABLE sessions (
    id              text        NOT NULL PRIMARY KEY,
    expires_at      timestamptz NOT NULL,
    token           text        NOT NULL UNIQUE,
    created_at      timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      timestamptz NOT NULL,
    ip_address      text,
    user_agent      text,
    user_id         text        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    impersonated_by text
);
CREATE INDEX sessions_user_id_idx ON sessions (user_id);

CREATE TABLE accounts (
    id                       text        NOT NULL PRIMARY KEY,
    account_id               text        NOT NULL,
    provider_id              text        NOT NULL,
    user_id                  text        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    access_token             text,
    refresh_token            text,
    id_token                 text,
    access_token_expires_at  timestamptz,
    refresh_token_expires_at timestamptz,
    scope                    text,
    password                 text,
    created_at               timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               timestamptz NOT NULL
);
CREATE INDEX accounts_user_id_idx ON accounts (user_id);

CREATE TABLE verifications (
    id         text        NOT NULL PRIMARY KEY,
    identifier text        NOT NULL,
    value      text        NOT NULL,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX verifications_identifier_idx ON verifications (identifier);

-- Only users is audited: sessions, accounts and verifications hold tokens and
-- password hashes that must not be copied into audit_log.
SELECT track('users');

-- Per-language review rights (CON-2, DQ-3). Not a global role, because a
-- reviewer is approved for specific languages only.
CREATE TYPE reviewer_level AS ENUM ('reviewer', 'senior');

CREATE TABLE reviewer_approvals (
    user_id       text           NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    language_code text           NOT NULL REFERENCES languages (code),
    level         reviewer_level NOT NULL,
    approved_by   text           REFERENCES users (id),
    created_at    timestamptz    NOT NULL DEFAULT now(),
    updated_at    timestamptz    NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, language_code)
);
SELECT track('reviewer_approvals', '{user_id,language_code}');

-- The foreign keys 00003 left for this migration.
ALTER TABLE items   ADD FOREIGN KEY (author_id)   REFERENCES users (id);
ALTER TABLE reviews ADD FOREIGN KEY (reviewer_id) REFERENCES users (id);

-- +goose Down
ALTER TABLE reviews DROP CONSTRAINT reviews_reviewer_id_fkey;
ALTER TABLE items   DROP CONSTRAINT items_author_id_fkey;
DROP TABLE reviewer_approvals;
DROP TYPE reviewer_level;
DROP TABLE verifications;
DROP TABLE accounts;
DROP TABLE sessions;
DROP TABLE users;
