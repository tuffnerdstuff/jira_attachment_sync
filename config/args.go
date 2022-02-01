package config

import (
	"flag"
	"fmt"
)

var Username string
var Password string
var BaseUrl string
var OutputDir string = "."
var IssueKey string
var ShowHelp bool
var HttpRetryCount int

func ParseArgs() bool {
	flag.StringVar(&Username, "username", "", "Your Jira username")
	flag.StringVar(&Password, "password", "", "Your Jira password")
	flag.StringVar(&BaseUrl, "base_url", "", "Your Jira base URL (e.g. https://example.com/rest/api/2)")
	flag.StringVar(&IssueKey, "issue", "", "The issue ID")
	flag.StringVar(&OutputDir, "output", ".", "The base path to your downloaded attachments")
	flag.IntVar(&HttpRetryCount, "retry_count", 5, "How many retries should be performed until HTTP request is abandoned")
	flag.BoolVar(&ShowHelp, "help", false, "Show usage info")
	flag.Parse()

	if ShowHelp || !validateArgs() {
		flag.Usage()
		return false
	}

	return true
}

func validateStringArg(value string, error_msg string, valid_flag *bool) {
	if value == "" {
		fmt.Println(error_msg)
		*valid_flag = false
	}
}

func validateArgs() bool {
	valid := true
	validateStringArg(Username, "- Please provide username", &valid)
	validateStringArg(Password, "- Please provide password", &valid)
	validateStringArg(BaseUrl, "- Please provide base_url", &valid)
	validateStringArg(IssueKey, "- Please provide issue_key", &valid)
	return valid
}
