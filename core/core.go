package core

import (
	"biter/internal/config"
	"biter/internal/errors"
	"biter/internal/log"
	"biter/internal/utils"
	"biter/pkg/hash"
	"biter/pkg/model"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Core struct {
	Config *config.GlobalConfig
}

func (c *Core) getAcid() (acid int, err error) {
	var (
		client = http.DefaultClient
		reg, _ = regexp.Compile(`index_[\d]\.html`)
		demoUrl = "http://t.cn"
		request *http.Request
		response *http.Response
	)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if strings.Contains(req.URL.String(), c.Config.Http.GetPurePortal()) {
			if reg.MatchString(req.URL.String()) {
				res := reg.FindString(req.URL.String())
				acids := strings.TrimRight(strings.TrimLeft(res, "index_"), ".html")
				acid, err = strconv.Atoi(acids)
				if err != nil {
					log.Errorf("Invalid URL:%s with error:%v", acids, err)
				}
				return nil
			}
		}
		return errors.ErrConnected
	}

	request, err = http.NewRequest(http.MethodGet, demoUrl, nil)
	if err != nil {
		log.Error(err)
		return acid, err
	}

	response, err = client.Do(request)
	switch err {
	case errors.ErrConnected:
		return
	case nil:
		_ = response.Body.Close()
		return acid, nil
	default:
		err = errors.ErrRequest
		return
	}
}

func (c *Core) getChallenge() (res model.ChallengeResponse, err error) {
	val := model.Challenge(c.Config.Basic.Username)
	err = utils.ParseRequest(c.Config.Http.GetChallengeUrl(), val, &res)
	return
}

func (*Core) parseHtml(url string, val url.Values) (err error) {
	resp, err := utils.DoRequest(url, val)
	if err != nil {
		log.Debug(err)
		err = errors.ErrRequest
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Error(err)
		err = errors.ErrRequest
		return
	}

	bytes := doc.Find("span#sum_bytes").Last().Text()
	times := doc.Find("span#sum_seconds").Text()
	balance := doc.Find("span#user_balance").Text()
	fmt.Println("已用流量:", bytes)
	fmt.Println("已用时长:", times)
	fmt.Println("账户余额:", balance)
	return
}

func (c *Core) Login(account *model.AccountRequest) (res model.InfoRequest, err error) {
	log.Infof("当前账号:%s", c.Config.Basic.Username)
	acid, err := c.getAcid()
	if err != nil {
		log.Errorf("Get acid found error:%v", err)
		err = errors.ErrConnected
		return
	}
	if acid == 1 && account.Server != model.ServerTypeOrigin {
		log.Warning(errors.ErrAcid)
		account.Server = model.ServerTypeOrigin
	}

	username := account.GenUsername()
	formLogin := model.Login(username, account.Password, acid)

	challenge, err := c.getChallenge()
	if err != nil {
		log.Errorf("Get challenge found error:%v", err)
		err = errors.ErrRequest
		return
	}

	token := challenge.Challenge
	ip := challenge.ClientIp

	formLogin.Set("ip", ip)
	formLogin.Set("info", hash.GenInfo(formLogin, token))
	formLogin.Set("password", hash.PwdHmd5("", token))
	formLogin.Set("chksum", hash.Checksum(formLogin, token))

	actionResponse := model.ActionResponse{}
	if err := utils.ParseRequest(c.Config.Http.GetSrunPortalUrl(), formLogin, &actionResponse); err != nil {
		log.Errorf("Fetch request found error:%v", err)
		err = errors.ErrRequest
		return model.InfoRequest{}, nil
	}
	if actionResponse.Res != "ok" {
		msg := actionResponse.Res
		if msg == "" {
			msg = actionResponse.ErrorMsg
		}
		log.Errorf("Login failed:%v", msg)
		err = errors.ErrFailed
		return
	}

	res = model.InfoRequest{
		Acid:        acid,
		Username:    username,
		ClientIp:    challenge.ClientIp,
		AccessToken: challenge.Challenge,
	}
	return
}

func (c *Core) AccountInfo(account model.AccountRequest) (err error) {
	info := model.Info(1, account.Username, account.Ip, account.AccessToken)
	log.Infof("Server:%v", account.Server)
	if account.Server != model.ServerTypeOrigin {
		return nil
	}
	err = c.parseHtml(c.Config.Http.GetSucceedUrlOrigin(), info)
	return
}

func (c *Core) Logout() (err error) {
	logout := model.Logout(c.Config.Basic.Username)
	actionResponse := model.ActionResponse{}
	if err = utils.ParseRequest(c.Config.Http.GetSrunPortalUrl(), logout, &actionResponse); err != nil {
		log.Errorf("Do logout found error:%v", err)
		err = errors.ErrRequest
		return
	}
	if actionResponse.Error != "ok" {
		log.Errorf("Logout error:%v", err)
		err = errors.ErrRequest
	}
	return
}
