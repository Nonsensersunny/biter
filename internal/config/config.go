package config

import (
	"biter/internal/log"
	"biter/pkg/model"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	PersistentConfig = "biter"
	DefaultGlobalConfig = "config.yaml"
)

var (
	RootPath string
	DefaultHttpConfig = HttpConfig{
		Portal:           "10.0.0.55",
		ChallengeUrl:     "cgi-bin/get_challenge",
		SrunPortalUrl:    "cgi-bin/srun_portal",
		SucceedUrlOrigin: "srun_portal_pc_succeed.php",
		SucceedUrlCMCC:   "srun_portal_pc_succeed_yys.php",
		SucceedUrlWCDMA:  "srun_portal_pc_succeed_yys_cucc.php",
	}
)

type GlobalConfig struct {
	Basic *BasicConfig `yaml:"basic" json:"basic"`
	Http  *HttpConfig  `yaml:"http" json:"http"`
}

type BasicConfig struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

type HttpConfig struct {
	Portal           string `yaml:"portal" json:"portal"`
	ChallengeUrl     string `yaml:"challenge-url" json:"challenge_url"`
	SrunPortalUrl    string `yaml:"srun-portal-url" json:"srun_portal_url"`
	SucceedUrlOrigin string `yaml:"succeed-url-origin" json:"succeed_url_origin"`
	SucceedUrlCMCC   string `yaml:"succeed-url-cmcc" json:"succeed_url_cmcc"`
	SucceedUrlWCDMA  string `yaml:"succeed-url-wcdma" json:"succeed_url_wcdma"`
}

// GetPortal get http portal
func (h *HttpConfig) GetPortal() string {
	if !strings.Contains(h.Portal, "http://") {
		h.Portal = fmt.Sprintf("http://%s", h.Portal)
	}
	if string(h.Portal[len(h.Portal) - 1]) != "/" {
		h.Portal = fmt.Sprintf("%s/", h.Portal)
	}
	return h.Portal
}

// GetPurePortal return pure portal config not containing `http://` prefix
func (h *HttpConfig) GetPurePortal() string {
	reg := regexp.MustCompile(`(^\d.*\d$)`)
	return reg.FindStringSubmatch(h.Portal)[0]
}

// GetChallengeUrl get challenge url
func (h *HttpConfig) GetChallengeUrl() string {
	if !strings.Contains(h.ChallengeUrl, h.GetPortal()) {
		h.ChallengeUrl = fmt.Sprintf("%s%s", h.GetPortal(), h.ChallengeUrl)
	}
	return h.ChallengeUrl
}

// GetSrunPortalUrl get srun portal url
func (h *HttpConfig) GetSrunPortalUrl() string {
	if !strings.Contains(h.SrunPortalUrl, h.GetPortal()) {
		h.SrunPortalUrl = fmt.Sprintf("%s%s", h.GetPortal(), h.SrunPortalUrl)
	}
	return h.SrunPortalUrl
}

// GetSucceedUrlOrigin get succeed url origin
func (h *HttpConfig) GetSucceedUrlOrigin() string {
	if !strings.Contains(h.SucceedUrlOrigin, h.GetPortal()) {
		h.SucceedUrlOrigin = fmt.Sprintf("%s%s", h.GetPortal(), h.SucceedUrlOrigin)
	}
	return h.SucceedUrlOrigin
}

// GetSucceedUrlCMCC get succeed url cmcc
func (h *HttpConfig) GetSucceedUrlCMCC() string {
	if !strings.Contains(h.SucceedUrlCMCC, h.GetPortal()) {
		h.SucceedUrlCMCC = fmt.Sprintf("%s%s", h.GetPortal(), h.SucceedUrlCMCC)
	}
	return h.SucceedUrlCMCC
}

// GetSucceedUrlWCDMA get succeed url wcdma
func (h *HttpConfig) GetSucceedUrlWCDMA() string {
	if !strings.Contains(h.SucceedUrlWCDMA, h.GetPortal()) {
		h.SucceedUrlWCDMA = fmt.Sprintf("%s%s", h.GetPortal(), h.SucceedUrlWCDMA)
	}
	return h.SucceedUrlWCDMA
}

// GetConfig get global config
func GetConfig(path string) *GlobalConfig {
	log.Info("加载配置")
	if path == "" {
		return GetDefaultConfig()
	}
	var config *GlobalConfig
	ymlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Reading config found error:%v, using default config", err)
		return GetDefaultConfig()
	}
	if err = yaml.Unmarshal(ymlFile, &config); err != nil {
		log.Fatalf("Unmarshal config found error:%v", err)
		return GetDefaultConfig()
	}
	return config
}

// GetDefaultConfig get default global config
func GetDefaultConfig() *GlobalConfig {
	log.Info("Loading default config...")
	conf, err := ReadAccount()
	if err != nil {
		log.Errorf("Reading default config failed:%v", err)
		os.Exit(1)
	}
	return &GlobalConfig{
		Basic: &BasicConfig{
			Username: conf.Username,
			Password: conf.Password,
		},
		Http:  &DefaultHttpConfig,
	}
}

// GetDefaultGlobalConfig get default config
func GetDefaultGlobalConfigPath() string {
	path := filepath.Join(RootPath, ".biter")
	return fmt.Sprintf("%s/%s", path, DefaultGlobalConfig)
}

func getAccountFilePath() (src string, err error) {
	path := filepath.Join(RootPath, ".biter")
	if _, se := os.Stat(path); se != nil {
		if me := os.MkdirAll(path, 0755); me != nil {
			log.Errorf("mkdir:%s found error:%v", path, err)
			return
		}
	}
	src = filepath.Join(path, PersistentConfig)
	return
}

func PersistAccount(account *model.AccountRequest) (err error) {
	src, err := getAccountFilePath()
	if err != nil {
		log.Errorf("Get config file found error:%v", err)
		return
	}
	file, err := os.OpenFile(src, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		log.Errorf("Open config file found error:%v", err)
		return
	}
	defer file.Close()

	jsonBytes, err := account.JSONBytes()
	if err != nil {
		log.Errorf("Get account json bytes found error:%v", err)
		return
	}
	str := base64.StdEncoding.EncodeToString(jsonBytes)
	if _, err := file.WriteString(str); err != nil {
		log.Errorf("Writing config file found error:%v", err)
		return err
	}
	return nil
}

func ReadAccount() (account model.AccountRequest, err error) {
	src, err := getAccountFilePath()
	if err != nil {
		log.Errorf("Get config file found error:%v", err)
		return
	}
	file, err := os.Open(src)
	if err != nil {
		log.Errorf("Open config file found error:%v", err)
		return
	}
	defer file.Close()

	readBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf("Read config file found error:%v", err)
		return
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(string(readBytes))
	if err != nil {
		log.Errorf("Decode config found error:%v", err)
		return
	}

	if err = json.Unmarshal(decodedBytes, &account); err != nil {
		log.Errorf("Unmarshal config found error:%v", err)
		return
	}
	return account, nil
}

func init() {
	curUser, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to read system user info:%v", err)
	} else {
		RootPath = curUser.HomeDir
	}
}