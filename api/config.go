package api

type Config struct {
	TargetSchema    string           `yaml:"targetSchema"`
	TemplateConfigs []TemplateConfig `yaml:"templateConfigs"`
}

type DbConfig struct {
	// Connection configuration
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

type TemplateConfig struct {
	TemplatePath string `yaml:"templatePath"`
}
