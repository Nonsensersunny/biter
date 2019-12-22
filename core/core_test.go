package core

import (
	"biter/internal/config"
	"biter/pkg/model"
	"testing"
)

func TestCore_Login(t *testing.T) {
	path := "/home/zyven/go/src/biter/internal/config/config.yaml"
	core := Core{Config: config.GetConfig(path)}
	core.Login(&model.AccountRequest{

	})
}
