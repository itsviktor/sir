CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE basic_types (
    id                  BIGSERIAL PRIMARY KEY,

    field_smallint      SMALLINT      NOT NULL,
    field_integer       INTEGER       NOT NULL,
    field_bigint        BIGINT        NOT NULL,

    field_real          REAL          NOT NULL,
    field_double        DOUBLE PRECISION NOT NULL,
    field_numeric       NUMERIC(18,4) NOT NULL,

    field_boolean       BOOLEAN       NOT NULL,

    field_text          TEXT          NOT NULL,
    field_varchar       VARCHAR(255)  NOT NULL,
    field_char          CHAR(8)       NOT NULL,

    field_uuid          UUID          NOT NULL DEFAULT gen_random_uuid(),

    created_at          TIMESTAMP     NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE TABLE nullable_types (
    id                      BIGSERIAL PRIMARY KEY,

    field_smallint          SMALLINT,
    field_integer           INTEGER,
    field_bigint            BIGINT,

    field_real              REAL,
    field_double            DOUBLE PRECISION,
    field_numeric           NUMERIC(18,4),

    field_boolean           BOOLEAN,

    field_text              TEXT,
    field_varchar           VARCHAR(255),
    field_char              CHAR(8),

    field_uuid              UUID,

    field_timestamp         TIMESTAMP,
    field_timestamptz       TIMESTAMPTZ,

    field_date              DATE,
    field_time              TIME,
    field_timetz            TIME WITH TIME ZONE,

    field_interval          INTERVAL,

    field_json              JSON,
    field_jsonb             JSONB,

    field_bytea             BYTEA
);

CREATE TABLE array_types (
    id                  BIGSERIAL PRIMARY KEY,

    field_int_array     INTEGER[]      NOT NULL,
    field_text_array    TEXT[]         NOT NULL,
    field_uuid_array    UUID[]         NOT NULL,
    field_bool_array    BOOLEAN[]      NOT NULL
);

CREATE TABLE network_types (
    id                  BIGSERIAL PRIMARY KEY,

    field_inet          INET        NOT NULL,
    field_cidr          CIDR        NOT NULL,
    field_macaddr       MACADDR     NOT NULL
);

CREATE TYPE order_status AS ENUM (
    'pending',
    'paid',
    'cancelled'
);

CREATE TABLE enum_types (
    id                  BIGSERIAL PRIMARY KEY,

    field_status        order_status NOT NULL
);

CREATE TABLE relation_types (
    id                  BIGSERIAL PRIMARY KEY,

    basic_id            BIGINT NOT NULL
        REFERENCES basic_types(id)
        ON DELETE CASCADE,

    nullable_id         BIGINT
        REFERENCES nullable_types(id)
        ON DELETE SET NULL
);