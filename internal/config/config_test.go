package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	path := "/home/zyven/go/src/biter/config/config.yaml"
	config := GetConfig(path)
	t.Log(config.Http)
	t.Log(config.Basic)
	//PersistAccount(&model.AccountRequest{
	//	Username: "3220191000",
	//	Password: "zyven",
	//})
	ac, _ := ReadAccount()
	t.Log(ac)
	t.Log(GetDefaultGlobalConfigPath())
}