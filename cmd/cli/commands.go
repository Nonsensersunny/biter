package cli

import (
	"biter/core"
	"biter/internal/config"
	"biter/internal/log"
	"biter/pkg/model"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const VERSION = "V0.1.1"

type Func func(cmd string, params map[string]string)

type Client struct{}

var (
	DefaultClient = &Client{}
	serverTypes   = map[string]string{
		"campus": model.ServerTypeOrigin,
		"mobile": model.ServerTypeCMCC,
		"unicom": model.ServerTypeWCDMA,
	}
	CmdMap = map[string]Func{
		"config":  DefaultClient.Config,
		"login":   DefaultClient.Login,
		"logout":  DefaultClient.Logout,
		"account": DefaultClient.Account,
		"update":  DefaultClient.Update,
		"help":    DefaultClient.CmdHelp,
	}
	optionDocs = map[string]string{
		"-a": "customize account configuration file path",
		"-c": "customize account settings",
		"-g": "customize full configuration file path",
		"-h": "customize server configuration file path",
	}
	CmdHelper = map[string][]string{
		"config":  {"biter config", "Account settings"},
		"login":   {"biter [login] [campus|mobile|unicom]", "Login to network"},
		"logout":  {"biter logout", "Network logout"},
		"account": {"biter account", "Get account info"},
		"update":  {"biter update", "Update biter tool"},
	}
)

func getCustomizedBasicConfig() model.AccountRequest {
	account := model.AccountRequest{}
	fmt.Print("Username:")
	_, _ = fmt.Scanln(&account.Username)
	fmt.Print("Password:")
	_, _ = fmt.Scanln(&account.Password)
	return account
}

func (c *Client) Login(cmd string, params map[string]string) {
	var account model.AccountRequest
	ok, err := strconv.ParseBool(params["create"])
	if err != nil {
		log.Errorf("Parse parameter found error:%v", err)
		os.Exit(1)
	}
	if ok {
		account = getCustomizedBasicConfig()
	} else {
		account, err = config.ReadAccount()
		if err != nil {
			log.Error("Config missing")
			os.Exit(1)
		}
	}
	account.Server = params["type"]
	cli := &core.Core{Config: config.GetConfig(params["global"])}
	info, err := cli.Login(&account)
	if err != nil {
		log.Errorf("Failed to login:%v", err)
		os.Exit(1)
	}
	log.Info("登录成功!")
	log.Infof("当前IP:%v", info.ClientIp)
}

func (c *Client) Logout(cmd string, params map[string]string) {
	cli := &core.Core{Config: config.GetConfig(params["global"])}
	if err := cli.Logout(); err != nil {
		log.Errorf("登出失败!")
		os.Exit(1)
	}
	log.Info("登出成功!")
}

func (c *Client) Account(cmd string, params map[string]string) {
	conf := config.GetConfig(params["global"])
	cli := &core.Core{Config: conf,}
	cli.AccountInfo(model.AccountRequest{
		Username: conf.Basic.Username,
		Password: conf.Basic.Password,
	})
}

func (c *Client) Config(cmd string, params map[string]string) {
	account := getCustomizedBasicConfig()
	if err := config.PersistAccount(&account); err != nil {
		log.Errorf("Persist config found error:%v", err)
		os.Exit(1)
	}
	log.Info("Persist config succeeded!")
}

func (c *Client) CmdHelp(cmd string, params map[string]string) {
	fmt.Println(c.CmdList())
}

func (Client) CmdList() string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("biter %s\r\n", VERSION))
	sb.WriteString("A simple tool for login to campus network in BIT\r\n")
	sb.WriteString(fmt.Sprint("\r\nUsage:	biter COMMAND [OPTIONS]\r\n"))
	sb.WriteString("\r\nOptions:\r\n")
	for k, v := range optionDocs {
		sb.WriteString(fmt.Sprintf("  %-10s%-20s\r\n", k, v))
	}
	sb.WriteString("\r\nCommands:\r\n")
	for k, v := range CmdHelper {
		sb.WriteString(fmt.Sprintf("  %-10s%-20s\r\n", k, v[1]))
	}
	return sb.String()
}

func (c *Client) Update(cmd string, params map[string]string) {

}