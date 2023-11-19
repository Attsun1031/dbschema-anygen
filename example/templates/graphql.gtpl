{{ range $i, $param := .TableParams }}

{{ $tableName := $param.TableName -}}
{{ $tableNameCamelFU := $param.TableNameCamelFU -}}
{{ $columns := $param.Columns -}}
{{ $typeBody := GenGraphQLTypeBody $param -}}

type {{ $tableNameCamelFU }} implements Node {
    {{- range $i, $line := $typeBody }}
    {{ $line }}
    {{- end }}
}
{{- end }}
