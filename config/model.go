package config

import "fmt"

type Config struct {
	BaseUrl    string
	Username   string
	Password   string
	OutputDir  string
	RetryCount int
}

type Arguments struct {
	ConfigPath   string
	IssueKey     string
	Script       string
	ShowHelp     bool
	ShowProgress bool
}

func validateStringArg(value string, error_msg string, valid_flag *bool) {
	if value == "" {
		fmt.Println(error_msg)
		*valid_flag = false
	}
}

func (c *Config) Validate() bool {
	valid := true
	validateStringArg(c.Username, "- Please provide username", &valid)
	validateStringArg(c.Password, "- Please provide password", &valid)
	validateStringArg(c.BaseUrl, "- Please provide base_url", &valid)
	return valid
}

func (a *Arguments) Validate() bool {
	valid := true
	validateStringArg(a.IssueKey, "- Please provide issue_key", &valid)
	return valid
}
