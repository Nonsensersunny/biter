package cli

import (
	"biter/core"
	"biter/internal/config"
	"biter/internal/log"
	"biter/pkg/model"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// VERSION specify the version
const VERSION = "v0.0.2"

// RELEASE release link
const RELEASE = "https://github.com/Nonsensersunny/biter/releases/latest"
const updateTimeout = 3 * time.Second

// Func type of command function
type Func func(cmd string, params map[string]string)

// Client is an empty struct
type Client struct{}

var (
	// DefaultClient default client
	DefaultClient = &Client{}
	serverTypes   = map[string]string{
		"campus": model.ServerTypeOrigin,
		"mobile": model.ServerTypeCMCC,
		"unicom": model.ServerTypeWCDMA,
	}
	// CmdMap command map
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
	// CmdHelper command helper
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
SetServer:
	fmt.Print("Server[Default campus](0 campus | 1 mobile | 2 unicom):")
	server := 0
	_, _ = fmt.Scanln(&server)
	switch server {
	case 0:
		account.Server = model.ServerTypeOrigin
	case 1:
		account.Server = model.ServerTypeCMCC
	case 2:
		account.Server = model.ServerTypeWCDMA
	default:
		goto SetServer
	}
	return account
}

// Login login
func (c *Client) Login(cmd string, params map[string]string) {
	var account model.AccountRequest
	ok, err := strconv.ParseBool(params["create"])
	if err != nil {
		log.Errorf("解析参数失败:%v", err)
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
		log.Errorf("登录失败:%v", err)
		os.Exit(1)
	}
	log.Info("登录成功!")
	log.Infof("当前IP:%v", info.ClientIp)
}

// Logout logout
func (c *Client) Logout(cmd string, params map[string]string) {
	cli := &core.Core{Config: config.GetConfig(params["global"])}
	if err := cli.Logout(); err != nil {
		log.Errorf("登出失败!")
		os.Exit(1)
	}
	log.Info("登出成功!")
}

// Account account
func (c *Client) Account(cmd string, params map[string]string) {
	conf := config.GetConfig(params["global"])
	cli := &core.Core{Config: conf}
	cli.AccountInfo(model.AccountRequest{
		Username: conf.Basic.Username,
		Password: conf.Basic.Password,
	})
}

// Config config
func (c *Client) Config(cmd string, params map[string]string) {
	account := getCustomizedBasicConfig()
	if err := config.PersistAccount(&account); err != nil {
		log.Errorf("配置保存失败:%v", err)
		os.Exit(1)
	}
	log.Info("成功保存配置!")
}

// CmdHelp command help
func (c *Client) CmdHelp(cmd string, params map[string]string) {
	fmt.Println(c.CmdList())
}

// CmdList list of commands
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

// Update update
func (c *Client) Update(cmd string, params map[string]string) {
	client := http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
			conn, err = net.DialTimeout(network, addr, updateTimeout)
			if err != nil {
				return nil, err
			}
			_ = conn.SetDeadline(time.Now().Add(updateTimeout))
			return conn, nil
		},
	}

	req, err := http.NewRequest(http.MethodGet, RELEASE, nil)
	if err != nil {
		log.Errorf("获取更新失败:%v", err)
		return
	}

	res, err := client.RoundTrip(req)
	if err != nil {
		log.Errorf("请求错误:%v", err)
		return
	}

	newAddr := res.Header.Get("Location")
	arr := strings.Split(newAddr, "/")
	newVersion := arr[len(arr)-1]
	log.Infof("最新版本:%s, 当前版本:%s", newVersion, VERSION)
	if newVersion != VERSION {
		log.Infof("新版本链接:%s", newAddr)
	}
}
