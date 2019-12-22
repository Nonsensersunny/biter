package utils

import (
	"biter/internal/config"
	"net/url"
	"testing"
)

func TestDoRequest(t *testing.T) {
	testUrl := "http://baidu.com"
	res, err := DoRequest(testUrl, url.Values{})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(res)
}

func TestParseRequest(t *testing.T) {
	conf := config.GetConfig()
	t.Log(conf)
}