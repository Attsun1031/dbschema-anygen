package api

type Config struct {
	TargetSchema    string           `yaml:"targetSchema" json:"targetSchema"`
	TemplateConfigs []TemplateConfig `yaml:"templateConfigs" json:"templateConfigs"`
	DbConfig        DbConfig         `yaml:"dbConfig" json:"dbConfig"`
}

type DbConfig struct {
	// Connection configuration
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	DbName   string `yaml:"dbName" json:"dbName"`
}

type TemplateConfig struct {
	TemplatePath string `yaml:"templatePath" json:"templatePath"`
	OutputPath   string `yaml:"outputPath" json:"outputPath"`
}
