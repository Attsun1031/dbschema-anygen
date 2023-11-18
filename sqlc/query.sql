-- name: GetColumnDefinitions :many
SELECT
    CAST(table_schema AS TEXT) AS table_schema,
    cast(table_name AS TEXT) AS table_name,
    cast(column_name AS TEXT) AS column_name,
    cast(data_type AS TEXT) AS data_type,
    cast(ordinal_position AS INTEGER) AS ordinal_position
FROM information_schema.columns
WHERE table_schema = $1
ORDER BY table_schema, table_name, ordinal_position;
