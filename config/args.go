package config

import (
	"flag"
)

func ParseArgs(args *Arguments) {
	// Set args
	flag.StringVar(&args.ConfigPath, "configPath", "./jas-config.toml", "The path to the configuration file (toml)")
	flag.StringVar(&args.URI, "URI", "", "The URI of the Jira issue(s)")
	flag.StringVar(&args.Script, "script", "", "The path of the script to run after download is finished (download folder is passed as argument)")
	flag.BoolVar(&args.ShowProgress, "progress", false, "Show animated progress")
	flag.BoolVar(&args.ShowHelp, "help", false, "Show usage info")
	flag.Parse()
}
