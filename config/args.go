package config

import (
	"flag"
)

func ParseArgs(args *Arguments) {
	// Set args
	flag.StringVar(&args.ConfigPath, "config_path", "./jsa-config.toml", "The path to the configuration file (toml)")
	flag.StringVar(&args.IssueKey, "issue", "", "The issue ID")
	flag.BoolVar(&args.ShowHelp, "help", false, "Show usage info")

	flag.Parse()

}

func ParseConfig(config *Config) {
	// Override config with args
	flag.StringVar(&config.Username, "username", config.Username, "Your Jira username")
	flag.StringVar(&config.Password, "password", config.Password, "Your Jira password")
	flag.StringVar(&config.BaseUrl, "base_url", config.BaseUrl, "Your Jira base URL (e.g. https://example.com/rest/api/2)")
	flag.StringVar(&config.OutputDir, "output", config.OutputDir, "The base path to your downloaded attachments")
	flag.IntVar(&config.RetryCount, "retry_count", config.RetryCount, "How many retries should be performed until HTTP request is abandoned")
}
