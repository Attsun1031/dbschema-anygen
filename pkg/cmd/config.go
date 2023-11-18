package cmd

type Config struct {
	Filterings []FilteringConfig `yaml:"filterings"`
}

type FilteringConfig struct {
	SchemaName string `yaml:"schemaName"`
}
