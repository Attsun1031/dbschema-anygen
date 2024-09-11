package api

type Config struct {
	TargetSchema    string           `yaml:"targetSchema" json:"targetSchema"`
	TemplateConfigs []TemplateConfig `yaml:"templateConfigs" json:"templateConfigs"`
}

type TemplateConfig struct {
	TemplatePath string `yaml:"templatePath" json:"templatePath"`
	OutputPath   string `yaml:"outputPath" json:"outputPath"`
}
