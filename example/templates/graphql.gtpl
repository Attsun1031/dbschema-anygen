{{ range $i, $param := .TableParams }}

{{ $tableName := $param.TableName -}}
{{ $tableNameFCU := $param.TableNameFCU -}}
{{ $columns := $param.Columns -}}
{{ $typeBody := GenGraphQLTypeBody $param -}}

type {{ $tableNameFCU }} implements Node {
    {{- range $i, $line := $typeBody }}
    {{ $line }}
    {{- end }}
}
{{- end }}
