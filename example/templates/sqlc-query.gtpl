{{- range $i, $param := .TableParams -}}

{{ $tableName := $param.TableName -}}
{{ $tableNameFCU := $param.TableNameFCU -}}
{{ $columns := $param.Columns -}}

-- name: Get{{ $tableNameFCU }} :one
SELECT * FROM {{ $tableName }} WHERE id = $1;

-- name: Get{{ $tableNameFCU }}ByIds :many
SELECT * FROM {{ $tableName }} WHERE id = ANY($1::UUID []);

-- name: Create{{ $tableNameFCU }} :one
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

-- name: Delete{{ $tableNameFCU }} :exec
DELETE FROM {{ $tableName }} WHERE id = $1;

{{ end }}
