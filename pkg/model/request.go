package model

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	ServerTypeCMCC   = "移动"
	ServerTypeWCDMA  = "联通"
	ServerTypeOrigin = "校园网"

	suffixCMCC  = "@yidong"
	suffixWCDMA = "@liantong"
)

type ChallengeRequest struct {
	Username string `json:"username"`
	Ip       string `json:"ip"`
}

type LoginRequest struct {
	Action   string `json:"action"`
	Username string `json:"username"`
	Password string `json:"password"`
	Acid     int    `json:"ac_id"`
	Ip       string `json:"ip"`
	Info     string `json:"info"`
	Chksum   string `json:"chksum"`
	N        int    `json:"n"`
	Type     int    `json:"type"`
}

type InfoRequest struct {
	Acid        int    `json:"ac_id"`
	Username    string `json:"username"`
	ClientIp    string `json:"client_ip"`
	AccessToken string `json:"access_token"`
}

type LogoutRequest struct {
	Action   string `json:"action"`
	Username string `json:"username"`
	Acid     int    `json:"ac_id"`
	Ip       string `json:"ip"`
}

type AccountRequest struct {
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	AccessToken string `json:"access_token"`
	Ip          string `json:"ip"`
	Server      string `json:"server"`
}

func Challenge(username string) url.Values {
	return url.Values{
		"username": {username},
		"ip":       {""},
	}
}

func Info(acid int, username, clientIp, accessToken string) url.Values {
	return url.Values{
		"ac_id":        {fmt.Sprint(acid)},
		"username":     {username},
		"client_ip":    {clientIp},
		"access_token": {accessToken},
	}
}

func Login(username, password string, acid int) url.Values {
	return url.Values{
		"action":   {"login"},
		"username": {username},
		"password": {password},
		"ac_id":    {fmt.Sprint(acid)},
		"ip":       {""},
		"info":     {},
		"chksum":   {},
		"n":        {"200"},
		"type":     {"1"},
	}
}

func Logout(username string) url.Values {
	return url.Values{
		"action":   {"logout"},
		"username": {username},
	}
}

func (a *AccountRequest) JSONString() (jsonStr string, err error) {
	jsonData, err := json.Marshal(a)
	if err != nil {
		return
	}
	jsonStr = string(jsonData)
	return
}

func (a *AccountRequest) JSONBytes() (jsonData []byte, err error) {
	return json.Marshal(a)
}

func (a *AccountRequest) String() string {
	return fmt.Sprintln("用户名:", a.Username)
}

func (a *AccountRequest) GenUsername() string {
	return addSuffix(a.Username, a.Server)
}

func addSuffix(name, server string) string {
	switch server {
	case ServerTypeCMCC:
		return name + suffixCMCC
	case ServerTypeWCDMA:
		return name + suffixWCDMA
	default:
		return name
	}
}
