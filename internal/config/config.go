package config

import (
	"biter/internal/log"
	"biter/pkg/model"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	// PersistentConfig default config name
	PersistentConfig = "biter"
	// DefaultGlobalConfig default global config name
	DefaultGlobalConfig = "config.yaml"
)

var (
	// RootPath root path
	RootPath string
	// DefaultHttpConfig default HTTP config
	DefaultHttpConfig = HttpConfig{
		Portal:           "10.0.0.55",
		ChallengeUrl:     "cgi-bin/get_challenge",
		SrunPortalUrl:    "cgi-bin/srun_portal",
		SucceedUrlOrigin: "srun_portal_pc_succeed.php",
		SucceedUrlCMCC:   "srun_portal_pc_succeed_yys.php",
		SucceedUrlWCDMA:  "srun_portal_pc_succeed_yys_cucc.php",
	}
)

// GlobalConfig global config struct
type GlobalConfig struct {
	Basic *BasicConfig `yaml:"basic" json:"basic"`
	Http  *HttpConfig  `yaml:"http" json:"http"`
}

// BasicConfig basic config struct
type BasicConfig struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

// HttpConfig HTTP config struct
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
	if string(h.Portal[len(h.Portal)-1]) != "/" {
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
		log.Warningf("读取配置失败:%v, 应用默认配置", err)
		return GetDefaultConfig()
	}
	if err = yaml.Unmarshal(ymlFile, &config); err != nil {
		log.Fatalf("配置解析失败:%v", err)
		return GetDefaultConfig()
	}
	return config
}

// GetDefaultConfig get default global config
func GetDefaultConfig() *GlobalConfig {
	log.Info("加载默认配置...")
	conf, err := ReadAccount()
	if err != nil {
		log.Errorf("读取默认配置失败:%v", err)
		os.Exit(1)
	}
	return &GlobalConfig{
		Basic: &BasicConfig{
			Username: conf.Username,
			Password: conf.Password,
		},
		Http: &DefaultHttpConfig,
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
			log.Errorf("创建文件夹:%s 失败:%v", path, err)
			return
		}
	}
	src = filepath.Join(path, PersistentConfig)
	return
}

func PersistAccount(account *model.AccountRequest) (err error) {
	src, err := getAccountFilePath()
	if err != nil {
		log.Errorf("读取配置文件失败:%v", err)
		return
	}
	file, err := os.OpenFile(src, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		log.Errorf("无法打开配置文件:%v", err)
		return
	}
	defer file.Close()

	jsonBytes, err := account.JSONBytes()
	if err != nil {
		log.Errorf("无法读取配置字节流:%v", err)
		return
	}
	str := base64.StdEncoding.EncodeToString(jsonBytes)
	if _, err := file.WriteString(str); err != nil {
		log.Errorf("配置写入失败:%v", err)
		return err
	}
	return nil
}

func ReadAccount() (account model.AccountRequest, err error) {
	src, err := getAccountFilePath()
	if err != nil {
		log.Errorf("配置文件获取失败:%v", err)
		return
	}
	file, err := os.Open(src)
	if err != nil {
		log.Errorf("配置文件打开失败:%v", err)
		return
	}
	defer file.Close()

	readBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf("配置文件读取失败:%v", err)
		return
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(string(readBytes))
	if err != nil {
		log.Errorf("配置解码失败:%v", err)
		return
	}

	if err = json.Unmarshal(decodedBytes, &account); err != nil {
		log.Errorf("配置解析失败:%v", err)
		return
	}
	return account, nil
}

func init() {
	curUser, err := user.Current()
	if err != nil {
		log.Fatalf("无法获取系统用户信息:%v", err)
	} else {
		RootPath = curUser.HomeDir
	}
}
