package cli

func Execute(cmd string, params map[string]string) {
	if handle, ok := CmdMap[cmd]; ok {
		handle(cmd, params)
	} else {
		DefaultClient.CmdHelp(cmd, params)
	}
}