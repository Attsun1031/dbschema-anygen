{{ $tableName := .TableName }}
{{ $tableNameFCU := .TableNameFCU }}
-- name: Get{{ $tableNameFCU }} :one
SELECT * FROM {{ $tableName }} WHERE id = $1;

-- name: Get{{ $tableNameFCU }}ByIds :many
SELECT * FROM {{ $tableName }} WHERE id = ANY($1::UUID []);
