-- Тестовая схема для интроспекции и вывода типов (SQLite).
-- Скрипт идемпотентный: можно прогонять повторно.
--
-- Соглашение по именам: колонка без суффикса — NOT NULL,
-- колонка с суффиксом _null — nullable-версия того же типа.
-- В комментариях: affinity по правилам SQLite -> ожидаемый Go-тип.

DROP TABLE IF EXISTS all_types;
DROP TABLE IF EXISTS all_types_strict;
DROP TABLE IF EXISTS quirks;

-- 1. Обычная (не STRICT) таблица: типичные декларированные типы.
CREATE TABLE all_types (
    id              INTEGER PRIMARY KEY,                              -- алиас rowid, всегда NOT NULL -> int64

    -- INTEGER affinity
    int_val         INTEGER NOT NULL DEFAULT 0,                       -- int64
    int_null        INTEGER,                                          -- *int64
    bigint_val      BIGINT NOT NULL DEFAULT 0,                        -- int64
    bigint_null     BIGINT,                                           -- *int64
    smallint_val    SMALLINT NOT NULL DEFAULT 0,                      -- int64 (или int16 по соглашению)
    smallint_null   SMALLINT,                                         -- *int64

    -- REAL affinity
    real_val        REAL NOT NULL DEFAULT 0,                          -- float64
    real_null       REAL,                                             -- *float64
    double_val      DOUBLE NOT NULL DEFAULT 0,                        -- float64
    double_null     DOUBLE,                                           -- *float64
    float_val       FLOAT NOT NULL DEFAULT 0,                         -- float64
    float_null      FLOAT,                                            -- *float64

    -- NUMERIC affinity
    numeric_val     NUMERIC NOT NULL DEFAULT 0,                       -- float64 / string / decimal — решать конфигом
    numeric_null    NUMERIC,
    decimal_val     DECIMAL(10, 2) NOT NULL DEFAULT 0,                -- то же самое
    decimal_null    DECIMAL(10, 2),

    -- BOOLEAN: affinity NUMERIC, хранится как 0/1
    bool_val        BOOLEAN NOT NULL DEFAULT 0 CHECK (bool_val IN (0, 1)),            -- bool
    bool_null       BOOLEAN CHECK (bool_null IS NULL OR bool_null IN (0, 1)),         -- *bool

    -- TEXT affinity
    text_val        TEXT NOT NULL DEFAULT '',                         -- string
    text_null       TEXT,                                             -- *string
    varchar_val     VARCHAR(255) NOT NULL DEFAULT '',                 -- string
    varchar_null    VARCHAR(255),                                     -- *string
    char_val        CHAR(10) NOT NULL DEFAULT '',                     -- string
    char_null       CHAR(10),                                         -- *string

    -- BLOB affinity
    blob_val        BLOB NOT NULL DEFAULT x'',                        -- []byte
    blob_null       BLOB,                                             -- []byte (nil = NULL)

    -- Дата/время: affinity NUMERIC, по факту обычно TEXT в ISO-8601
    date_val        DATE NOT NULL DEFAULT CURRENT_DATE,               -- time.Time
    date_null       DATE,                                             -- *time.Time
    datetime_val    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,      -- time.Time
    datetime_null   DATETIME,                                         -- *time.Time
    timestamp_val   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,     -- time.Time
    timestamp_null  TIMESTAMP,                                        -- *time.Time

    -- JSON: affinity NUMERIC, хранится как TEXT
    json_val        JSON NOT NULL DEFAULT '{}',                       -- json.RawMessage / string
    json_null       JSON,

    -- Без типа: affinity BLOB, может хранить что угодно
    untyped_val     NOT NULL DEFAULT 0,                               -- any
    untyped_null                                                      -- any
);

-- 2. STRICT-таблица (SQLite 3.37+): разрешены только INT, INTEGER, REAL, TEXT, BLOB, ANY,
-- и движок реально проверяет типы при вставке.
CREATE TABLE all_types_strict (
    id              INTEGER PRIMARY KEY,

    int_val         INT NOT NULL DEFAULT 0,                           -- int64
    int_null        INT,                                              -- *int64
    integer_val     INTEGER NOT NULL DEFAULT 0,                       -- int64
    integer_null    INTEGER,                                          -- *int64
    real_val        REAL NOT NULL DEFAULT 0,                          -- float64
    real_null       REAL,                                             -- *float64
    text_val        TEXT NOT NULL DEFAULT '',                         -- string
    text_null       TEXT,                                             -- *string
    blob_val        BLOB NOT NULL DEFAULT x'',                        -- []byte
    blob_null       BLOB,                                             -- []byte
    any_val         ANY NOT NULL DEFAULT 0,                           -- any
    any_null        ANY                                               -- any
) STRICT;

-- 3. Подводные камни для вывода типов и nullability.
CREATE TABLE quirks (
    -- PRIMARY KEY, который НЕ INTEGER, в обычной таблице допускает NULL
    -- (исторический баг SQLite). pragma_table_info покажет notnull = 0.
    code            TEXT PRIMARY KEY,                                 -- *string (!), а не string

    -- Правило "содержит INT -> INTEGER" срабатывает раньше остальных:
    floating_point  FLOATING POINT,                                   -- affinity INTEGER (!)
    charint         CHARINT,                                          -- affinity INTEGER, а не TEXT (!)

    -- Неизвестное имя типа -> NUMERIC:
    string_val      STRING                                            -- affinity NUMERIC (!), не TEXT: '007' сохранится как 7
);

-- Тестовые данные: одна полностью заполненная строка и одна, где все nullable-колонки NULL.
INSERT INTO all_types (
    int_val, int_null, bigint_val, bigint_null, smallint_val, smallint_null,
    real_val, real_null, double_val, double_null, float_val, float_null,
    numeric_val, numeric_null, decimal_val, decimal_null,
    bool_val, bool_null,
    text_val, text_null, varchar_val, varchar_null, char_val, char_null,
    blob_val, blob_null,
    date_val, date_null, datetime_val, datetime_null, timestamp_val, timestamp_null,
    json_val, json_null,
    untyped_val, untyped_null
) VALUES (
    42, 42, 9223372036854775807, -9223372036854775808, 7, 7,
    3.14, 3.14, 2.718281828, 2.718281828, 1.5, 1.5,
    100.5, 100.5, 19.99, 19.99,
    1, 0,
    'hello', 'hello', 'varchar', 'varchar', 'char', 'char',
    x'DEADBEEF', x'CAFEBABE',
    '2026-11-27', '2026-11-27', '2026-11-27 02:30:00', '2026-11-27 02:30:00', '2026-11-27 02:30:00', '2026-11-27 02:30:00',
    '{"a": 1}', '[1, 2, 3]',
    'text in untyped', 123
);

INSERT INTO all_types DEFAULT VALUES;

INSERT INTO all_types_strict (
    int_val, int_null, integer_val, integer_null,
    real_val, real_null, text_val, text_null,
    blob_val, blob_null, any_val, any_null
) VALUES (
    1, 1, 2, 2,
    1.5, 1.5, 'strict', 'strict',
    x'01', x'02', 'any text', 3.5
);

INSERT INTO all_types_strict DEFAULT VALUES;

INSERT INTO quirks (code, floating_point, charint, string_val) VALUES
    ('A1', 1.5, '12', '007'),
    (NULL, NULL, NULL, NULL);  -- да, NULL в PRIMARY KEY проходит