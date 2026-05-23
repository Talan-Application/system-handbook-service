CREATE TABLE common_subjects
(
    id         BIGSERIAL    PRIMARY KEY,
    name_key   VARCHAR(255) NOT NULL,
    is_deleted BOOLEAN      NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
