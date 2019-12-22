package utils

import (
	"biter/internal/log"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

func jsonpParam() string {
	return fmt.Sprintf("jsonp%d", time.Now().Unix())
}

// DoRequest do request with given url and params
func DoRequest(url string, val url.Values) (*http.Response, error) {
	val.Add("callback", jsonpParam())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("Request:%s found error:%v", url, err)
		return nil, err
	}

	req.URL.RawQuery = val.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Do request found error:%v", err)
		return nil, err
	}
	return resp, nil
}

func cutJsonp(jsonp string) (string, error) {
	infoRegex := regexp.MustCompile(`\{(.*)\}`)
	slices := infoRegex.FindStringSubmatch(jsonp)
	if len(slices) < 1 {
		return "", errors.New("invalid jsonp")
	}
	return slices[0], nil
}

// ParseRequest parse response json
func ParseRequest(url string, val url.Values, res interface{}) error {
	resp, err := DoRequest(url, val)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Reading response body found error:%v", err)
		return err
	}

	keyInfo, err := cutJsonp(string(raw))
	if err != nil {
		log.Errorf("Parse key info found error:%v", err)
		return err
	}

	if err := json.Unmarshal([]byte(keyInfo), &res); err != nil {
		log.Errorf("Unmarshal response body found error:%v", err)
		return err
	}
	return nil
}