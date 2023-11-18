{{ $tableName := .TableName }}
-- name: Get{{ .tableName | ToUpper }} :one
SELECT * FROM {{ .tableName }} WHERE id = $1;
