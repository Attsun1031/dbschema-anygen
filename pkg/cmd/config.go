package cmd

type Config struct {
	TemplatePath string `yaml:"templatePath"`
	TargetSchema string `yaml:"targetSchema"`
}
