-- name: GetColumnDefinitions :many
SELECT
    CAST(table_schema AS TEXT) AS table_schema,
    CAST(table_name AS TEXT) AS table_name,
    CAST(column_name AS TEXT) AS column_name,
    CAST(data_type AS TEXT) AS data_type,
    CASE (is_nullable)
        WHEN 'YES' THEN TRUE
        ELSE FALSE
    END AS is_nullable,
    CAST(ordinal_position AS INTEGER) AS ordinal_position
FROM information_schema.columns
WHERE table_schema = $1
ORDER BY table_schema, table_name, ordinal_position;
