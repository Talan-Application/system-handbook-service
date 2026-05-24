CREATE TABLE locale_entity
(
    language_code VARCHAR(10) PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    native_name   VARCHAR(255) NOT NULL,
    is_active     BOOLEAN      NOT NULL DEFAULT true,
    is_deleted    BOOLEAN      NOT NULL DEFAULT false,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE translation
(
    id                BIGSERIAL PRIMARY KEY,
    translation_key   VARCHAR(255) NOT NULL,
    language_code     VARCHAR(10)  NOT NULL REFERENCES locale_entity (language_code),
    translation_value TEXT         NOT NULL,
    output_channel    VARCHAR(50),
    description       VARCHAR(255),
    is_deleted        BOOLEAN      NOT NULL DEFAULT false,
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    UNIQUE (translation_key, language_code, output_channel)
);

CREATE UNIQUE INDEX translation_key_locale_idx
    ON translation (translation_key, language_code);

CREATE SEQUENCE IF NOT EXISTS translation_key_seq
    INCREMENT BY 1
    MINVALUE 1
    START 1;