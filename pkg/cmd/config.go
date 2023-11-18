package cmd

type Config struct {
	TemplateConfigs []TemplateConfig `yaml:"templateConfigs"`
}

type TemplateConfig struct {
	TemplatePath string `yaml:"templatePath"`
	TargetSchema string `yaml:"targetSchema"`
}
