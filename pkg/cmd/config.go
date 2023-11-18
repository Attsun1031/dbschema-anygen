package cmd

type Config struct {
	TemplatePath string            `yaml:"templatePath"`
	Filterings   []FilteringConfig `yaml:"filterings"`
}

type FilteringConfig struct {
	SchemaName string `yaml:"schemaName"`
}
