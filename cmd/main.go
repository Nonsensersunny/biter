package main

import (
	"biter/cmd/cli"
	"biter/internal/config"
	"biter/pkg/model"
	"flag"
	"strconv"
)

func main() {
	var (
		cmd    string
		params map[string]string

		globalConfigPath string
		httpConfigPath string
		basicConfigPath string
		serverType string
		customizeBasicConfig bool
	)

	flag.StringVar(&globalConfigPath, "g", config.GetDefaultGlobalConfigPath(), "customize configuration file")
	flag.StringVar(&httpConfigPath, "s", config.GetDefaultGlobalConfigPath(), "customize server configuration file")
	flag.StringVar(&httpConfigPath, "t", model.ServerTypeOrigin, "customize server type")
	flag.StringVar(&basicConfigPath, "a", config.GetDefaultGlobalConfigPath(), "customize account configuration file")
	flag.BoolVar(&customizeBasicConfig, "c", false, "customize account settings")

	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		cmd = args[0]
		params = map[string]string{
			"global": globalConfigPath,
			"http": httpConfigPath,
			"basic": basicConfigPath,
			"type": serverType,
			"create": strconv.FormatBool(customizeBasicConfig),
		}
	} else {
		cmd = "login"
	}
	cli.Execute(cmd, params)
}
