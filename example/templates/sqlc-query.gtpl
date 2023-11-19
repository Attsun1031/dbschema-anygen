{{- range $i, $param := .TableParams -}}

{{ $tableName := $param.TableName -}}
{{ $tableNameCamelFU := $param.TableNameCamelFU -}}
{{ $columns := $param.Columns -}}

-- name: Get{{ $tableNameCamelFU }} :one
SELECT * FROM {{ $tableName }} WHERE id = $1;

-- name: Get{{ $tableNameCamelFU }}ByIds :many
SELECT * FROM {{ $tableName }} WHERE id = ANY($1::UUID []);

-- name: Create{{ $tableNameCamelFU }} :one
INSERT INTO {{ $tableName }} (
{{- range $i, $col := $columns }}
    {{- if ne $i 0 }}, {{ end }}
    {{- $col.ColumnName }}
{{- end -}}
) VALUES (
{{- range $i, $col := $columns }}
    {{- if ne $i 0 }}, {{ end -}}
    ${{- AddNum $i 1 }}
{{- end -}}
) RETURNING *;

-- name: Delete{{ $tableNameCamelFU }} :exec
DELETE FROM {{ $tableName }} WHERE id = $1;

{{ end }}
